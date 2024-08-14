package server

import (
	"errors"
	"fmt"
	"rush/attendance"
	"rush/session"
	"sort"
	"time"
)

func (s *Server) GetSession(id string) (*Session, error) {
	session, err := s.sessionRepo.Get(id)
	if err != nil {
		return nil, newNotFoundError(fmt.Errorf("failed to get session: %w", err))
	}
	return fromSession(session), nil
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
		converted = append(converted, *fromSession(&session))
	}

	return &ListSessionsResult{
		Sessions:   converted,
		IsEnd:      listResult.IsEnd,
		TotalCount: listResult.TotalCount,
	}, nil
}

func (s *Server) AddSession(name string, description string, startsAt time.Time, score int) (string, error) {
	id, err := s.sessionRepo.Add(name, description, 0, 0, startsAt, score)
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to add session: %w", err))
	}
	return id, nil
}

func (s *Server) CloseSession(sessionId string) error {
	dbSession, err := s.sessionRepo.Get(sessionId)
	if err != nil {
		return newNotFoundError(fmt.Errorf("failed to get session: %w", err))
	}

	if dbSession.IsClosed {
		return newBadRequestError(errors.New("session is already closed"))
	}

	formSubmissions, err := s.attendanceFormHandler.GetSubmissions(dbSession.GoogleFormId)
	if err != nil {
		return newInternalServerError(fmt.Errorf("failed to get form submissions: %w", err))
	}

	submissionsOnTime := []attendance.FormSubmission{}
	for _, submission := range formSubmissions {
		if submission.SubmissionTime.Before(dbSession.StartsAt) {
			submissionsOnTime = append(submissionsOnTime, submission)
		}
	}

	externalIds := []string{}
	for _, submissionOnTime := range submissionsOnTime {
		externalIds = append(externalIds, submissionOnTime.UserExternalId)
	}

	users, err := s.userRepo.GetAllByExternalIds(externalIds)
	if err != nil {
		return newInternalServerError(fmt.Errorf("failed to get users by external IDs: %w", err))
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ExternalId < users[j].ExternalId
	})
	sort.Slice(submissionsOnTime, func(i, j int) bool {
		return submissionsOnTime[i].UserExternalId < submissionsOnTime[j].UserExternalId
	})

	addAttendanceReqs := []attendance.AddAttendanceReq{}
	for index, submissionOnTime := range submissionsOnTime {
		user := users[index]
		addAttendanceReqs = append(addAttendanceReqs, attendance.AddAttendanceReq{
			SessionId:        sessionId,
			SessionName:      dbSession.Name,
			SessionScore:     dbSession.Score,
			SessionStartedAt: dbSession.StartsAt,
			UserId:           user.Id,
			UserName:         user.Name,
			UserGeneration:   user.Generation,
			UserJoinedAt:     submissionOnTime.SubmissionTime,
		})
	}

	if err := s.attendanceRepo.BulkInsert(addAttendanceReqs); err != nil {
		return newInternalServerError(fmt.Errorf("failed to bulk insert attendance: %w", err))
	}

	isClosed := true
	_, err = s.sessionRepo.Update(sessionId, &session.UpdateForm{IsClosed: &isClosed, ReturnUpdatedSession: false})

	if err != nil {
		return newInternalServerError(fmt.Errorf("failed to update session to be closed: %w", err))
	}

	return nil
}
