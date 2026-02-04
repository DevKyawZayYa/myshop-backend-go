package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterBasicMiddleware(r *gin.Engine) {
	// gin.Default() already includes Logger and Recovery middleware
	// Add custom RequestID middleware
	r.Use(RequestIDMiddleware())
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in header
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context and response header
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
