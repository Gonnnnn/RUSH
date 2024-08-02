package auth

import (
	"context"
	"fmt"
	"net/mail"

	fbAuth "firebase.google.com/go/auth"
)

// firebaseAuthClient is an interface for Firebase Auth client. https://pkg.go.dev/firebase.google.com/go/auth#Client
type firebaseAuthClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*fbAuth.Token, error)
}

type firebaseAuth struct {
	// The firebase auth client to verify the token.
	client firebaseAuthClient
}

func NewFbAuth(client firebaseAuthClient) *firebaseAuth {
	return &firebaseAuth{
		client: client,
	}
}

func (f *firebaseAuth) GetUserIdentifier(token string) (UserIdentifier, error) {
	decodedToken, err := f.client.VerifyIDToken(context.Background(), token)
	if err != nil {
		return UserIdentifier{}, fmt.Errorf("failed to verify the token: %w", err)
	}

	email := decodedToken.Claims["email"]
	if email == nil {
		return UserIdentifier{}, fmt.Errorf("failed to verify the token: invalid email in claim")
	}

	emailStr, ok := email.(string)
	if !ok {
		return UserIdentifier{}, fmt.Errorf("failed to verify the token: invalid email in claim")
	}

	_, err = mail.ParseAddress(emailStr)
	if err != nil {
		return UserIdentifier{}, fmt.Errorf("failed to verify the token: invalid email format")
	}

	firestoreId := decodedToken.Claims["user_id"]
	if firestoreId == nil {
		return UserIdentifier{}, fmt.Errorf("failed to verify the token: invalid user_id in claim")
	}

	firestoreIdStr, ok := firestoreId.(string)
	if !ok {
		return UserIdentifier{}, fmt.Errorf("failed to verify the token: invalid user_id in claim")
	}

	return UserIdentifier{
		ProviderIds: map[Provider]string{ProviderFirebase: firestoreIdStr},
		Emails:      map[Provider]string{ProviderFirebase: emailStr},
	}, nil
}
