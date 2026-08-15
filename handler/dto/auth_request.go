package dto

import (
	"errors"
	"strings"
	"user-mgmt/utils"
)

// LoginRequest struct represents the login request in the system
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the login request.
func (r *LoginRequest) Validate() []error {
	var errs []error
	if !utils.IsEmailValid(strings.TrimSpace(r.Email)) {
		errs = append(errs, errors.New("email is invalid"))
	}
	if strings.TrimSpace(r.Password) == "" {
		errs = append(errs, errors.New("password is required"))
	}
	return errs
}
