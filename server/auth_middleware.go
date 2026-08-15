package server

import (
	"context"
	"net/http"
	"strings"
	"time"
	"user-mgmt/utils"

	"github.com/gin-gonic/gin"
)

// authMiddleware validates the access token in the request header and sets the user ID in the request context.
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get token from Authorization header.
		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, utils.NewErrorResponseStruct("authentication failed", "invalid access token"))
			c.Abort()
			return
		}

		// Verify access token.
		claims, err := utils.VerifyAccessToken(&s.config.JWT, parts[1], time.Now().UTC())
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.NewErrorResponseStruct("authentication failed", "invalid access token"))
			c.Abort()
			return
		}

		// Set user ID in request context.
		c.Set("user_id", claims.Subject)
		requestContext := context.WithValue(c.Request.Context(), utils.UserIDContextKey, claims.Subject)
		c.Request = c.Request.WithContext(requestContext)
		c.Next()
	}
}
