package server

import (
	"errors"
	"fmt"
	"rush/attendance"
	"rush/golang/array"
	"rush/user"
	"strings"
	"time"
)

func (s *Server) AdminGetSession(id string) (SessionForAdmin, error) {
	session, err := s.sessionRepo.Get(id)
	if err != nil {
		return SessionForAdmin{}, newNotFoundError(fmt.Errorf("failed to get session: %w", err))
	}
	return fromSessionToSessionForAdmin(session), nil
}

func (s *Server) GetSession(id string) (Session, error) {
	session, err := s.sessionRepo.Get(id)
	if err != nil {
		return Session{}, newNotFoundError(fmt.Errorf("failed to get session: %w", err))
	}
	return fromSessionToSessionForUser(session), nil
}

type AdminListSessionsResult struct {
	Sessions   []SessionForAdmin `json:"sessions"`
	IsEnd      bool              `json:"is_end"`
	TotalCount int               `json:"total_count"`
}

func (s *Server) AdminListSessions(offset int, pageSize int) (*AdminListSessionsResult, error) {
	listResult, err := s.sessionRepo.List(offset, pageSize)
	if err != nil {
		return nil, newInternalServerError(fmt.Errorf("failed to list sessions: %w", err))
	}

	converted := []SessionForAdmin{}
	for _, session := range listResult.Sessions {
		converted = append(converted, fromSessionToSessionForAdmin(session))
	}

	return &AdminListSessionsResult{
		Sessions:   converted,
		IsEnd:      listResult.IsEnd,
		TotalCount: listResult.TotalCount,
	}, nil
}

type ListSessionsResult struct {
	Sessions   []Session `json:"sessions"`
	IsEnd      bool      `json:"is_end"`
	TotalCount int       `json:"total_count"`
}

func (s *Server) ListSessions(offset int, pageSize int) (*ListSessionsResult, error) {
	listResult, err := s.sessionRepo.List(offset, pageSize)
	if err != nil {
		return nil, newInternalServerError(fmt.Errorf("failed to list sessions: %w", err))
	}

	converted := []Session{}
	for _, session := range listResult.Sessions {
		converted = append(converted, fromSessionToSessionForUser(session))
	}

	return &ListSessionsResult{
		Sessions:   converted,
		IsEnd:      listResult.IsEnd,
		TotalCount: listResult.TotalCount,
	}, nil
}

func (s *Server) AddSession(name string, description string, createdBy string, startsAt time.Time, score int) (string, error) {
	id, err := s.sessionRepo.Add(name, description, createdBy, startsAt, score)
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to add session: %w", err))
	}
	return id, nil
}

func (s *Server) DeleteSession(id string) error {
	if err := s.openSessionRepo.DeleteOpenSession(id); err != nil {
		return newInternalServerError(fmt.Errorf("failed to delete session: %w", err))
	}
	return nil
}

func (s *Server) ApplyAttendanceByFormSubmissions(sessionId string, calledBy string) error {
	dbSession, err := s.sessionRepo.Get(sessionId)
	if err != nil {
		return newNotFoundError(fmt.Errorf("failed to get session: %w", err))
	}

	if !dbSession.CanUpdateMetadata() {
		return newBadRequestError(errors.New("session is already closed"))
	}

	if !dbSession.CanApplyAttendanceByFormSubmissions() {
		return newBadRequestError(errors.New("session cannot apply attendance by form submissions"))
	}

	formSubmissions, err := s.attendanceFormHandler.GetSubmissions(dbSession.GoogleFormId)
	if err != nil {
		return newInternalServerError(fmt.Errorf("failed to get form submissions: %w", err))
	}

	if len(formSubmissions) == 0 {
		if err := s.openSessionRepo.MarkAttendanceIsIgnored(sessionId, "no form submissions"); err != nil {
			return newInternalServerError(fmt.Errorf("failed to mark the session's attendance as ignored: %w", err))
		}
		return nil
	}

	submissionsOnTime := array.Filter(formSubmissions, func(submission attendance.FormSubmission) bool {
		return submission.SubmissionTime.Before(dbSession.StartsAt)
	})

	externalNames := array.Map(submissionsOnTime, func(submission attendance.FormSubmission) string {
		return submission.UserExternalName
	})

	users, err := s.userRepo.GetAllByExternalNames(externalNames)
	if err != nil {
		return newInternalServerError(fmt.Errorf("failed to get users by external names: %w", err))
	}

	externalNameToUserMap := make(map[string]user.User)
	for _, user := range users {
		externalNameToUserMap[user.ExternalName] = user
	}

	notFoundExternalNames := []string{}
	addAttendanceReqs := array.Map(submissionsOnTime, func(submission attendance.FormSubmission) attendance.AddAttendanceReq {
		user, exists := externalNameToUserMap[submission.UserExternalName]
		if !exists {
			notFoundExternalNames = append(notFoundExternalNames, submission.UserExternalName)
		}
		return attendance.AddAttendanceReq{
			SessionId:        sessionId,
			SessionName:      dbSession.Name,
			SessionScore:     dbSession.Score,
			SessionStartedAt: dbSession.StartsAt,
			UserId:           user.Id,
			UserExternalName: user.ExternalName,
			UserGeneration:   user.Generation,
			UserJoinedAt:     submission.SubmissionTime,
			CreatedBy:        calledBy,
		}
	})

	if len(notFoundExternalNames) > 0 {
		if err := s.openSessionRepo.MarkAttendanceIsIgnored(sessionId, fmt.Sprintf("some users (%s) were not found although there are form submissions",
			strings.Join(notFoundExternalNames, ", "))); err != nil {
			return newInternalServerError(fmt.Errorf(
				"some users (%s) were not found although there are form submissions and it has failed to mark the session's attendance as ignored: %w",
				strings.Join(notFoundExternalNames, ", "), err))
		}
		return newInternalServerError(fmt.Errorf("some users (%s) were not found although there are form submissions",
			strings.Join(notFoundExternalNames, ", ")))
	}

	if err := s.attendanceRepo.BulkInsert(addAttendanceReqs); err != nil {
		return newInternalServerError(fmt.Errorf("failed to bulk insert attendance: %w", err))
	}

	if err := s.openSessionRepo.MarkAsAttendanceApplied(sessionId); err != nil {
		return newInternalServerError(fmt.Errorf("failed to close open session: %w", err))
	}

	return nil
}
