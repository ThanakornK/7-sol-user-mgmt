// package service for service interfaces and implementation in service layer
package service

import (
	"context"
	"user-mgmt/domain"
	"user-mgmt/utils"
)

// UserService interface for user management
type UserService interface {
	// User Management
	Create(ctx context.Context, email, name, password string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetUserList(ctx context.Context, page, pageSize int) ([]*domain.User, utils.Pagination, error)
	Update(ctx context.Context, id string, email, name *string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

// AuthService interface for authentication
type AuthService interface {
	Login(ctx context.Context, email, password string) (*domain.User, error)
	Logout(ctx context.Context, userID string) error
	RefreshToken(ctx context.Context, userID string) (*domain.User, error)
}
