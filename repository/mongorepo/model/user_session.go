package model

import (
	"time"
	"user-mgmt/domain"
)

// UserSession struct represents a user session in the MongoDB database.
type UserSession struct {
	ID         string     `bson:"_id"`
	UserID     string     `bson:"userId"`
	TokenHash  string     `bson:"tokenHash"`
	ExpiresAt  time.Time  `bson:"expiresAt"`
	CreatedAt  time.Time  `bson:"createdAt"`
	LastUsedAt *time.Time `bson:"lastUsedAt,omitempty"`
	RevokedAt  *time.Time `bson:"revokedAt,omitempty"`
}

// FromUserSessionDomain converts a domain.UserSession to a model.UserSession.
func FromUserSessionDomain(session *domain.UserSession) *UserSession {
	return &UserSession{ID: session.ID, UserID: session.UserID, TokenHash: session.TokenHash, ExpiresAt: session.ExpiresAt, CreatedAt: session.CreatedAt, LastUsedAt: session.LastUsedAt, RevokedAt: session.RevokedAt}
}

// ToDomain converts a model.UserSession to a domain.UserSession.
func (s *UserSession) ToDomain() *domain.UserSession {
	return &domain.UserSession{ID: s.ID, UserID: s.UserID, TokenHash: s.TokenHash, ExpiresAt: s.ExpiresAt, CreatedAt: s.CreatedAt, LastUsedAt: s.LastUsedAt, RevokedAt: s.RevokedAt}
}
