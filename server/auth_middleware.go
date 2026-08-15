package server

import (
	"strings"
	"user-mgmt/utils"

	"github.com/gin-gonic/gin"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Authorization: Bearer JWT type
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Split Authorization header to get JWT token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid Authorization header"})
			c.Abort()
			return
		}

		// Verify token
		claims, err := utils.VerifyToken(&s.config.JWT, tokenParts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set claims to context
		c.Set("user_id", claims.ID)
		c.Set("user_name", claims.Name)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
