// package server for initializing server and setting up routes
package server

import (
	"user-mgmt/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Server struct set config, db and logger
type Server struct {
	config *config.Config
	db     *mongo.Database
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, db *mongo.Database) *Server {
	return &Server{config: cfg, db: db}
}

// SetupRoutes sets up the routes for the server
func (s *Server) SetupRoutes() *gin.Engine {

	// Init Router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// Init db

	v1 := router.Group("/api/v1")
	{
		// auth := v1.Group("/auth")
		// {
		// 	// auth.POST("/register", s.register)
		// 	// auth.POST("/login", s.login)
		// 	// auth.POST("/refresh", s.refreshToken)
		// 	// auth.POST("/logout", s.logout)
		// }

		// Protected routes (authenticated users)
		protected := v1.Group("/")
		protected.Use(s.authMiddleware())
		{
			// User routes
			// userRoute := protected.Group("/user")
			// {
			// 	// userRoute.GET("/:id", s.getUser)
			// 	// userRoute.GET("/", s.getUserList)
			// 	// userRoute.PATCH("/:id", s.updateUser)
			// 	// userRoute.DELETE("/:id", s.deleteUser)
			// }
		}
	}

	// Health routes
	router.GET("/health", func(c *gin.Context) {
		// Get host name
		c.JSON(200, gin.H{"status": "healthy"})
	})

	return router
}

// corsMiddleware for CORS frontend in different origin will can call this server API
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
