package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akrindev/gonely/internal/application/feeds"
	domain "github.com/akrindev/gonely/internal/domain/feeds"
	"github.com/akrindev/gonely/internal/infrastructure/logger"
	"github.com/akrindev/gonely/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handlerStubPostRepository struct {
	posts []*domain.Post
	err   error
}

func (s *handlerStubPostRepository) Create(_ context.Context, post *domain.Post) error {
	if s.err != nil {
		return s.err
	}
	s.posts = append(s.posts, post)
	return nil
}

func (s *handlerStubPostRepository) FindByID(_ context.Context, id string) (*domain.Post, error) {
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

func (s *handlerStubPostRepository) FindAll(_ context.Context, limit, offset int) ([]*domain.Post, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.posts, nil
}

func (s *handlerStubPostRepository) Update(_ context.Context, post *domain.Post) error {
	for i, p := range s.posts {
		if p.ID == post.ID {
			s.posts[i] = post
			break
		}
	}
	return s.err
}

func (s *handlerStubPostRepository) Delete(_ context.Context, id string) error {
	return s.err
}

type handlerStubCommentRepository struct {
	comments []*domain.Comment
	err      error
}

func (s *handlerStubCommentRepository) Create(_ context.Context, comment *domain.Comment) error {
	if s.err != nil {
		return s.err
	}
	s.comments = append(s.comments, comment)
	return nil
}

func (s *handlerStubCommentRepository) FindByID(_ context.Context, id string) (*domain.Comment, error) {
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

func (s *handlerStubCommentRepository) FindByPostID(_ context.Context, postID string, limit, offset int) ([]*domain.Comment, error) {
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

func (s *handlerStubCommentRepository) Update(_ context.Context, comment *domain.Comment) error {
	return s.err
}

func (s *handlerStubCommentRepository) Delete(_ context.Context, id string) error {
	return s.err
}

func newHandlerTestService(initialPosts []*domain.Post) *feeds.FeedService {
	postRepo := &handlerStubPostRepository{posts: initialPosts}
	commentRepo := &handlerStubCommentRepository{}
	return feeds.NewFeedService(postRepo, commentRepo, &logger.Logger{Logger: zap.NewNop()})
}

func TestFeedHandler_GetPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := newHandlerTestService([]*domain.Post{domain.NewPost("hello")})
	handler := NewFeedHandler(service)

	r := gin.New()
	r.GET("/posts", handler.GetPosts)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected error decoding response: %v", err)
	}

	data, ok := response["data"].([]interface{})
	if !ok || len(data) != 1 {
		t.Fatalf("expected data array with one post")
	}
}

func TestFeedHandler_GetPosts_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	postRepo := &handlerStubPostRepository{err: errors.New("db error")}
	service := feeds.NewFeedService(postRepo, &handlerStubCommentRepository{}, &logger.Logger{Logger: zap.NewNop()})
	handler := NewFeedHandler(service)

	r := gin.New()
	r.Use(middleware.ErrorHandler(&logger.Logger{Logger: zap.NewNop()}))
	r.GET("/posts", handler.GetPosts)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
