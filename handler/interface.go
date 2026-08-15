// package handler provides interfaces and implementations for handlers.
package handler

import "github.com/gin-gonic/gin"

// AuthHandler interface for authentication handlers.
type AuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
}

// UserHandler interface for user handlers.
type UserHandler interface {
	CreateUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetUserList(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}
