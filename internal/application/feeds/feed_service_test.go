package feeds

import (
	"context"
	"errors"
	"testing"

	domain "github.com/akrindev/gonely/internal/domain/feeds"
	"github.com/akrindev/gonely/internal/infrastructure/logger"
	"go.uber.org/zap"
)

type stubPostRepository struct {
	posts []*domain.Post
	err   error
}

func (s *stubPostRepository) Create(ctx context.Context, post *domain.Post) error {
	s.posts = append(s.posts, post)
	return s.err
}

func (s *stubPostRepository) FindByID(ctx context.Context, id string) (*domain.Post, error) {
	if s.err != nil {
		return nil, s.err
	}
	for _, post := range s.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, nil
}

func (s *stubPostRepository) FindAll(ctx context.Context, limit, offset int) ([]*domain.Post, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.posts, nil
}

func (s *stubPostRepository) Update(ctx context.Context, post *domain.Post) error {
	return s.err
}

func (s *stubPostRepository) Delete(ctx context.Context, id string) error {
	return s.err
}

type stubCommentRepository struct {
	comments []*domain.Comment
	err      error
}

func (s *stubCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	s.comments = append(s.comments, comment)
	return s.err
}

func (s *stubCommentRepository) FindByID(ctx context.Context, id string) (*domain.Comment, error) {
	if s.err != nil {
		return nil, s.err
	}
	for _, comment := range s.comments {
		if comment.ID == id {
			return comment, nil
		}
	}
	return nil, nil
}

func (s *stubCommentRepository) FindByPostID(ctx context.Context, postID string, limit, offset int) ([]*domain.Comment, error) {
	if s.err != nil {
		return nil, s.err
	}
	var result []*domain.Comment
	for _, comment := range s.comments {
		if comment.PostID == postID {
			result = append(result, comment)
		}
	}
	return result, nil
}

func (s *stubCommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	return s.err
}

func (s *stubCommentRepository) Delete(ctx context.Context, id string) error {
	return s.err
}

func newTestLogger() *logger.Logger {
	return &logger.Logger{Logger: zap.NewNop()}
}

func TestFeedService_GetAllPosts(t *testing.T) {
	postRepo := &stubPostRepository{}
	commentRepo := &stubCommentRepository{}

	postRepo.Create(context.Background(), domain.NewPost("hello world"))

	service := NewFeedService(postRepo, commentRepo, newTestLogger())

	posts, err := service.GetAllPosts(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if posts[0].Content != "hello world" {
		t.Fatalf("unexpected post content: %s", posts[0].Content)
	}
}

func TestFeedService_GetPostByID_NotFound(t *testing.T) {
	postRepo := &stubPostRepository{}
	commentRepo := &stubCommentRepository{}
	service := NewFeedService(postRepo, commentRepo, newTestLogger())

	if _, err := service.GetPostByID(context.Background(), "missing"); err == nil {
		t.Fatalf("expected error when post is missing")
	}
}

func TestFeedService_GetPostByID_Error(t *testing.T) {
	postRepo := &stubPostRepository{err: errors.New("db error")}
	commentRepo := &stubCommentRepository{}
	service := NewFeedService(postRepo, commentRepo, newTestLogger())

	if _, err := service.GetPostByID(context.Background(), "id"); err == nil {
		t.Fatalf("expected error when repository fails")
	}
}

func TestFeedService_CreatePost_Validation(t *testing.T) {
	postRepo := &stubPostRepository{}
	commentRepo := &stubCommentRepository{}
	service := NewFeedService(postRepo, commentRepo, newTestLogger())

	if _, err := service.CreatePost(context.Background(), ""); err == nil {
		t.Fatalf("expected validation error for empty content")
	}
}

func TestFeedService_AddComment_PostMissing(t *testing.T) {
	postRepo := &stubPostRepository{}
	commentRepo := &stubCommentRepository{}
	service := NewFeedService(postRepo, commentRepo, newTestLogger())

	if _, err := service.AddComment(context.Background(), "missing", "hi"); err == nil {
		t.Fatalf("expected error when post is missing")
	}
}
