package auth

import (
	"rush/permission"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
)

func TestNewRushAuth(t *testing.T) {
	t.Run("Returns a new rushAuth instance", func(t *testing.T) {
		mockClock := clock.NewMock()
		auth := NewRushAuth("admin-token", "secret", mockClock)

		assert.Equal(t, &rushAuth{adminToken: "admin-token", secretKey: []byte("secret"), clock: mockClock}, auth)
	})
}

func TestSignInAndVerifyIdentifier(t *testing.T) {
	t.Run("Fails if user ID is empty", func(t *testing.T) {
		rushAuth := NewRushAuth("admin-token", "secret", clock.NewMock())
		token, err := rushAuth.SignIn("", permission.RoleAdmin)

		assert.EqualError(t, err, "user ID is empty")
		assert.Empty(t, token)
	})

	t.Run("Generates a token that has the correct claim and parse it successfully", func(t *testing.T) {
		mockClock := clock.NewMock()
		mockClock.Set(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

		rushAuth := NewRushAuth("admin-token", "secret", mockClock)
		token, err := rushAuth.SignIn("John Doe", permission.RoleAdmin)

		assert.Nil(t, err)

		session, err := rushAuth.GetSession(token)
		assert.Nil(t, err)
		assert.Equal(t, "John Doe", session.Id)
		assert.Equal(t, permission.RoleAdmin, session.Role)
	})
}

func TestGetSession(t *testing.T) {
	t.Run("Returns TokenExpiredError if token is expired", func(t *testing.T) {
		mockClock := clock.NewMock()
		mockClock.Set(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

		rushAuth := NewRushAuth("admin-token", "secret", mockClock)
		token, err := rushAuth.SignIn("John Doe", permission.RoleAdmin)
		assert.Nil(t, err)

		mockClock.Add(14 * 24 * time.Hour)

		session, err := rushAuth.GetSession(token)
		tokenExpiredErr, ok := err.(*TokenExpiredError)
		assert.True(t, ok)
		assert.NotNil(t, tokenExpiredErr.Err)
		assert.Equal(t, Session{}, session)
	})

	t.Run("Returns InvalidTokenError if token is invalid", func(t *testing.T) {
		mockClock := clock.NewMock()
		mockClock.Set(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

		rushAuth := NewRushAuth("admin-token", "secret", mockClock)
		session, err := rushAuth.GetSession("invalid token")

		invalidTokenErr, ok := err.(*InvalidTokenError)
		assert.True(t, ok)
		assert.NotNil(t, invalidTokenErr.Err)
		assert.Equal(t, Session{}, session)
	})

	t.Run("Returns the session if the token is valid", func(t *testing.T) {
		mockClock := clock.NewMock()
		mockClock.Set(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

		rushAuth := NewRushAuth("admin-token", "secret", mockClock)
		token, err := rushAuth.SignIn("John Doe", permission.RoleAdmin)
		assert.Nil(t, err)

		session, err := rushAuth.GetSession(token)
		assert.Nil(t, err)
		assert.Equal(t, "John Doe", session.Id)
		assert.Equal(t, permission.RoleAdmin, session.Role)
		assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Add(7*24*time.Hour).UnixNano(), session.ExpiresAt.UnixNano())
	})

	t.Run("Returns the admin session if the token is the admin token", func(t *testing.T) {
		mockClock := clock.NewMock()
		mockClock.Set(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

		rushAuth := NewRushAuth("admin-token", "secret", mockClock)
		session, err := rushAuth.GetSession("admin-token")
		assert.Nil(t, err)
		assert.Equal(t, "admin-token", session.Id)
		assert.Equal(t, permission.RoleSuperAdmin, session.Role)
	})
}
