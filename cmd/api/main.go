// Package main Gonely API
//
//	@title			Gonely API
//	@version		1.0
//	@description	A social media API built with Go and Gin
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host		localhost:8080
//	@BasePath	/api
//
//	@securityDefinitions.apikey BearerAuth
//	@description				JWT Authorization header using the Bearer scheme.
//	@in							header
//	@name						Authorization
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/akrindev/gonely/docs"
	authApp "github.com/akrindev/gonely/internal/application/auth"
	feedsApp "github.com/akrindev/gonely/internal/application/feeds"
	"github.com/akrindev/gonely/internal/infrastructure/config"
	"github.com/akrindev/gonely/internal/infrastructure/database"
	"github.com/akrindev/gonely/internal/infrastructure/logger"
	authRepo "github.com/akrindev/gonely/internal/infrastructure/repository/auth"
	feedsRepo "github.com/akrindev/gonely/internal/infrastructure/repository/feeds"
	httpInterface "github.com/akrindev/gonely/internal/interfaces/http"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting application",
		zap.String("name", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.App.Port),
	)

	// Connect to database
	db, err := database.New(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Run migrations
	logger.Info("Running database migrations...")
	if err := db.Migrate(); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}
	logger.Info("Database migrations completed successfully")

	// Seed database (development only)
	if cfg.App.IsDevelopment() {
		logger.Info("Seeding database with sample data...")
		if err := db.Seed(); err != nil {
			logger.Warn("Failed to seed database", zap.Error(err))
		}
	}

	// Initialize repositories
	userRepo := authRepo.NewUserRepository(db.DB)
	accountRepo := authRepo.NewAccountRepository(db.DB)
	sessionRepo := authRepo.NewSessionRepository(db.DB)
	postRepo := feedsRepo.NewPostRepository(db.DB)
	commentRepo := feedsRepo.NewCommentRepository(db.DB)

	// Initialize services
	authService := authApp.NewAuthService(userRepo, sessionRepo, accountRepo, cfg, logger)
	feedService := feedsApp.NewFeedService(postRepo, commentRepo, logger)

	// Setup HTTP router
	router := httpInterface.SetupRouter(cfg, logger, authService, feedService)

	// Setup HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("HTTP server starting", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info(fmt.Sprintf("Server is ready to handle requests at %s", cfg.App.URL))

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server is shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
