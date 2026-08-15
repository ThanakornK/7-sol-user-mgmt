package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, email, passwordHashed string) *User {
	now := time.Now().UTC()

	return &User{
		Name:      name,
		Email:     email,
		Password:  passwordHashed,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
