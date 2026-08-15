package dto

import "github.com/google/uuid"

// UserResponse is the response struct for a user.
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

// UserListResponse is the response struct for a list of users.
type UserListResponse struct {
	Users    []UserResponse `json:"users"`
	Total    int64          `json:"total"`
	Page     int64          `json:"page"`
	PageSize int64          `json:"pageSize"`
}
