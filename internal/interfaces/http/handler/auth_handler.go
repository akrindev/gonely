package handler

import (
	"net/http"

	"github.com/akrindev/gonely/internal/application/auth"
	authDomain "github.com/akrindev/gonely/internal/domain/auth"
	"github.com/akrindev/gonely/internal/interfaces/http/dto"
	pkgerrors "github.com/akrindev/gonely/pkg/errors"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /auth/register
// @Summary		Register a new user
// @Description	Register a new user with name, email and password
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		dto.RegisterRequest	true	"User registration data"
// @Success		201		{object}	dto.Response
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(pkgerrors.BadRequest("Invalid request body"))
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.Error(pkgerrors.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewResponse(user))
}

// Login handles POST /auth/login
// @Summary		Login user
// @Description	Authenticate user with email and password
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		dto.LoginRequest	true	"User login credentials"
// @Success		200		{object}	dto.Response
// @Failure		400		{object}	dto.ErrorResponse
// @Failure		401		{object}	dto.ErrorResponse
// @Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(pkgerrors.BadRequest("Invalid request body"))
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	token, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password, userAgent, ipAddress)
	if err != nil {
		c.Error(pkgerrors.Unauthorized(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewResponse(dto.AuthResponse{
		Token: token,
		User:  user,
	}))
}

// CreateAnonymous handles POST /auth/anonymous
// @Summary		Create anonymous user
// @Description	Create a new anonymous user session
// @Tags			auth
// @Accept			json
// @Produce		json
// @Success		201	{object}	dto.Response
// @Failure		500	{object}	dto.ErrorResponse
// @Router			/auth/anonymous [post]
func (h *AuthHandler) CreateAnonymous(c *gin.Context) {
	user, token, err := h.authService.CreateAnonymousUser(c.Request.Context())
	if err != nil {
		c.Error(pkgerrors.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewResponse(dto.AuthResponse{
		Token: token,
		User:  user,
	}))
}

// Logout handles POST /auth/logout
// @Summary		Logout user
// @Description	Invalidate the current user session
// @Tags			auth
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	map[string]string
// @Failure		400	{object}	dto.ErrorResponse
// @Failure		500	{object}	dto.ErrorResponse
// @Router			/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get the session from context (already validated by middleware)
	session, exists := c.Get("session")
	if !exists {
		c.Error(pkgerrors.InternalError("Session not found in context"))
		return
	}

	authSession, ok := session.(*authDomain.Session)
	if !ok {
		c.Error(pkgerrors.InternalError("Invalid session type"))
		return
	}

	// Revoke the session
	authSession.Revoke("user logout")
	if err := h.authService.RevokeSession(c.Request.Context(), authSession); err != nil {
		c.Error(pkgerrors.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Me handles GET /auth/me
// @Summary		Get current user
// @Description	Get information about the currently authenticated user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	dto.Response
// @Failure		401	{object}	dto.ErrorResponse
// @Router			/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.Error(pkgerrors.Unauthorized("User not authenticated"))
		return
	}

	c.JSON(http.StatusOK, dto.NewResponse(user))
}
