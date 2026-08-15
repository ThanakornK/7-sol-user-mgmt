package utils

import (
	"regexp"

	"github.com/google/uuid"
)

func IsNameValid(name string) bool {
	regex := regexp.MustCompile("^[a-zA-Z]+$")
	return regex.MatchString(name)
}

func IsEmailValid(email string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	return regex.MatchString(email)
}

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

func IsUUIDValid(id string) bool {
	err := uuid.Validate(id)
	return err == nil
}
