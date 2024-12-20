package oauth

import (
	"context"
	"fmt"
	"log"
	"net/mail"

	fbAuth "firebase.google.com/go/auth"
)

// firebaseAuthClient is an interface for Firebase Auth client. https://pkg.go.dev/firebase.google.com/go/auth#Client
type firebaseAuthClient interface {
	// Verfies the OpenID token and returns the information about the token such as payload.
	VerifyIDToken(ctx context.Context, idToken string) (*fbAuth.Token, error)
}

type firebaseOauth struct {
	// The firebase auth client to verify the token.
	client firebaseAuthClient
}

func NewFbClient(client firebaseAuthClient) *firebaseOauth {
	return &firebaseOauth{
		client: client,
	}
}

func (f *firebaseOauth) GetEmail(token string) (string, error) {
	decodedToken, err := f.client.VerifyIDToken(context.Background(), token)
	if err != nil {
		return "", fmt.Errorf("failed to verify the token: %w", err)
	}

	log.Printf("decodedToken: %+v", decodedToken)

	email := decodedToken.Claims["email"]
	if email == nil {
		return "", fmt.Errorf("failed to verify the token: invalid email in claim")
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", fmt.Errorf("failed to verify the token: invalid email in claim")
	}

	_, err = mail.ParseAddress(emailStr)
	if err != nil {
		return "", fmt.Errorf("failed to verify the token: invalid email format")
	}

	return emailStr, nil
}
