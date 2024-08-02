package auth

import (
	"errors"
)

type rushAuth struct{}

func NewRushAuth() *rushAuth {
	return &rushAuth{}
}

// TODO(#23): Implement the methods.
func (f *rushAuth) SignIn(userIdentifier UserIdentifier) (string, error) {
	return "temp-token", nil
}

func (f *rushAuth) GetUserIdentifier(token string) (UserIdentifier, error) {
	if token == "temp-token" {
		return UserIdentifier{
			ProviderIds: map[Provider]string{ProviderRush: "temp-user-id"},
			Emails:      map[Provider]string{ProviderRush: "temp-email"},
		}, nil
	}
	return UserIdentifier{}, errors.New("invalid token")
}
