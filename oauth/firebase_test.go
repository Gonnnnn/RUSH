package oauth

import (
	"context"
	"fmt"
	"testing"

	fbClient "firebase.google.com/go/auth"
	"github.com/stretchr/testify/assert"
)

type mockFbAuthClient struct {
	verifyIDTokenToken *fbClient.Token
	verifyIDTokenError error
}

func (m *mockFbAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*fbClient.Token, error) {
	return m.verifyIDTokenToken, m.verifyIDTokenError
}

func TestGetEmail(t *testing.T) {
	t.Run("Fails if firebase auth client fails to verify the token", func(t *testing.T) {
		mockFbAuthClient := &mockFbAuthClient{verifyIDTokenToken: nil, verifyIDTokenError: fmt.Errorf("mock error")}

		fbClient := NewFbClient(mockFbAuthClient)
		email, err := fbClient.GetEmail("token")

		assert.Empty(t, email)
		assert.EqualError(t, err, "failed to verify the token: mock error")
	})

	t.Run("Invalid claim", func(t *testing.T) {
		t.Run("Fails if email is missing in the token", func(t *testing.T) {
			mockFbAuthClient := &mockFbAuthClient{verifyIDTokenToken: &fbClient.Token{}, verifyIDTokenError: nil}

			fbClient := NewFbClient(mockFbAuthClient)
			email, err := fbClient.GetEmail("token")

			assert.Empty(t, email)
			assert.EqualError(t, err, "failed to verify the token: invalid email in claim")
		})

		t.Run("Fails if email is not a string", func(t *testing.T) {
			mockFbAuthClient := &mockFbAuthClient{verifyIDTokenToken: &fbClient.Token{Claims: map[string]interface{}{"email": 1}}, verifyIDTokenError: nil}

			fbClient := NewFbClient(mockFbAuthClient)
			email, err := fbClient.GetEmail("token")

			assert.Empty(t, email)
			assert.EqualError(t, err, "failed to verify the token: invalid email in claim")
		})
	})

	t.Run("Successfully fetches the firebase email of the user", func(t *testing.T) {
		mockFbAuthClient := &mockFbAuthClient{verifyIDTokenToken: &fbClient.Token{Claims: map[string]interface{}{"email": "john.doe@gmail.com", "user_id": "abcdefg"}}, verifyIDTokenError: nil}

		fbClient := NewFbClient(mockFbAuthClient)
		email, err := fbClient.GetEmail("token")

		assert.Equal(t, "john.doe@gmail.com", email)
		assert.Nil(t, err)
	})
}
