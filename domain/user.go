package domain

import (
	"time"

	"github.com/google/uuid"
)

// User struct represents a user in the system
type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user and assigns a new unique ID with name, email and hashed password
func NewUser(name, email, passwordHashed string) *User {
	now := time.Now().UTC()

	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  passwordHashed,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
