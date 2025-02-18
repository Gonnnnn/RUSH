package server

import (
	"errors"
	"fmt"
	"rush/auth"
	"rush/permission"
	"rush/user"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestSignIn(t *testing.T) {
	t.Run("Returns bad request error if failed to get user identifier", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockOauthClient := NewMockoauthClient(controller)
		server := New(mockOauthClient, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		mockOauthClient.EXPECT().GetEmail("token").Return("", assert.AnError)
		token, err := server.SignIn("token")

		assert.Equal(t, "", token)
		assert.Equal(t, newBadRequestError(fmt.Errorf("failed to get user identifier: %w", assert.AnError)), err)
	})

	t.Run("Returns not found error if user is not found", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockOauthClient := NewMockoauthClient(controller)
		mockUserRepo := NewMockuserRepo(controller)
		server := New(mockOauthClient, nil, mockUserRepo, nil, nil, nil, nil, nil, nil, nil, nil)

		mockOauthClient.EXPECT().GetEmail("token").Return("email@example.com", nil)
		mockUserRepo.EXPECT().GetByEmail("email@example.com").Return(nil, user.ErrNotFound)
		token, err := server.SignIn("token")

		assert.Equal(t, "", token)
		assert.Equal(t, newNotFoundError(fmt.Errorf("failed to get user by email (%s): %w", "email@example.com", user.ErrNotFound)), err)
	})

	t.Run("Returns internal server error if failed to get user", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockOauthClient := NewMockoauthClient(controller)
		mockUserRepo := NewMockuserRepo(controller)
		server := New(mockOauthClient, nil, mockUserRepo, nil, nil, nil, nil, nil, nil, nil, nil)

		mockOauthClient.EXPECT().GetEmail("token").Return("email@example.com", nil)
		mockUserRepo.EXPECT().GetByEmail("email@example.com").Return(nil, assert.AnError)
		token, err := server.SignIn("token")

		assert.Equal(t, "", token)
		assert.Equal(t, newInternalServerError(fmt.Errorf("failed to get user by email (%s): %w", "email@example.com", assert.AnError)), err)
	})

	t.Run("Returns internal server error if failed to sign in", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockOauthClient := NewMockoauthClient(controller)
		mockUserRepo := NewMockuserRepo(controller)
		mockAuthHandler := NewMockauthHandler(controller)
		server := New(mockOauthClient, mockAuthHandler, mockUserRepo, nil, nil, nil, nil, nil, nil, nil, nil)

		mockOauthClient.EXPECT().GetEmail("token").Return("email@example.com", nil)
		mockUserRepo.EXPECT().GetByEmail("email@example.com").Return(&user.User{
			Id:   "user_id",
			Role: permission.RoleMember,
		}, nil)
		mockAuthHandler.EXPECT().SignIn("user_id", permission.RoleMember).Return("", assert.AnError)
		token, err := server.SignIn("token")

		assert.Equal(t, "", token)
		assert.Equal(t, newInternalServerError(fmt.Errorf("failed to sign in: %w", assert.AnError)), err)
	})

	t.Run("Returns rush token if user is found and signed in", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockOauthClient := NewMockoauthClient(controller)
		mockUserRepo := NewMockuserRepo(controller)
		mockAuthHandler := NewMockauthHandler(controller)
		server := New(mockOauthClient, mockAuthHandler, mockUserRepo, nil, nil, nil, nil, nil, nil, nil, nil)

		mockOauthClient.EXPECT().GetEmail("token").Return("email@example.com", nil)
		mockUserRepo.EXPECT().GetByEmail("email@example.com").Return(&user.User{
			Id:   "user_id",
			Role: permission.RoleMember,
		}, nil)
		mockAuthHandler.EXPECT().SignIn("user_id", permission.RoleMember).Return("rush_token", nil)
		token, err := server.SignIn("token")

		assert.Equal(t, "rush_token", token)
		assert.Nil(t, err)
	})
}

func TestGetUserSession(t *testing.T) {
	t.Run("Returns bad request error if auth handler returns token expired error", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockAuthHandler := NewMockauthHandler(controller)
		server := New(nil, mockAuthHandler, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		mockAuthHandler.EXPECT().GetSession("token").Return(auth.Session{}, &auth.TokenExpiredError{})
		userSession, newToken, err := server.GetUserSession("token")

		assert.Equal(t, UserSession{}, userSession)
		assert.Equal(t, "", newToken)
		assert.Equal(t, newBadRequestError(fmt.Errorf("token expired: %w", &auth.TokenExpiredError{})), err)
	})

	t.Run("Returns bad request error if auth handler returns invalid token error", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockAuthHandler := NewMockauthHandler(controller)
		server := New(nil, mockAuthHandler, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		mockAuthHandler.EXPECT().GetSession("token").Return(auth.Session{}, &auth.InvalidTokenError{})
		userSession, newToken, err := server.GetUserSession("token")

		assert.Equal(t, UserSession{}, userSession)
		assert.Equal(t, "", newToken)
		assert.Equal(t, newBadRequestError(fmt.Errorf("invalid token: %w", &auth.InvalidTokenError{})), err)
	})

	t.Run("Returns internal server error if auth handler returns other error", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockAuthHandler := NewMockauthHandler(controller)
		server := New(nil, mockAuthHandler, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		mockAuthHandler.EXPECT().GetSession("token").Return(auth.Session{}, errors.New("unknown error"))
		userSession, newToken, err := server.GetUserSession("token")

		assert.Equal(t, UserSession{}, userSession)
		assert.Equal(t, "", newToken)
		assert.Equal(t, newInternalServerError(fmt.Errorf("failed to get user session: %w", errors.New("unknown error"))), err)
	})

	t.Run("Refreshes if the session is going to be expired within 24 hours", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockAuthHandler := NewMockauthHandler(controller)
		mockClock := clock.NewMock()
		server := New(nil, mockAuthHandler, nil, nil, nil, nil, nil, nil, nil, nil, mockClock)

		mockClock.Set(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC))
		mockAuthHandler.EXPECT().GetSession("token").Return(auth.Session{
			Id:        "user_id",
			Role:      permission.RoleMember,
			ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)
		mockAuthHandler.EXPECT().SignIn("user_id", permission.RoleMember).Return("new_token", nil)
		userSession, newToken, err := server.GetUserSession("token")

		assert.Equal(t, UserSession{
			UserId:    "user_id",
			Role:      permission.RoleMember,
			ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}, userSession)
		assert.Equal(t, "new_token", newToken)
		assert.Nil(t, err)
	})

	t.Run("Returns error if failed to refresh token", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockAuthHandler := NewMockauthHandler(controller)
		mockClock := clock.NewMock()
		server := New(nil, mockAuthHandler, nil, nil, nil, nil, nil, nil, nil, nil, mockClock)

		mockClock.Set(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC))
		mockAuthHandler.EXPECT().GetSession("token").Return(auth.Session{
			Id:        "user_id",
			Role:      permission.RoleMember,
			ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)
		mockAuthHandler.EXPECT().SignIn("user_id", permission.RoleMember).Return("", errors.New("unknown error"))
		userSession, newToken, err := server.GetUserSession("token")

		assert.Equal(t, UserSession{}, userSession)
		assert.Equal(t, "", newToken)
		assert.Equal(t, newInternalServerError(fmt.Errorf("failed to refresh token: %w", errors.New("unknown error"))), err)
	})

	t.Run("Returns user session without a new token if the session is not going to be expired within 24 hours", func(t *testing.T) {
		controller := gomock.NewController(t)
		mockAuthHandler := NewMockauthHandler(controller)
		mockClock := clock.NewMock()
		server := New(nil, mockAuthHandler, nil, nil, nil, nil, nil, nil, nil, nil, mockClock)

		mockClock.Set(time.Date(2023, 12, 30, 23, 59, 59, 0, time.UTC))
		mockAuthHandler.EXPECT().GetSession("token").Return(auth.Session{
			Id:        "user_id",
			Role:      permission.RoleMember,
			ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)
		userSession, newToken, err := server.GetUserSession("token")

		assert.Equal(t, UserSession{
			UserId:    "user_id",
			Role:      permission.RoleMember,
			ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		}, userSession)
		assert.Equal(t, "token", newToken)
		assert.Nil(t, err)
	})
}
