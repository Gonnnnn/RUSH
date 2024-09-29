package server

import (
	"errors"
	"fmt"
	"rush/auth"
	"rush/permission"
)

func (s *Server) SignIn(token string) (string, error) {
	email, err := s.oauthClient.GetEmail(token)
	if err != nil {
		return "", newBadRequestError(fmt.Errorf("failed to get user identifier: %w", err))
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", newNotFoundError(fmt.Errorf("failed to get user by email (%s): %w", email, err))
	}

	rushToken, err := s.authHandler.SignIn(user.Id, user.Role)
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to sign in: %w", err))
	}

	return rushToken, nil
}

func (s *Server) IsTokenValid(token string) bool {
	if _, err := s.authHandler.GetSession(token); err != nil {
		return false
	}
	return true
}

func (s *Server) GetUserIdentifier(token string) (string, permission.Role, error) {
	session, err := s.authHandler.GetSession(token)
	if err != nil {
		if errors.Is(err, &auth.TokenExpiredError{}) {
			// TODO(#105): Refresh the token based on the time left.
			return "", permission.RoleNotSpecified, newBadRequestError(fmt.Errorf("token expired: %w", err))
		}
		if errors.Is(err, &auth.InvalidTokenError{}) {
			return "", permission.RoleNotSpecified, newBadRequestError(fmt.Errorf("invalid token: %w", err))
		}
		return "", permission.RoleNotSpecified, newBadRequestError(fmt.Errorf("failed to get user session: %w", err))
	}

	return session.Id, session.Role, nil
}
