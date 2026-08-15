package utils

import (
	"regexp"

	"github.com/google/uuid"
)

// IsNameValid validates the name.
func IsNameValid(name string) bool {
	regex := regexp.MustCompile("^[a-zA-Z]+$")
	return regex.MatchString(name)
}

// IsEmailValid validates the email.
func IsEmailValid(email string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	return regex.MatchString(email)
}

// IsPasswordValid validates the password.
func IsPasswordValid(password string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if len(password) < 6 || len(password) > 20 {
		return false
	}
	if !regex.MatchString(password) {
		return false
	}
	return true
}

// IsUUIDValid validates the UUID.
func IsUUIDValid(id string) bool {
	err := uuid.Validate(id)
	return err == nil
}
