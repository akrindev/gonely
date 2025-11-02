package handler

import (
	"net/http"
	"strconv"

	"github.com/akrindev/gonely/internal/application/feeds"
	"github.com/akrindev/gonely/internal/interfaces/http/dto"
	pkgerrors "github.com/akrindev/gonely/pkg/errors"
	"github.com/gin-gonic/gin"
)

// FeedHandler handles HTTP requests for feeds
type FeedHandler struct {
	feedService *feeds.FeedService
}

// NewFeedHandler creates a new feed handler
func NewFeedHandler(feedService *feeds.FeedService) *FeedHandler {
	return &FeedHandler{feedService: feedService}
}

// GetPosts handles GET /v1/posts
func (h *FeedHandler) GetPosts(c *gin.Context) {
	// Parse pagination parameters
	page := parseIntOrDefault(c.Query("page"), 0)
	limit := parseIntOrDefault(c.Query("limit"), 20)
	offset := page * limit

	posts, err := h.feedService.GetAllPosts(c.Request.Context(), limit, offset)
	if err != nil {
		c.Error(pkgerrors.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewListResponse(posts, page, limit, offset))
}

// GetPostByID handles GET /v1/posts/:id
func (h *FeedHandler) GetPostByID(c *gin.Context) {
	id := c.Param("id")

	post, err := h.feedService.GetPostByID(c.Request.Context(), id)
	if err != nil {
		c.Error(pkgerrors.NotFound("Post not found"))
		return
	}

	c.JSON(http.StatusOK, dto.NewResponse(post))
}

// CreatePost handles POST /v1/posts
func (h *FeedHandler) CreatePost(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(pkgerrors.BadRequest("Invalid request body"))
		return
	}

	post, err := h.feedService.CreatePost(c.Request.Context(), req.Content)
	if err != nil {
		c.Error(pkgerrors.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewResponse(post))
}

// AddComment handles POST /v1/posts/:id/comments
func (h *FeedHandler) AddComment(c *gin.Context) {
	postID := c.Param("id")

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(pkgerrors.BadRequest("Invalid request body"))
		return
	}

	comment, err := h.feedService.AddComment(c.Request.Context(), postID, req.Content)
	if err != nil {
		c.Error(pkgerrors.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewResponse(comment))
}

// GetComments handles GET /v1/posts/:id/comments
func (h *FeedHandler) GetComments(c *gin.Context) {
	postID := c.Param("id")

	page := parseIntOrDefault(c.Query("page"), 0)
	limit := parseIntOrDefault(c.Query("limit"), 50)
	offset := page * limit

	comments, err := h.feedService.GetCommentsByPostID(c.Request.Context(), postID, limit, offset)
	if err != nil {
		c.Error(pkgerrors.InternalError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewListResponse(comments, page, limit, offset))
}

// Helper function to parse int or return default
func parseIntOrDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}
