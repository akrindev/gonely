package database

import (
	"fmt"
	"time"

	"github.com/akrindev/gonely/internal/domain/auth"
	"github.com/akrindev/gonely/internal/domain/feeds"
	"github.com/akrindev/gonely/internal/infrastructure/config"
	"github.com/akrindev/gonely/internal/infrastructure/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Database wraps the GORM DB instance
type Database struct {
	*gorm.DB
}

// New creates a new database connection
func New(cfg *config.Config, log *logger.Logger) (*Database, error) {
	// Configure GORM logger
	var gormLogLevel gormlogger.LogLevel
	switch cfg.Log.Level {
	case "debug":
		gormLogLevel = gormlogger.Info
	case "error":
		gormLogLevel = gormlogger.Error
	default:
		gormLogLevel = gormlogger.Warn
	}

	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	}

	dsn := cfg.Database.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info("Database connection established",
		zap.String("host", cfg.Database.Host),
		zap.String("database", cfg.Database.DBName),
	)

	return &Database{db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate runs database migrations
func (d *Database) Migrate() error {
	// Migrate auth domain models
	if err := d.AutoMigrate(
		&auth.User{},
		&auth.Account{},
		&auth.Session{},
		&auth.Verification{},
	); err != nil {
		return fmt.Errorf("failed to migrate auth models: %w", err)
	}

	// Migrate feeds domain models
	if err := d.AutoMigrate(
		&feeds.Post{},
		&feeds.Comment{},
	); err != nil {
		return fmt.Errorf("failed to migrate feeds models: %w", err)
	}

	return nil
}

// Seed seeds the database with initial data (for development/testing)
func (d *Database) Seed() error {
	// Check if we already have data
	var count int64
	if err := d.Model(&feeds.Post{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // Already seeded
	}

	// Create sample posts
	posts := []*feeds.Post{
		feeds.NewPost("Welcome to Gonely! This is a sample post."),
		feeds.NewPost("Another sample post to demonstrate the feeds functionality."),
		feeds.NewPost("Learning Golang with DDD principles is awesome!"),
	}

	for _, post := range posts {
		if err := d.Create(post).Error; err != nil {
			return fmt.Errorf("failed to seed post: %w", err)
		}
	}

	return nil
}
