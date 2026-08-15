package domain

import "errors"

// centralized errors.
var (
	ErrEmailExists         = errors.New("email already exists.")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
