package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appauth "github.com/akrindev/gonely/internal/application/auth"
	domain "github.com/akrindev/gonely/internal/domain/auth"
	"github.com/akrindev/gonely/internal/infrastructure/config"
	"github.com/akrindev/gonely/internal/infrastructure/logger"
	"github.com/akrindev/gonely/internal/interfaces/http/dto"
	"github.com/akrindev/gonely/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handlerUserRepository struct {
	byID    map[string]*domain.User
	byEmail map[string]*domain.User
	fail    error
}

func newHandlerUserRepository() *handlerUserRepository {
	return &handlerUserRepository{
		byID:    make(map[string]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

func (r *handlerUserRepository) Create(_ context.Context, user *domain.User) error {
	if r.fail != nil {
		return r.fail
	}
	r.byID[user.ID] = user
	r.byEmail[user.Email] = user
	return nil
}

func (r *handlerUserRepository) FindByID(_ context.Context, id string) (*domain.User, error) {
	if r.fail != nil {
		return nil, r.fail
	}
	return r.byID[id], nil
}

func (r *handlerUserRepository) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	if r.fail != nil {
		return nil, r.fail
	}
	return r.byEmail[email], nil
}

func (r *handlerUserRepository) FindByHandle(_ context.Context, handle string) (*domain.User, error) {
	return nil, r.fail
}

func (r *handlerUserRepository) Update(_ context.Context, user *domain.User) error {
	if r.fail != nil {
		return r.fail
	}
	r.byID[user.ID] = user
	r.byEmail[user.Email] = user
	return nil
}

func (r *handlerUserRepository) Delete(_ context.Context, id string) error {
	return r.fail
}

func (r *handlerUserRepository) List(_ context.Context, limit, offset int) ([]*domain.User, error) {
	if r.fail != nil {
		return nil, r.fail
	}
	var result []*domain.User
	for _, user := range r.byID {
		result = append(result, user)
	}
	return result, nil
}

type handlerAccountRepository struct {
	accounts []*domain.Account
	fail     error
}

func (r *handlerAccountRepository) Create(_ context.Context, account *domain.Account) error {
	if r.fail != nil {
		return r.fail
	}
	r.accounts = append(r.accounts, account)
	return nil
}

func (r *handlerAccountRepository) FindByID(_ context.Context, id string) (*domain.Account, error) {
	return nil, r.fail
}

func (r *handlerAccountRepository) FindByUserID(_ context.Context, userID string) ([]*domain.Account, error) {
	if r.fail != nil {
		return nil, r.fail
	}
	var result []*domain.Account
	for _, account := range r.accounts {
		if account.UserID == userID {
			result = append(result, account)
		}
	}
	return result, nil
}

func (r *handlerAccountRepository) FindByProviderAndAccountID(_ context.Context, providerID, accountID string) (*domain.Account, error) {
	return nil, r.fail
}

func (r *handlerAccountRepository) Update(_ context.Context, account *domain.Account) error {
	return r.fail
}

func (r *handlerAccountRepository) Delete(_ context.Context, id string) error {
	return r.fail
}

type handlerSessionRepository struct {
	sessions map[string]*domain.Session
	fail     error
}

func newHandlerSessionRepository() *handlerSessionRepository {
	return &handlerSessionRepository{sessions: make(map[string]*domain.Session)}
}

func (r *handlerSessionRepository) Create(_ context.Context, session *domain.Session) error {
	if r.fail != nil {
		return r.fail
	}
	r.sessions[session.Token] = session
	return nil
}

func (r *handlerSessionRepository) FindByID(_ context.Context, id string) (*domain.Session, error) {
	if r.fail != nil {
		return nil, r.fail
	}
	for _, session := range r.sessions {
		if session.ID == id {
			return session, nil
		}
	}
	return nil, nil
}

func (r *handlerSessionRepository) FindByToken(_ context.Context, token string) (*domain.Session, error) {
	if r.fail != nil {
		return nil, r.fail
	}
	return r.sessions[token], nil
}

func (r *handlerSessionRepository) FindByUserID(_ context.Context, userID string) ([]*domain.Session, error) {
	if r.fail != nil {
		return nil, r.fail
	}
	var result []*domain.Session
	for _, session := range r.sessions {
		if session.UserID == userID {
			result = append(result, session)
		}
	}
	return result, nil
}

func (r *handlerSessionRepository) Update(_ context.Context, session *domain.Session) error {
	if r.fail != nil {
		return r.fail
	}
	r.sessions[session.Token] = session
	return nil
}

func (r *handlerSessionRepository) Delete(_ context.Context, id string) error {
	return r.fail
}

func (r *handlerSessionRepository) DeleteExpired(_ context.Context) error {
	return r.fail
}

func (r *handlerSessionRepository) RevokeAllByUserID(_ context.Context, userID string) error {
	return r.fail
}

func newAuthHandlerService() (*appauth.AuthService, *handlerUserRepository, *handlerAccountRepository, *handlerSessionRepository) {
	userRepo := newHandlerUserRepository()
	accountRepo := &handlerAccountRepository{}
	sessionRepo := newHandlerSessionRepository()
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:     "test-secret-must-be-at-least-32-chars-long",
			JWTExpiry:     time.Hour,
			SessionExpiry: time.Hour,
		},
	}
	log := &logger.Logger{Logger: zap.NewNop()}
	service := appauth.NewAuthService(userRepo, sessionRepo, accountRepo, cfg, log)
	return service, userRepo, accountRepo, sessionRepo
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service, userRepo, _, _ := newAuthHandlerService()
	handler := NewAuthHandler(service)

	r := gin.New()
	r.Use(middleware.ErrorHandler(&logger.Logger{Logger: zap.NewNop()}))
	r.POST("/auth/register", handler.Register)

	payload := dto.RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "SecurePass123"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	if _, exists := userRepo.byEmail[payload.Email]; !exists {
		t.Fatalf("expected user to be stored")
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service, _, accountRepo, sessionRepo := newAuthHandlerService()
	handler := NewAuthHandler(service)

	// Seed user via service.Register to ensure password hashing
	_, err := service.Register(context.Background(), "Alice", "alice@example.com", "SecurePass123")
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	r := gin.New()
	r.Use(middleware.ErrorHandler(&logger.Logger{Logger: zap.NewNop()}))
	r.POST("/auth/login", handler.Login)

	payload := dto.LoginRequest{Email: "alice@example.com", Password: "SecurePass123"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if len(sessionRepo.sessions) != 1 {
		t.Fatalf("expected session to be created")
	}

	if len(accountRepo.accounts) == 0 || accountRepo.accounts[0].Password == nil {
		t.Fatalf("expected hashed password in account")
	}
}

func TestAuthHandler_Login_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service, _, _, _ := newAuthHandlerService()
	handler := NewAuthHandler(service)

	r := gin.New()
	r.Use(middleware.ErrorHandler(&logger.Logger{Logger: zap.NewNop()}))
	r.POST("/auth/login", handler.Login)

	payload := dto.LoginRequest{Email: "unknown@example.com", Password: "bad"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service, _, _, sessionRepo := newAuthHandlerService()
	handler := NewAuthHandler(service)

	_, err := service.Register(context.Background(), "Alice", "alice@example.com", "SecurePass123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	token, _, err := service.Login(context.Background(), "alice@example.com", "SecurePass123", "UA", "127.0.0.1")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	r := gin.New()
	r.Use(middleware.ErrorHandler(&logger.Logger{Logger: zap.NewNop()}))
	r.POST("/auth/logout", handler.Logout)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	session := sessionRepo.sessions[token]
	if session == nil || !session.Revoked {
		t.Fatalf("expected session to be revoked")
	}
}

func TestAuthHandler_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(nil)

	r := gin.New()
	r.Use(middleware.ErrorHandler(&logger.Logger{Logger: zap.NewNop()}))
	r.GET("/auth/me", func(c *gin.Context) {
		user := domain.NewUser("Alice", "alice@example.com", false)
		c.Set("user", user)
		handler.Me(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestAuthHandler_Register_Duplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service, userRepo, _, _ := newAuthHandlerService()
	handler := NewAuthHandler(service)

	user := domain.NewUser("Alice", "alice@example.com", false)
	userRepo.Create(context.Background(), user)

	r := gin.New()
	r.Use(middleware.ErrorHandler(&logger.Logger{Logger: zap.NewNop()}))
	r.POST("/auth/register", handler.Register)

	payload := dto.RegisterRequest{Name: "Alice", Email: "alice@example.com", Password: "SecurePass123"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestAuthHandler_CreateAnonymous(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service, _, _, sessionRepo := newAuthHandlerService()
	handler := NewAuthHandler(service)

	r := gin.New()
	r.Use(middleware.ErrorHandler(&logger.Logger{Logger: zap.NewNop()}))
	r.POST("/auth/anonymous", handler.CreateAnonymous)

	req := httptest.NewRequest(http.MethodPost, "/auth/anonymous", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	if len(sessionRepo.sessions) == 0 {
		t.Fatalf("expected anonymous session to be created")
	}
}

