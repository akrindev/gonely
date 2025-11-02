package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/akrindev/gonely/internal/domain/auth"
	"github.com/akrindev/gonely/internal/infrastructure/config"
	"github.com/akrindev/gonely/internal/infrastructure/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication-related business logic
type AuthService struct {
	userRepo    auth.UserRepository
	sessionRepo auth.SessionRepository
	accountRepo auth.AccountRepository
	config      *config.Config
	logger      *logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo auth.UserRepository,
	sessionRepo auth.SessionRepository,
	accountRepo auth.AccountRepository,
	cfg *config.Config,
	log *logger.Logger,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		accountRepo: accountRepo,
		config:      cfg,
		logger:      log.Named("auth-service"),
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, name, email, password string) (*auth.User, error) {
	s.logger.Info("Registering new user", zap.String("email", email))

	// Check if user already exists
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Create user
	user := auth.NewUser(name, email, false)
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create account
	account := auth.NewAccount(user.Email, "credentials", user.ID)
	account.SetPassword(string(hashedPassword))
	if err := s.accountRepo.Create(ctx, account); err != nil {
		s.logger.Error("Failed to create account", zap.Error(err))
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	s.logger.Info("User registered successfully", zap.String("user_id", user.ID))
	return user, nil
}

// Login authenticates a user and creates a session
func (s *AuthService) Login(ctx context.Context, email, password, userAgent, ipAddress string) (string, *auth.User, error) {
	s.logger.Info("User login attempt", zap.String("email", email))

	// Find user
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Find account
	accounts, err := s.accountRepo.FindByUserID(ctx, user.ID)
	if err != nil || len(accounts) == 0 {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	account := accounts[0]
	if account.Password == nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Generate token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create session
	expiresAt := time.Now().Add(s.config.Auth.SessionExpiry)
	session := auth.NewSession(user.ID, token, expiresAt, &userAgent, &ipAddress)
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		s.logger.Error("Failed to create session", zap.Error(err))
		return "", nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info("User logged in successfully", zap.String("user_id", user.ID))
	return token, user, nil
}

// ValidateToken validates a JWT token and returns the session
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*auth.Session, error) {
	// Find session by token
	session, err := s.sessionRepo.FindByToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	// Check if session is valid
	if !session.IsValid() {
		return nil, fmt.Errorf("session is invalid or expired")
	}

	// Update last used timestamp
	session.UpdateLastUsed()
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		s.logger.Warn("Failed to update session last used", zap.Error(err))
	}

	return session, nil
}

// Logout revokes a session
func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	s.logger.Info("User logout")

	session, err := s.sessionRepo.FindByToken(ctx, tokenString)
	if err != nil || session == nil {
		return fmt.Errorf("session not found")
	}

	session.Revoke("user logout")
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		s.logger.Error("Failed to revoke session", zap.Error(err))
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	return nil
}

// CreateAnonymousUser creates an anonymous user
func (s *AuthService) CreateAnonymousUser(ctx context.Context) (*auth.User, string, error) {
	s.logger.Info("Creating anonymous user")

	name := fmt.Sprintf("Anonymous_%s", uuid.New().String()[:8])
	email := fmt.Sprintf("anon_%s@anonymous.local", uuid.New().String())

	user := auth.NewUser(name, email, true)
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create anonymous user", zap.Error(err))
		return nil, "", fmt.Errorf("failed to create anonymous user: %w", err)
	}

	// Generate token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Create session
	expiresAt := time.Now().Add(s.config.Auth.SessionExpiry)
	session := auth.NewSession(user.ID, token, expiresAt, nil, nil)
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		s.logger.Error("Failed to create session", zap.Error(err))
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	return user, token, nil
}

// generateToken generates a JWT token
func (s *AuthService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.config.Auth.JWTExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Auth.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
