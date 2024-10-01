package server

import (
	"errors"
	"fmt"
	"rush/auth"
	"time"
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

// Returns the user session and the new token if it was refreshed.
func (s *Server) GetUserSession(token string) (UserSession, string, error) {
	session, err := s.authHandler.GetSession(token)
	if err != nil {
		var tokenExpiredError *auth.TokenExpiredError
		if errors.As(err, &tokenExpiredError) {
			return UserSession{}, "", newBadRequestError(fmt.Errorf("token expired: %w", tokenExpiredError))
		}
		var invalidTokenError *auth.InvalidTokenError
		if errors.As(err, &invalidTokenError) {
			return UserSession{}, "", newBadRequestError(fmt.Errorf("invalid token: %w", invalidTokenError))
		}
		return UserSession{}, "", newInternalServerError(fmt.Errorf("failed to get user session: %w", err))
	}

	if session.ExpiresAt.Sub(s.clock.Now()) > 24*time.Hour {
		return UserSession{
			UserId:    session.Id,
			Role:      session.Role,
			ExpiresAt: session.ExpiresAt,
		}, token, nil
	}

	newToken, err := s.authHandler.SignIn(session.Id, session.Role)
	if err != nil {
		return UserSession{}, "", newInternalServerError(fmt.Errorf("failed to refresh token: %w", err))
	}
	return UserSession{
		UserId:    session.Id,
		Role:      session.Role,
		ExpiresAt: session.ExpiresAt,
	}, newToken, nil
}
