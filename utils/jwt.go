package utils

import (
	"errors"
	"time"
	"user-mgmt/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AccessClaims are the claims in an access token.
type AccessClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a new access token.
func GenerateAccessToken(cfg *config.JWTConfig, userID string, now time.Time) (string, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return "", errors.New("invalid user ID")
	}
	claims := AccessClaims{
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        uuid.NewString(),
			Issuer:    cfg.Issuer,
			Audience:  jwt.ClaimStrings{cfg.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.ExpiresIn)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.Secret))
}

// VerifyAccessToken verifies the access token.
func VerifyAccessToken(cfg *config.JWTConfig, raw string, now time.Time) (*AccessClaims, error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(cfg.Issuer),
		jwt.WithAudience(cfg.Audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(func() time.Time { return now }),
	)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid access token")
	}
	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type")
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return nil, errors.New("invalid token subject")
	}
	return claims, nil
}
