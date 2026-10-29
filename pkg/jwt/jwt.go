package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

var ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

type Manager struct {
	secret   string
	duration time.Duration
}

func New(secret string, duration time.Duration) *Manager {
	return &Manager{
		secret:   secret,
		duration: duration,
	}
}

func (m *Manager) GenerateToken() (string, error) {
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtlib.RegisteredClaims{
		Subject:   "admin",
		ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(m.duration)),
	})

	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("jwt - GenerateToken - token.SignedString: %w", err)
	}

	return tokenString, nil
}

func (m *Manager) ParseToken(tokenString string) (string, error) {
	token, err := jwtlib.Parse(tokenString,
		func(token *jwtlib.Token) (any, error) {
			if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
			}

			return []byte(m.secret), nil
		})

	if err != nil {
		return "", fmt.Errorf("jwt - ParseToken - jwtlib.Parse: %w", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("jwt - ParseToken - GetSubject: %w", err)
	}

	return subject, nil
}
