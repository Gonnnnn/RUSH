package server

import (
	"errors"
	"fmt"
	"rush/auth"
)

func (s *Server) SignIn(token string) (string, error) {
	userIdentifier, err := s.tokenInspector.GetUserIdentifier(token)
	if err != nil {
		return "", newBadRequestError(fmt.Errorf("failed to get user identifier: %w", err))
	}

	email, ok := userIdentifier.Email(s.tokenInspector.Provider())
	if !ok {
		return "", newInternalServerError(errors.New("failed to get email from user identifier although there should be"))
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", newNotFoundError(fmt.Errorf("failed to get user by email (%s): %w", email, err))
	}

	rushToken, err := s.authHandler.SignIn(
		auth.NewUserIdentifier(
			map[auth.Provider]string{auth.ProviderRush: user.Id},
			map[auth.Provider]string{auth.ProviderRush: email},
		),
	)
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to sign in: %w", err))
	}

	return rushToken, nil
}

func (s *Server) IsTokenValid(token string) bool {
	if _, err := s.authHandler.GetUserIdentifier(token); err != nil {
		return false
	}
	return true
}

func (s *Server) GetUserIdentifier(token string) (string, error) {
	userIdentifier, err := s.authHandler.GetUserIdentifier(token)
	if err != nil {
		return "", newBadRequestError(fmt.Errorf("failed to get user identifier: %w", err))
	}

	userId, ok := userIdentifier.ProviderId(auth.ProviderRush)
	if !ok {
		return "", newInternalServerError(errors.New("failed to get user ID from user identifier although there should be"))
	}

	return userId, nil
}
