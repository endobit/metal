package auth

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenUsername = "username"
	tokenAdmin    = "admin"
)

type tokenManager struct {
	key      []byte
	tokenTTL time.Duration
}

func newTokenManager(tokenTTL time.Duration) *tokenManager {
	key := make([]byte, 64)
	_, _ = rand.Read(key) // panics on error

	return &tokenManager{
		key:      key,
		tokenTTL: tokenTTL,
	}
}

func (m *tokenManager) Generate(user string, admin bool) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": now.Add(m.tokenTTL).Unix(), // Expires At
		"iat": now.Unix(),                 // Issued At
		"iss": "stackd",                   // Issuer
		"sub": "access",                   // Subject

		tokenUsername: user,
		tokenAdmin:    admin,
	})

	return token.SignedString(m.key)
}

func (m *tokenManager) Verify(t string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(t, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != "HS256" {
			return nil, errors.New("unexpected signing method")
		}

		return m.key, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
