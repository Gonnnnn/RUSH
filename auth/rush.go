package auth

import (
	"errors"
	"rush/permission"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/golang-jwt/jwt/v5"
)

type rushAuth struct {
	// The secret key to sign and verify the JWT.
	secretKey []byte
	// The clock to get the current time. It's used to mock the time in tests.
	clock clock.Clock
}

type rushClaims struct {
	jwt.RegisteredClaims
	Role permission.Role `json:"role"`
}

func NewRushAuth(secretKey string, clock clock.Clock) *rushAuth {
	return &rushAuth{secretKey: []byte(secretKey), clock: clock}
}

func (r *rushAuth) SignIn(userIdentifier UserIdentifier) (string, error) {
	rushUserId, ok := userIdentifier.ProviderId(ProviderRush)
	if !ok {
		return "", errors.New("invalid user identifier")
	}

	tokenSpec := jwt.NewWithClaims(jwt.SigningMethodHS256, rushClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   rushUserId,
			IssuedAt:  jwt.NewNumericDate(r.clock.Now()),
			ExpiresAt: jwt.NewNumericDate(r.clock.Now().Add(7 * 24 * time.Hour)),
		},
		Role: userIdentifier.RushRole(),
	})

	// Can not return an error because the secret key is byte slice and SHA256 is a basic golang hash function.
	// https://github.com/golang-jwt/jwt/blob/v5.2.1/token.go#L63. https://github.com/golang-jwt/jwt/blob/v5.2.1/hmac.go#L83.
	signedToken, _ := tokenSpec.SignedString(r.secretKey)
	return signedToken, nil
}

func (r *rushAuth) GetUserIdentifier(token string) (UserIdentifier, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &rushClaims{}, func(token *jwt.Token) (interface{}, error) {
		return r.secretKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithTimeFunc(func() time.Time {
		// Use clock to get the current time, not the standard "time" package.
		return r.clock.Now()
	}))
	if errors.Is(err, jwt.ErrTokenExpired) {
		return UserIdentifier{}, errors.New("The token has expired")
	}
	if err != nil {
		return UserIdentifier{}, errors.New("Failed to parse the token")
	}

	claims := parsedToken.Claims
	rushClaims, ok := claims.(*rushClaims)
	if !ok {
		return UserIdentifier{}, errors.New("Failed to parse the token")
	}
	subject, err := rushClaims.GetSubject()
	if subject == "" || err != nil {
		return UserIdentifier{}, errors.New("Failed to get information from the token")
	}

	return NewUserIdentifier(
		map[Provider]string{ProviderRush: subject},
		nil,
		map[Provider]permission.Role{ProviderRush: rushClaims.GetRole()},
	), nil
}

func (r *rushClaims) GetRole() permission.Role {
	return r.Role
}
