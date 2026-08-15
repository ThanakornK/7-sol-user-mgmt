package server

import (
	"context"
	"user-mgmt/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// traceMiddleware adds a trace ID to the request context.
func traceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if _, err := uuid.Parse(traceID); err != nil {
			traceID = uuid.NewString()
		}
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)
		requestContext := context.WithValue(c.Request.Context(), utils.TraceIDContextKey, traceID)
		c.Request = c.Request.WithContext(requestContext)
		c.Next()
	}
}
