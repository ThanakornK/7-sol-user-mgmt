// package dto provides data transfer objects for the application.
package dto

import (
	"errors"
	"user-mgmt/utils"
)

// CreateUserRequest is the request struct for creating a user.
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserRequest is the request struct for updating a user.
type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

// GetUserListRequest is the request struct for getting a list of users.
type GetUserListRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

// Validate CreateUserRequest manual for custom error message
func (r *CreateUserRequest) Validate() []error {
	var errs []error

	// Validate name
	if r.Name == "" {
		errs = append(errs, errors.New("name is required"))
	} else if len(r.Name) < 3 || len(r.Name) > 100 {
		errs = append(errs, errors.New("name must be between 3 and 100 characters"))
	}
	if !utils.IsNameValid(r.Name) {
		errs = append(errs, errors.New("name must contain only letters"))
	}

	// Validate email
	if r.Email == "" {
		errs = append(errs, errors.New("email is required"))
	} else if !utils.IsEmailValid(r.Email) {
		errs = append(errs, errors.New("email is invalid"))
	}

	// Validate password
	if r.Password == "" {
		errs = append(errs, errors.New("password is required"))
	} else if !utils.IsPasswordValid(r.Password) {
		errs = append(errs, errors.New("password must contain only letters and numbers and be between 6 and 20 characters"))
	}
	return errs
}

// Validate UpdateUserRequest manual for custom error message
func (r *UpdateUserRequest) Validate() []error {
	var errs []error
	if r.Name != nil {
		if !utils.IsNameValid(*r.Name) {
			errs = append(errs, errors.New("name must contain only letters"))
		}
	}
	if r.Email != nil {
		if !utils.IsEmailValid(*r.Email) {
			errs = append(errs, errors.New("email is invalid"))
		}
	}
	return errs
}

// Validate GetUserListRequest manual for custom error message
func (r *GetUserListRequest) Validate() []error {
	var errs []error
	if r.Page <= 0 {
		errs = append(errs, errors.New("page must be greater than 0"))
	}
	if r.PageSize <= 0 {
		errs = append(errs, errors.New("pageSize must be greater than 0"))
	}
	return errs
}
