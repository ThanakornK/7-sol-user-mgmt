// package repository for repository interface and implementation in repository layer
package repository

import (
	"context"
	"time"
	"user-mgmt/domain"
	"user-mgmt/utils"
)

// UserRepository interface for user repository
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserList(ctx context.Context, page, pageSize int) ([]*domain.User, utils.Pagination, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

// UserSessionRepository interface for user session repository
type UserSessionRepository interface {
	Create(ctx context.Context, session *domain.UserSession) error
	GetByID(ctx context.Context, id string) (*domain.UserSession, error)
	Rotate(ctx context.Context, id, currentHash, nextHash string, usedAt time.Time) error
	Revoke(ctx context.Context, id string, revokedAt time.Time) error
	RevokeAllByUserID(ctx context.Context, userID string, revokedAt time.Time) error
}
