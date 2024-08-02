package auth

import (
	"context"
	"fmt"
	"testing"

	fbAuth "firebase.google.com/go/auth"
	"github.com/stretchr/testify/assert"
)

type mockFbAuthClient struct {
	verifyIDTokenToken *fbAuth.Token
	verifyIDTokenError error
}

func (m *mockFbAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*fbAuth.Token, error) {
	return m.verifyIDTokenToken, m.verifyIDTokenError
}

func TestGetUserIdentifier(t *testing.T) {
	t.Run("Fails if firebase auth client fails to verify the token", func(t *testing.T) {
		mockFbAuthClient := &mockFbAuthClient{verifyIDTokenToken: nil, verifyIDTokenError: fmt.Errorf("mock error")}

		fbAuth := NewFbAuth(mockFbAuthClient)
		identifier, err := fbAuth.GetUserIdentifier("token")

		assert.Empty(t, identifier)
		assert.EqualError(t, err, "failed to verify the token: mock error")
	})

	t.Run("Invalid claim", func(t *testing.T) {
		t.Run("Fails if email is missing in the token", func(t *testing.T) {
			mockFbAuthClient := &mockFbAuthClient{verifyIDTokenToken: &fbAuth.Token{}, verifyIDTokenError: nil}

			fbAuth := NewFbAuth(mockFbAuthClient)
			identifier, err := fbAuth.GetUserIdentifier("token")

			assert.Empty(t, identifier)
			assert.EqualError(t, err, "failed to verify the token: invalid email in claim")
		})

		t.Run("Fails if email is not a string", func(t *testing.T) {
			mockFbAuthClient := &mockFbAuthClient{verifyIDTokenToken: &fbAuth.Token{Claims: map[string]interface{}{"email": 1}}, verifyIDTokenError: nil}

			fbAuth := NewFbAuth(mockFbAuthClient)
			identifier, err := fbAuth.GetUserIdentifier("token")

			assert.Empty(t, identifier)
			assert.EqualError(t, err, "failed to verify the token: invalid email in claim")
		})

		t.Run("Fails if user_id is missing in the token", func(t *testing.T) {
			mockFbAuthClient := &mockFbAuthClient{verifyIDTokenToken: &fbAuth.Token{Claims: map[string]interface{}{"email": "john.doe@gmail.conm"}}, verifyIDTokenError: nil}

			fbAuth := NewFbAuth(mockFbAuthClient)
			identifier, err := fbAuth.GetUserIdentifier("token")

			assert.Empty(t, identifier)
			assert.EqualError(t, err, "failed to verify the token: invalid user_id in claim")
		})

		t.Run("Fails if user_id is not a string", func(t *testing.T) {
			mockFbAuthClient := &mockFbAuthClient{verifyIDTokenToken: &fbAuth.Token{Claims: map[string]interface{}{"email": "john.doe@gmail.conm", "user_id": 1}}, verifyIDTokenError: nil}

			fbAuth := NewFbAuth(mockFbAuthClient)
			identifier, err := fbAuth.GetUserIdentifier("token")

			assert.Empty(t, identifier)
			assert.EqualError(t, err, "failed to verify the token: invalid user_id in claim")
		})
	})

	t.Run("Successfully fetches the firebase identifier of the user", func(t *testing.T) {
		mockFbAuthClient := &mockFbAuthClient{verifyIDTokenToken: &fbAuth.Token{Claims: map[string]interface{}{"email": "john.doe@gmail.com", "user_id": "abcdefg"}}, verifyIDTokenError: nil}

		fbAuth := NewFbAuth(mockFbAuthClient)
		identifier, err := fbAuth.GetUserIdentifier("token")

		providerId, ok := identifier.ProviderId(ProviderFirebase)
		assert.True(t, ok)
		assert.Equal(t, "abcdefg", providerId)
		email, ok := identifier.Email(ProviderFirebase)
		assert.True(t, ok)
		assert.Equal(t, "john.doe@gmail.com", email)
		assert.Nil(t, err)
	})
}
