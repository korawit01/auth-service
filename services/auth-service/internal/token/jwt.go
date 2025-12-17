package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Provider exposes token generation/validation.
type Provider interface {
	Generate(userID string, email string) (string, error)
	Verify(tokenStr string) (*Claims, error)
}

// Claims represents the JWT payload.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTProvider struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTProvider(secret string) Provider {
	if secret == "" {
		secret = "dev-secret"
	}
	return &JWTProvider{
		secret: []byte(secret),
		ttl:    24 * time.Hour,
	}
}

func (p *JWTProvider) Generate(userID string, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *JWTProvider) Verify(tokenStr string) (*Claims, error) {
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}

	tkn, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return p.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tkn.Claims.(*Claims)
	if !ok || !tkn.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
