package handler

import (
	"strings"
	"user-mgmt/handler/dto"
	"user-mgmt/service"
	"user-mgmt/utils"

	"github.com/gin-gonic/gin"
)

// userHandler struct implements the UserHandler interface.
type userHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{userService: userService}
}

// CreateUser creates a new user.
func (h *userHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponseStruct("invalid request body.", err.Error()))
		return
	}

	if errs := req.Validate(); len(errs) > 0 {
		c.JSON(400, utils.NewErrorResponseStruct("validation failed.", utils.ErrorMessages(errs)))
		return
	}

	user, err := h.userService.Create(c, req.Email, req.Name, req.Password)
	if err != nil {
		c.JSON(500, utils.NewErrorResponseStruct("failed to create user.", err.Error()))
		return
	}
	c.JSON(201, utils.NewSuccessResponseStruct("user created successfully.", dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}))
}

// GetUserByID retrieves a user by ID from the database.
func (h *userHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(400, utils.NewErrorResponseStruct("validation failed.", "id is required."))
		return
	}

	if !utils.IsUUIDValid(id) {
		c.JSON(400, utils.NewErrorResponseStruct("validation failed.", "id is not a valid uuid."))
		return
	}

	user, err := h.userService.GetByID(c, id)
	if err != nil {
		c.JSON(500, utils.NewErrorResponseStruct("failed to retrieve user.", err.Error()))
		return
	}
	c.JSON(200, utils.NewSuccessResponseStruct("user retrieved successfully.", dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}))
}

// GetUserList retrieves a list of users from the database with pagination.
func (h *userHandler) GetUserList(c *gin.Context) {

	var req dto.GetUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, utils.NewErrorResponseStruct("invalid query parameters.", err.Error()))
		return
	}

	if errs := req.Validate(); len(errs) > 0 {
		c.JSON(400, utils.NewErrorResponseStruct("validation failed.", utils.ErrorMessages(errs)))
		return
	}

	users, pagination, err := h.userService.GetUserList(c, req.Page, req.PageSize)
	if err != nil {
		c.JSON(500, utils.NewErrorResponseStruct("failed to retrieve users.", err.Error()))
		return
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(200, utils.NewSuccessResponseStruct("users retrieved successfully.", dto.UserListResponse{
		Users:    userResponses,
		Total:    pagination.Total,
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	}))
}

// UpdateUser updates a user's name, email in the database by ID.
func (h *userHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(400, utils.NewErrorResponseStruct("validation failed.", "id is required."))
		return
	}

	if !utils.IsUUIDValid(id) {
		c.JSON(400, utils.NewErrorResponseStruct("validation failed.", "id is not a valid uuid."))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponseStruct("invalid request body.", err.Error()))
		return
	}

	if errs := req.Validate(); len(errs) > 0 {
		c.JSON(400, utils.NewErrorResponseStruct("validation failed.", utils.ErrorMessages(errs)))
		return
	}

	user, err := h.userService.Update(c, id, req.Name, req.Email)
	if err != nil {
		c.JSON(500, utils.NewErrorResponseStruct("failed to update user.", err.Error()))
		return
	}
	c.JSON(200, utils.NewSuccessResponseStruct("user updated successfully.", dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}))
}

// DeleteUser deletes a user from the database by ID.
func (h *userHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(400, utils.NewErrorResponseStruct("validation failed.", "id is required."))
		return
	}

	if !utils.IsUUIDValid(id) {
		c.JSON(400, utils.NewErrorResponseStruct("validation failed.", "id is not a valid uuid."))
		return
	}

	err := h.userService.Delete(c, id)
	if err != nil {
		c.JSON(500, utils.NewErrorResponseStruct("failed to delete user.", err.Error()))
		return
	}
	c.JSON(200, utils.NewSuccessResponseStruct("user deleted successfully.", nil))
}
