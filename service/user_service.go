package service

import (
	"context"
	"errors"
	"user-mgmt/domain"
	"user-mgmt/repository"
	"user-mgmt/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// userService implements UserService interface
type userService struct {
	userRepository repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{userRepository: userRepository}
}

// Create creates a new user from email, name and raw password
func (s *userService) Create(ctx context.Context, email, name, password string) (*domain.User, error) {
	// Check duplicate email exists
	_, err := s.userRepository.GetByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrEmailExists
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	// Hash password
	passwordHashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := domain.NewUser(name, email, passwordHashed)
	return s.userRepository.Create(ctx, user)
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepository.GetByID(ctx, id)
}

// GetUserList retrieves a list of users with pagination
func (s *userService) GetUserList(ctx context.Context, page, pageSize int) ([]*domain.User, utils.Pagination, error) {
	return s.userRepository.GetUserList(ctx, page, pageSize)
}

// Update updates a user's email and name
func (s *userService) Update(ctx context.Context, id string, name *string, email *string) (*domain.User, error) {
	// Check existing user
	existingUser, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update info
	if email != nil {
		existingUser.Email = *email
	}
	if name != nil {
		existingUser.Name = *name
	}

	return s.userRepository.Update(ctx, existingUser)
}

// Delete deletes a user by ID
func (s *userService) Delete(ctx context.Context, id string) error {
	// Check existing user
	_, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.userRepository.Delete(ctx, id)
}
