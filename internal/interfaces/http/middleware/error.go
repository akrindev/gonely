package middleware

import (
	"net/http"

	"github.com/akrindev/gonely/internal/infrastructure/logger"
	"github.com/akrindev/gonely/internal/interfaces/http/dto"
	pkgerrors "github.com/akrindev/gonely/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler creates a middleware for centralized error handling
func ErrorHandler(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		// Handle AppError
		if appErr, ok := err.(*pkgerrors.AppError); ok {
			c.JSON(appErr.StatusCode, dto.NewErrorResponse(appErr.Code, appErr.Message))
			return
		}

		// Log unexpected errors
		log.Error("Unexpected error occurred",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		// Return generic error response
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
			"INTERNAL_ERROR",
			"An unexpected error occurred",
		))
	}
}
