package auth

import (
	"errors"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/golang-jwt/jwt/v5"
)

type rushAuth struct {
	secretKey []byte
	clock     clock.Clock
}

func NewRushAuth(secretKey string, clock clock.Clock) *rushAuth {
	return &rushAuth{secretKey: []byte(secretKey), clock: clock}
}

// TODO(#23): Implement the methods.
func (r *rushAuth) SignIn(userIdentifier UserIdentifier) (string, error) {
	rushUserId, ok := userIdentifier.ProviderId(ProviderRush)
	if !ok {
		return "", errors.New("invalid user identifier")
	}

	tokenSpec := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   rushUserId,
		IssuedAt:  jwt.NewNumericDate(r.clock.Now()),
		ExpiresAt: jwt.NewNumericDate(r.clock.Now().Add(7 * 24 * time.Hour)),
	})

	// Can not return an error because the secret key is byte slice and SHA256 is a basic golang hash function.
	// https://github.com/golang-jwt/jwt/blob/v5.2.1/token.go#L63. https://github.com/golang-jwt/jwt/blob/v5.2.1/hmac.go#L83.
	signedToken, _ := tokenSpec.SignedString(r.secretKey)
	return signedToken, nil
}

func (r *rushAuth) GetUserIdentifier(token string) (UserIdentifier, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
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
	subject, err := claims.GetSubject()
	if subject == "" || err != nil {
		return UserIdentifier{}, errors.New("Failed to get information from the token")
	}

	return NewUserIdentifier(map[Provider]string{ProviderRush: subject}, map[Provider]string{ProviderRush: subject}), nil
}
