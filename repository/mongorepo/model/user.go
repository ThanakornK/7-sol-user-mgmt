package model

import (
	"time"
	"user-mgmt/domain"

	"github.com/google/uuid"
)

// User struct represents a user in the MongoDB database.
type User struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"-" bson:"password"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// FromUserDomain converts a domain.User to a model.User.
func FromUserDomain(user *domain.User) *User {
	return &User{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToDomain converts a model.User to a domain.User.
func (u *User) ToDomain() *domain.User {
	return &domain.User{
		ID:        uuid.MustParse(u.ID),
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
