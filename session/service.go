package session

import (
	"errors"
	"fmt"
	"time"
)

type service struct {
	sessionRepo SessionRepo
}

func NewService(sessionRepo SessionRepo) *service {
	return &service{
		sessionRepo: sessionRepo,
	}
}

func (s *service) DeleteOpenSession(id string) error {
	session, err := s.sessionRepo.Get(id)
	if err != nil {
		return fmt.Errorf("repo failed to get session: %w", err)
	}
	if !session.CanUpdateMetadata() {
		return errors.New("session is already closed")
	}

	if err := s.sessionRepo.Delete(id); err != nil {
		return fmt.Errorf("repo failed to delete session: %w", err)
	}
	return nil
}

type OpenSessionUpdateForm struct {
	Title       *string
	Description *string
	StartsAt    *time.Time
	Score       *int

	GoogleFormId  *string
	GoogleFormUri *string

	ReturnUpdatedSession bool
}

func (s *service) UpdateOpenSession(id string, updateForm OpenSessionUpdateForm) (Session, error) {
	session, err := s.sessionRepo.Get(id)
	if err != nil {
		return Session{}, fmt.Errorf("repo failed to get session: %w", err)
	}
	if !session.CanUpdateMetadata() {
		return Session{}, errors.New("session is already closed")
	}

	updatedSession, err := s.sessionRepo.Update(id,
		UpdateForm{
			Title:         updateForm.Title,
			Description:   updateForm.Description,
			StartsAt:      updateForm.StartsAt,
			Score:         updateForm.Score,
			GoogleFormId:  updateForm.GoogleFormId,
			GoogleFormUri: updateForm.GoogleFormUri,

			ReturnUpdatedSession: updateForm.ReturnUpdatedSession,
		})
	if err != nil {
		return Session{}, fmt.Errorf("repo failed to update session: %w", err)
	}
	return updatedSession, nil
}

func (s *service) CloseOpenSession(id string) error {
	attendanceStatus := AttendanceStatusApplied
	_, err := s.sessionRepo.Update(id, UpdateForm{AttendanceStatus: &attendanceStatus})
	if err != nil {
		return fmt.Errorf("repo failed to update session: %w", err)
	}
	return nil
}

//go:generate mockgen -source=service.go -destination=service_mock.go -package=session
type SessionRepo interface {
	Get(id string) (Session, error)
	Update(id string, updateForm UpdateForm) (Session, error)
	Delete(id string) error
}
