package feeds

import (
	"context"
	"fmt"

	"github.com/akrindev/gonely/internal/domain/feeds"
	"github.com/akrindev/gonely/internal/infrastructure/logger"
	"go.uber.org/zap"
)

// FeedService handles feed-related business logic
type FeedService struct {
	postRepo    feeds.PostRepository
	commentRepo feeds.CommentRepository
	logger      *logger.Logger
}

// NewFeedService creates a new feed service
func NewFeedService(
	postRepo feeds.PostRepository,
	commentRepo feeds.CommentRepository,
	log *logger.Logger,
) *FeedService {
	return &FeedService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		logger:      log.Named("feed-service"),
	}
}

// GetAllPosts retrieves all posts with pagination
func (s *FeedService) GetAllPosts(ctx context.Context, limit, offset int) ([]*feeds.Post, error) {
	s.logger.Info("Fetching all posts", zap.Int("limit", limit), zap.Int("offset", offset))

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, err := s.postRepo.FindAll(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to fetch posts", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}

	return posts, nil
}

// GetPostByID retrieves a post by its ID with comments
func (s *FeedService) GetPostByID(ctx context.Context, id string) (*feeds.Post, error) {
	s.logger.Info("Fetching post by ID", zap.String("id", id))

	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch post", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch post: %w", err)
	}

	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	return post, nil
}

// CreatePost creates a new post
func (s *FeedService) CreatePost(ctx context.Context, content string) (*feeds.Post, error) {
	s.logger.Info("Creating new post", zap.Int("content_length", len(content)))

	if content == "" {
		return nil, fmt.Errorf("post content cannot be empty")
	}

	post := feeds.NewPost(content)
	if err := s.postRepo.Create(ctx, post); err != nil {
		s.logger.Error("Failed to create post", zap.Error(err))
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// AddComment adds a comment to a post
func (s *FeedService) AddComment(ctx context.Context, postID, content string) (*feeds.Comment, error) {
	s.logger.Info("Adding comment to post", zap.String("post_id", postID))

	if content == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}

	// Verify post exists
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to find post", zap.String("post_id", postID), zap.Error(err))
		return nil, fmt.Errorf("failed to find post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	comment := feeds.NewComment(postID, content)
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		s.logger.Error("Failed to create comment", zap.Error(err))
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// GetCommentsByPostID retrieves comments for a post
func (s *FeedService) GetCommentsByPostID(ctx context.Context, postID string, limit, offset int) ([]*feeds.Comment, error) {
	s.logger.Info("Fetching comments for post", zap.String("post_id", postID))

	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	comments, err := s.commentRepo.FindByPostID(ctx, postID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to fetch comments", zap.String("post_id", postID), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	return comments, nil
}
