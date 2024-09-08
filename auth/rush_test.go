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
		auth := NewRushAuth("secret", mockClock)

		assert.Equal(t, &rushAuth{secretKey: []byte("secret"), clock: mockClock}, auth)
	})
}

func TestSignInAndVerifyIdentifier(t *testing.T) {

	t.Run("Fails if it can not find the rush user id", func(t *testing.T) {
		rushAuth := NewRushAuth("secret", clock.NewMock())
		token, err := rushAuth.SignIn(
			NewUserIdentifier(
				map[Provider]string{ProviderFirebase: "John Doe"},
				nil, /* =emails */
				map[Provider]permission.Role{ProviderFirebase: permission.RoleAdmin},
			),
		)

		assert.EqualError(t, err, "invalid user identifier")
		assert.Empty(t, token)
	})

	t.Run("Generates a token that has the correct claim and parse it successfully", func(t *testing.T) {
		mockClock := clock.NewMock()
		mockClock.Set(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

		rushAuth := NewRushAuth("secret", mockClock)
		token, err := rushAuth.SignIn(
			NewUserIdentifier(
				map[Provider]string{ProviderRush: "John Doe"},
				nil, /* =emails */
				map[Provider]permission.Role{ProviderRush: permission.RoleAdmin},
			),
		)

		assert.Nil(t, err)

		userIdentifier, err := rushAuth.GetUserIdentifier(token)
		assert.Nil(t, err)

		rushUserId, ok := userIdentifier.ProviderId(ProviderRush)
		assert.True(t, ok)
		assert.Equal(t, "John Doe", rushUserId)
	})
}
