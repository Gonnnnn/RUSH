package session

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestDeleteOpenSession(t *testing.T) {
	t.Run("Fails if it fails to get session", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionRepo := NewMockSessionRepo(controller)
		service := NewService(sessionRepo)

		sessionRepo.EXPECT().Get("session-id").Return(Session{}, errors.New("failed to get session"))
		err := service.DeleteOpenSession("session-id")

		assert.Equal(t, fmt.Errorf("repo failed to get session: %w", errors.New("failed to get session")), err)
	})

	t.Run("Fails if session is already closed", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionRepo := NewMockSessionRepo(controller)
		service := NewService(sessionRepo)

		sessionRepo.EXPECT().Get("session-id").Return(Session{IsClosed: true}, nil)
		err := service.DeleteOpenSession("session-id")

		assert.Equal(t, errors.New("session is already closed"), err)
	})

	t.Run("Fails if it fails to delete session", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionRepo := NewMockSessionRepo(controller)
		service := NewService(sessionRepo)

		sessionRepo.EXPECT().Get("session-id").Return(Session{IsClosed: false}, nil)
		sessionRepo.EXPECT().Delete("session-id").Return(errors.New("failed to delete session"))
		err := service.DeleteOpenSession("session-id")

		assert.Equal(t, fmt.Errorf("repo failed to delete session: %w", errors.New("failed to delete session")), err)
	})

	t.Run("Success", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionRepo := NewMockSessionRepo(controller)
		service := NewService(sessionRepo)

		sessionRepo.EXPECT().Get("session-id").Return(Session{IsClosed: false}, nil)
		sessionRepo.EXPECT().Delete("session-id").Return(nil)
		err := service.DeleteOpenSession("session-id")

		assert.NoError(t, err)
	})
}

func TestUpdateOpenSession(t *testing.T) {
	t.Run("Fails if it fails to get session", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionRepo := NewMockSessionRepo(controller)
		service := NewService(sessionRepo)

		sessionRepo.EXPECT().Get("session-id").Return(Session{}, errors.New("failed to get session"))
		_, err := service.UpdateOpenSession("session-id", OpenSessionUpdateForm{})

		assert.Equal(t, fmt.Errorf("repo failed to get session: %w", errors.New("failed to get session")), err)
	})

	t.Run("Fails if session is already closed", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionRepo := NewMockSessionRepo(controller)
		service := NewService(sessionRepo)

		sessionRepo.EXPECT().Get("session-id").Return(Session{IsClosed: true}, nil)
		_, err := service.UpdateOpenSession("session-id", OpenSessionUpdateForm{})

		assert.Equal(t, errors.New("session is already closed"), err)
	})

	t.Run("Fails if it fails to update session", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionRepo := NewMockSessionRepo(controller)
		service := NewService(sessionRepo)

		sessionRepo.EXPECT().Get("session-id").Return(Session{IsClosed: false}, nil)
		newTitle := "new-title"
		newDescription := "new-description"
		newStartsAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		newScore := 100
		newGoogleFormId := "google-form-id"
		newGoogleFormUri := "google-form-uri"
		sessionRepo.EXPECT().Update("session-id", UpdateForm{
			Title:         &newTitle,
			Description:   &newDescription,
			StartsAt:      &newStartsAt,
			Score:         &newScore,
			GoogleFormId:  &newGoogleFormId,
			GoogleFormUri: &newGoogleFormUri,
		}).Return(Session{}, errors.New("failed to update session"))
		_, err := service.UpdateOpenSession("session-id", OpenSessionUpdateForm{
			Title:         &newTitle,
			Description:   &newDescription,
			StartsAt:      &newStartsAt,
			Score:         &newScore,
			GoogleFormId:  &newGoogleFormId,
			GoogleFormUri: &newGoogleFormUri,
		})

		assert.Equal(t, fmt.Errorf("repo failed to update session: %w", errors.New("failed to update session")), err)
	})

	t.Run("Success", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionRepo := NewMockSessionRepo(controller)
		service := NewService(sessionRepo)

		sessionRepo.EXPECT().Get("session-id").Return(Session{IsClosed: false}, nil)
		newTitle := "new-title"
		newDescription := "new-description"
		newStartsAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		newScore := 100
		newGoogleFormId := "google-form-id"
		newGoogleFormUri := "google-form-uri"
		sessionRepo.EXPECT().Update("session-id", UpdateForm{
			Title:         &newTitle,
			Description:   &newDescription,
			StartsAt:      &newStartsAt,
			Score:         &newScore,
			GoogleFormId:  &newGoogleFormId,
			GoogleFormUri: &newGoogleFormUri,

			ReturnUpdatedSession: true,
		}).Return(Session{
			Id:            "session-id",
			Name:          newTitle,
			Description:   newDescription,
			StartsAt:      newStartsAt,
			Score:         newScore,
			IsClosed:      false,
			CreatedBy:     "created-by",
			GoogleFormId:  newGoogleFormId,
			GoogleFormUri: newGoogleFormUri,
		}, nil)

		updatedSession, err := service.UpdateOpenSession("session-id", OpenSessionUpdateForm{
			Title:                &newTitle,
			Description:          &newDescription,
			StartsAt:             &newStartsAt,
			Score:                &newScore,
			GoogleFormId:         &newGoogleFormId,
			GoogleFormUri:        &newGoogleFormUri,
			ReturnUpdatedSession: true,
		})

		assert.NoError(t, err)
		assert.Equal(t, newTitle, updatedSession.Name)
		assert.Equal(t, newDescription, updatedSession.Description)
		assert.Equal(t, newStartsAt, updatedSession.StartsAt)
		assert.Equal(t, newScore, updatedSession.Score)
	})
}

func TestCloseOpenSession(t *testing.T) {
	t.Run("Fails if it fails to update session", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionRepo := NewMockSessionRepo(controller)
		service := NewService(sessionRepo)

		isClosed := true
		sessionRepo.EXPECT().Update("session-id", UpdateForm{
			IsClosed: &isClosed,
		}).Return(Session{}, errors.New("failed to update session"))
		err := service.CloseOpenSession("session-id")

		assert.Equal(t, fmt.Errorf("repo failed to update session: %w", errors.New("failed to update session")), err)
	})

	t.Run("Success", func(t *testing.T) {
		controller := gomock.NewController(t)
		sessionRepo := NewMockSessionRepo(controller)
		service := NewService(sessionRepo)

		isClosed := true
		sessionRepo.EXPECT().Update("session-id", UpdateForm{
			IsClosed: &isClosed,
		}).Return(Session{}, nil)

		err := service.CloseOpenSession("session-id")

		assert.NoError(t, err)
	})
}
