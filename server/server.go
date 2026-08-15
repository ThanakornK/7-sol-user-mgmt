// package server for initializing server and setting up routes
package server

import (
	"net/http"
	"user-mgmt/config"
	"user-mgmt/handler"
	"user-mgmt/repository/mongorepo"
	"user-mgmt/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Server struct set config, db and logger
type Server struct {
	config *config.Config
	db     *mongo.Database
	logger zerolog.Logger
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, db *mongo.Database, logger zerolog.Logger) *Server {
	return &Server{config: cfg, db: db, logger: logger}
}

// SetupRoutes sets up the routes for the server
func (s *Server) SetupRoutes() *gin.Engine {

	// Init Router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())
	router.Use(traceMiddleware())

	// Init repositories
	userRepository := mongorepo.NewUserRepository(s.db)
	userSessionRepository := mongorepo.NewUserSessionRepository(s.db)

	// Init services
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(userRepository, userSessionRepository, &s.config.JWT, s.config.Refresh)

	// Init handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService, s.config.Refresh)

	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)

		// Register route should not require authentication
		v1.POST("/users", userHandler.CreateUser)
		// User management routes
		users := v1.Group("/users")
		users.Use(s.authMiddleware())
		{

			users.GET("/:id", userHandler.GetUserByID)
			users.GET("", userHandler.GetUserList)
			users.PATCH("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

	}

	// Health routes
	router.GET("/health", func(c *gin.Context) {
		// Get host name
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	return router
}

// corsMiddleware for CORS frontend in different origin will can call this server API
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && origin == s.config.Server.AllowedOrigin {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Trace-ID")
		c.Header("Access-Control-Expose-Headers", "X-Trace-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
