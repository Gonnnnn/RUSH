package auth

import (
	"fmt"
	"rush/permission"
	"time"
)

// The session information of the user.
// It includes necessary user information for authentication
// and the time when the session will expire.
type Session struct {
	// The ID of the user. E.g., 1234567890
	Id string
	// The role of the user. It is used to determine the access level of the user.
	// E.g., member, admin, etc.
	Role permission.Role
	// The time when the session will expire.
	ExpiresAt time.Time
}

// Error to indicate token has been expired.
// Method doc will specify it if it returns this error.
type TokenExpiredError struct {
	Err error
}

func (e *TokenExpiredError) Error() string {
	if e.Err == nil {
		return "token expired"
	}
	return fmt.Sprintf("token expired: %v", e.Err)
}

// Error to indicate token is invalid.
// Method doc will specify it if it returns this error.
type InvalidTokenError struct {
	Err error
}

func (e *InvalidTokenError) Error() string {
	if e.Err == nil {
		return "invalid token"
	}
	return fmt.Sprintf("invalid token: %v", e.Err)
}
