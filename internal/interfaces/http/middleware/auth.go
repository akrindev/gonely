package middleware

import (
	"net/http"
	"strings"

	"github.com/akrindev/gonely/internal/application/auth"
	"github.com/akrindev/gonely/internal/interfaces/http/dto"
	pkgerrors "github.com/akrindev/gonely/pkg/errors"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserContextKey      = "user"
	SessionContextKey   = "session"
)

// AuthMiddleware creates a middleware for authentication
func AuthMiddleware(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Missing authorization header",
			))
			return
		}

		// Extract token - support both "Bearer <token>" and just "<token>" formats
		var token string
		if strings.HasPrefix(authHeader, BearerPrefix) {
			token = strings.TrimPrefix(authHeader, BearerPrefix)
		} else {
			// Assume the entire header value is the token (for Swagger UI compatibility)
			token = authHeader
		}

		// Validate token
		session, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Invalid or expired token",
			))
			return
		}

		// Set user and session in context
		c.Set(UserContextKey, session.User)
		c.Set(SessionContextKey, session)

		c.Next()
	}
}

// OptionalAuthMiddleware creates a middleware for optional authentication
func OptionalAuthMiddleware(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)
		session, err := authService.ValidateToken(c.Request.Context(), token)
		if err == nil && session != nil {
			c.Set(UserContextKey, session.User)
			c.Set(SessionContextKey, session)
		}

		c.Next()
	}
}

// GetUserFromContext retrieves the user from the context
func GetUserFromContext(c *gin.Context) (interface{}, error) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil, pkgerrors.Unauthorized("User not authenticated")
	}
	return user, nil
}
