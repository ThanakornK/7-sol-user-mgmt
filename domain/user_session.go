package domain

import "time"

// UserSession struct represents a user session in the system
type UserSession struct {
	ID         string
	UserID     string
	TokenHash  string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	LastUsedAt *time.Time
	RevokedAt  *time.Time
}

// NewUserSession creates a new user session and assigns a new unique ID with user ID, token hash, and expires at
func NewUserSession(id, userID, tokenHash string, expiresIn time.Duration) *UserSession {
	now := time.Now().UTC()
	return &UserSession{
		ID:        id,
		UserID:    userID,
		TokenHash: tokenHash,
		CreatedAt: now,
		ExpiresAt: now.Add(expiresIn)}
}

// IsActive checks if the session is active.
func (s UserSession) IsActive(now time.Time) bool {
	return s.RevokedAt == nil && now.Before(s.ExpiresAt)
}

// TokenPair struct represents a pair of access token and refresh token.
type TokenPair struct {
	AccessToken     string
	RefreshToken    string
	AccessExpiresIn time.Duration
}
