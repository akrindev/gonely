package http

import (
	"github.com/akrindev/gonely/internal/application/auth"
	"github.com/akrindev/gonely/internal/application/feeds"
	"github.com/akrindev/gonely/internal/infrastructure/config"
	"github.com/akrindev/gonely/internal/infrastructure/logger"
	"github.com/akrindev/gonely/internal/interfaces/http/handler"
	"github.com/akrindev/gonely/internal/interfaces/http/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter sets up the HTTP router with all routes
func SetupRouter(
	cfg *config.Config,
	log *logger.Logger,
	authService *auth.AuthService,
	feedService *feeds.FeedService,
) *gin.Engine {
	// Set Gin mode
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware(log))
	router.Use(middleware.ErrorHandler(log))

	// CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	feedHandler := handler.NewFeedHandler(feedService)

	// Base routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": cfg.App.Name,
		})
	})

	// Swagger routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group(cfg.App.BasePath)
	{
		// Auth routes (no auth required)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/anonymous", authHandler.CreateAnonymous)
		}

		// Protected auth routes
		authProtected := api.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(authService))
		{
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.GET("/me", authHandler.Me)
		}

		// V1 API routes
		v1 := api.Group("/v1")
		{
			// Public posts (read-only)
			v1.GET("/posts", feedHandler.GetPosts)
			v1.GET("/posts/:id", feedHandler.GetPostByID)
			v1.GET("/posts/:id/comments", feedHandler.GetComments)

			// Protected posts (write operations)
			postsProtected := v1.Group("/posts")
			postsProtected.Use(middleware.AuthMiddleware(authService))
			{
				postsProtected.POST("", feedHandler.CreatePost)
				postsProtected.POST("/:id/comments", feedHandler.AddComment)
			}
		}
	}

	return router
}
