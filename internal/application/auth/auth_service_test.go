package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/akrindev/gonely/internal/domain/auth"
	systemconfig "github.com/akrindev/gonely/internal/infrastructure/config"
	"github.com/akrindev/gonely/internal/infrastructure/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type stubUserRepository struct {
	stored []*domain.User
	byEmail map[string]*domain.User
	err error
}

func (s *stubUserRepository) Create(ctx context.Context, user *domain.User) error {
	if s.err != nil {
		return s.err
	}
	s.stored = append(s.stored, user)
	if s.byEmail == nil {
		s.byEmail = make(map[string]*domain.User)
	}
	s.byEmail[user.Email] = user
	return nil
}

func (s *stubUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if s.err != nil {
		return nil, s.err
	}
	for _, user := range s.stored {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func (s *stubUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if s.err != nil {
		return nil, s.err
	}
	if user, ok := s.byEmail[email]; ok {
		return user, nil
	}
	return nil, nil
}

func (s *stubUserRepository) FindByHandle(ctx context.Context, handle string) (*domain.User, error) {
	return nil, s.err
}

func (s *stubUserRepository) Update(ctx context.Context, user *domain.User) error {
	return s.err
}

func (s *stubUserRepository) Delete(ctx context.Context, id string) error {
	return s.err
}

func (s *stubUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return s.stored, s.err
}

type stubAccountRepository struct {
	accounts []*domain.Account
	err      error
}

func (s *stubAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	if s.err != nil {
		return s.err
	}
	s.accounts = append(s.accounts, account)
	return nil
}

func (s *stubAccountRepository) FindByID(ctx context.Context, id string) (*domain.Account, error) {
	return nil, s.err
}

func (s *stubAccountRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Account, error) {
	if s.err != nil {
		return nil, s.err
	}
	var result []*domain.Account
	for _, account := range s.accounts {
		if account.UserID == userID {
			result = append(result, account)
		}
	}
	return result, nil
}

func (s *stubAccountRepository) FindByProviderAndAccountID(ctx context.Context, providerID, accountID string) (*domain.Account, error) {
	return nil, s.err
}

func (s *stubAccountRepository) Update(ctx context.Context, account *domain.Account) error {
	return s.err
}

func (s *stubAccountRepository) Delete(ctx context.Context, id string) error {
	return s.err
}

type stubSessionRepository struct {
	sessions map[string]*domain.Session
	err      error
}

func (s *stubSessionRepository) ensureMap() {
	if s.sessions == nil {
		s.sessions = make(map[string]*domain.Session)
	}
}

func (s *stubSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	if s.err != nil {
		return s.err
	}
	s.ensureMap()
	s.sessions[session.Token] = session
	return nil
}

func (s *stubSessionRepository) FindByID(ctx context.Context, id string) (*domain.Session, error) {
	if s.err != nil {
		return nil, s.err
	}
	for _, session := range s.sessions {
		if session.ID == id {
			return session, nil
		}
	}
	return nil, nil
}

func (s *stubSessionRepository) FindByToken(ctx context.Context, token string) (*domain.Session, error) {
	if s.err != nil {
		return nil, s.err
	}
	s.ensureMap()
	return s.sessions[token], nil
}

func (s *stubSessionRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Session, error) {
	if s.err != nil {
		return nil, s.err
	}
	var result []*domain.Session
	for _, session := range s.sessions {
		if session.UserID == userID {
			result = append(result, session)
		}
	}
	return result, nil
}

func (s *stubSessionRepository) Update(ctx context.Context, session *domain.Session) error {
	if s.err != nil {
		return s.err
	}
	s.ensureMap()
	s.sessions[session.Token] = session
	return nil
}

func (s *stubSessionRepository) Delete(ctx context.Context, id string) error {
	return s.err
}

func (s *stubSessionRepository) DeleteExpired(ctx context.Context) error {
	return s.err
}

func (s *stubSessionRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	return s.err
}

func newAuthTestLogger() *logger.Logger {
	return &logger.Logger{Logger: zap.NewNop()}
}

func testConfig() *systemconfig.Config {
	return &systemconfig.Config{
		Auth: systemconfig.AuthConfig{
			JWTSecret:     "test-secret-should-be-long-enough",
			JWTExpiry:     time.Hour,
			SessionExpiry: time.Hour,
		},
	}
}

func TestAuthService_RegisterSuccess(t *testing.T) {
	userRepo := &stubUserRepository{}
	accountRepo := &stubAccountRepository{}
	sessionRepo := &stubSessionRepository{}
	service := NewAuthService(userRepo, sessionRepo, accountRepo, testConfig(), newAuthTestLogger())

	user, err := service.Register(context.Background(), "Alice", "alice@example.com", "SecurePass123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil || user.Email != "alice@example.com" {
		t.Fatalf("user not created correctly")
	}
	if len(accountRepo.accounts) != 1 {
		t.Fatalf("expected account to be created")
	}
	if accountRepo.accounts[0].Password == nil {
		t.Fatalf("expected password hash to be stored")
	}
}

func TestAuthService_RegisterDuplicateEmail(t *testing.T) {
	userRepo := &stubUserRepository{byEmail: map[string]*domain.User{
		"alice@example.com": domain.NewUser("Alice", "alice@example.com", false),
	}}
	accountRepo := &stubAccountRepository{}
	sessionRepo := &stubSessionRepository{}
	service := NewAuthService(userRepo, sessionRepo, accountRepo, testConfig(), newAuthTestLogger())

	if _, err := service.Register(context.Background(), "Alice", "alice@example.com", "SecurePass123"); err == nil {
		t.Fatalf("expected error for duplicate email")
	}
}

func TestAuthService_LoginSuccess(t *testing.T) {
	user := domain.NewUser("Alice", "alice@example.com", false)
	hash, _ := bcrypt.GenerateFromPassword([]byte("SecurePass123"), bcrypt.DefaultCost)
	account := domain.NewAccount(user.Email, "credentials", user.ID)
	account.SetPassword(string(hash))

	userRepo := &stubUserRepository{byEmail: map[string]*domain.User{user.Email: user}}
	accountRepo := &stubAccountRepository{accounts: []*domain.Account{account}}
	sessionRepo := &stubSessionRepository{}
	service := NewAuthService(userRepo, sessionRepo, accountRepo, testConfig(), newAuthTestLogger())

	token, loggedInUser, err := service.Login(context.Background(), user.Email, "SecurePass123", "UA", "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatalf("expected token to be returned")
	}
	if loggedInUser.Email != user.Email {
		t.Fatalf("unexpected user returned")
	}
	if sessionRepo.sessions[token] == nil {
		t.Fatalf("expected session to be stored")
	}
}

func TestAuthService_LoginInvalidPassword(t *testing.T) {
	user := domain.NewUser("Alice", "alice@example.com", false)
	hash, _ := bcrypt.GenerateFromPassword([]byte("SecurePass123"), bcrypt.DefaultCost)
	account := domain.NewAccount(user.Email, "credentials", user.ID)
	account.SetPassword(string(hash))

	userRepo := &stubUserRepository{byEmail: map[string]*domain.User{user.Email: user}}
	accountRepo := &stubAccountRepository{accounts: []*domain.Account{account}}
	sessionRepo := &stubSessionRepository{}
	service := NewAuthService(userRepo, sessionRepo, accountRepo, testConfig(), newAuthTestLogger())

	if _, _, err := service.Login(context.Background(), user.Email, "WrongPass", "UA", "127.0.0.1"); err == nil {
		t.Fatalf("expected error for invalid password")
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	user := domain.NewUser("Alice", "alice@example.com", false)
	hash, _ := bcrypt.GenerateFromPassword([]byte("SecurePass123"), bcrypt.DefaultCost)
	account := domain.NewAccount(user.Email, "credentials", user.ID)
	account.SetPassword(string(hash))

	userRepo := &stubUserRepository{byEmail: map[string]*domain.User{user.Email: user}}
	accountRepo := &stubAccountRepository{accounts: []*domain.Account{account}}
	sessionRepo := &stubSessionRepository{}
	service := NewAuthService(userRepo, sessionRepo, accountRepo, testConfig(), newAuthTestLogger())

	token, _, err := service.Login(context.Background(), user.Email, "SecurePass123", "UA", "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}

	session, err := service.ValidateToken(context.Background(), token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if session.UserID != user.ID {
		t.Fatalf("unexpected user ID on session")
	}
}

func TestAuthService_CreateAnonymousUser(t *testing.T) {
	userRepo := &stubUserRepository{}
	accountRepo := &stubAccountRepository{}
	sessionRepo := &stubSessionRepository{}
	service := NewAuthService(userRepo, sessionRepo, accountRepo, testConfig(), newAuthTestLogger())

	user, token, err := service.CreateAnonymousUser(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !user.IsAnonymous {
		t.Fatalf("expected anonymous user")
	}
	if token == "" {
		t.Fatalf("expected token to be returned")
	}
	if sessionRepo.sessions[token] == nil {
		t.Fatalf("expected session to be stored")
	}
}

func TestAuthService_Logout(t *testing.T) {
	user := domain.NewUser("Alice", "alice@example.com", false)
	session := domain.NewSession(user.ID, "token-123", time.Now().Add(time.Hour), nil, nil)

	sessionRepo := &stubSessionRepository{sessions: map[string]*domain.Session{
		session.Token: session,
	}}
	userRepo := &stubUserRepository{}
	accountRepo := &stubAccountRepository{}
	service := NewAuthService(userRepo, sessionRepo, accountRepo, testConfig(), newAuthTestLogger())

	if err := service.Logout(context.Background(), session.Token); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !session.Revoked {
		t.Fatalf("expected session to be revoked")
	}
}

func TestAuthService_RegisterUserRepoError(t *testing.T) {
	expectedErr := errors.New("write failed")
	userRepo := &stubUserRepository{err: expectedErr}
	accountRepo := &stubAccountRepository{}
	sessionRepo := &stubSessionRepository{}
	service := NewAuthService(userRepo, sessionRepo, accountRepo, testConfig(), newAuthTestLogger())

	if _, err := service.Register(context.Background(), "Alice", "alice@example.com", "SecurePass123"); err == nil {
		t.Fatalf("expected error when repository fails")
	}
}
