package dto

// CreatePostRequest represents a request to create a post
type CreatePostRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000"`
}

// CreateCommentRequest represents a request to create a comment
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=2000"`
}
