package auth

import (
	"errors"
	"fmt"
	"rush/permission"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/golang-jwt/jwt/v5"
)

type rushAuth struct {
	// The super admin token that passes everything. It's used for developers.
	superAdminToken string
	// The secret key to sign and verify the JWT.
	secretKey []byte
	// The clock to get the current time. It's used to mock the time in tests.
	clock clock.Clock
}

type rushClaims struct {
	jwt.RegisteredClaims
	Role permission.Role `json:"role"`
}

const SuperAdminId = "super-admin-token"

func NewRushAuth(superAdminToken string, secretKey string, clock clock.Clock) *rushAuth {
	return &rushAuth{superAdminToken: superAdminToken, secretKey: []byte(secretKey), clock: clock}
}

func (r *rushAuth) SignIn(userId string, role permission.Role) (string, error) {
	if userId == "" {
		return "", errors.New("user ID is empty")
	}

	tokenSpec := jwt.NewWithClaims(jwt.SigningMethodHS256, rushClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(r.clock.Now()),
			ExpiresAt: jwt.NewNumericDate(r.clock.Now().Add(7 * 24 * time.Hour)),
		},
		Role: role,
	})

	// Can not return an error because the secret key is byte slice and SHA256 is a basic golang hash function.
	// https://github.com/golang-jwt/jwt/blob/v5.2.1/token.go#L63. https://github.com/golang-jwt/jwt/blob/v5.2.1/hmac.go#L83.
	signedToken, _ := tokenSpec.SignedString(r.secretKey)
	return signedToken, nil
}

func (r *rushAuth) GetSession(token string) (Session, error) {
	if token == r.superAdminToken {
		return Session{
			Id:        SuperAdminId,
			Role:      permission.RoleSuperAdmin,
			ExpiresAt: r.clock.Now().Add(100000 * time.Hour),
		}, nil
	}

	parsedToken, err := jwt.ParseWithClaims(token, &rushClaims{}, func(token *jwt.Token) (interface{}, error) {
		return r.secretKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithTimeFunc(func() time.Time {
		// Use clock to get the current time, not the standard "time" package.
		return r.clock.Now()
	}))
	if errors.Is(err, jwt.ErrTokenExpired) {
		return Session{}, &TokenExpiredError{Err: err}
	}
	if err != nil {
		return Session{}, &InvalidTokenError{Err: err}
	}

	claims := parsedToken.Claims
	rushClaims, ok := claims.(*rushClaims)
	if !ok {
		return Session{}, errors.New("Failed to parse the token")
	}
	subject, err := rushClaims.GetSubject()
	if err != nil {
		return Session{}, fmt.Errorf("Failed to get information from the token: %w", err)
	}
	if subject == "" {
		return Session{}, errors.New("The token does not have a subject")
	}

	return Session{
		Id:        subject,
		Role:      rushClaims.GetRole(),
		ExpiresAt: rushClaims.ExpiresAt.Time,
	}, nil
}

func (r *rushClaims) GetRole() permission.Role {
	return r.Role
}
