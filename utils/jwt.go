// utils package for utils function use in project
package utils

import (
	"errors"
	"time"
	"user-mgmt/config"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims struct for form data token
type UserClaims struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken from user profile struct to token string
func GenerateToken(cfg *config.JWTConfig, id, name, email string) (accessToken, refreshToken string, err error) {

	// Generate Access token
	accessClaims := &UserClaims{
		ID:    id,
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.ExpiresIn)), // 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   id,
		},
	}

	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	// Generate Refresh token
	refreshClaims := &UserClaims{
		ID:    id,
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.RefreshExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   id,
		},
	}
	reToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := reToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// VerifyToken verify user profile token string
func VerifyToken(cfg *config.JWTConfig, tokenString string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is not valid.")
	}

	if claims, ok := token.Claims.(*UserClaims); ok {
		// Success
		return claims, nil
	}

	return nil, errors.New("invalid token claims.")
}
