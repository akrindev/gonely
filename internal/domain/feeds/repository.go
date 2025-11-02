package feeds

import "context"

// PostRepository defines the interface for post data operations
type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	FindByID(ctx context.Context, id string) (*Post, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id string) error
}

// CommentRepository defines the interface for comment data operations
type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	FindByID(ctx context.Context, id string) (*Comment, error)
	FindByPostID(ctx context.Context, postID string, limit, offset int) ([]*Comment, error)
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id string) error
}
