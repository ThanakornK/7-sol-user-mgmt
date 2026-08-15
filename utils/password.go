package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword for hash password we use separate method to hash password and other value for security
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return string(hash), err
}

// CheckPassword for check password input with hash
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}
