package main

import (
	"ai-content-creation/handlers"
	"ai-content-creation/models"
	"ai-content-creation/services"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Database connection
	db, err := gorm.Open(sqlite.Open("database.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize database schema
	if err := models.InitDB(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	authService := services.NewAuthService(db)
	userService := services.NewUserService(db)
	contentService := services.NewContentService(db)
	subscriptionService := services.NewSubscriptionService(db)

	// Initialize handlers
	h := handlers.NewHandler(authService, userService, contentService, subscriptionService)

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("FRONTEND_URL")} // Add your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Make auth service available to middleware
	r.Use(func(c *gin.Context) {
		c.Set("authService", authService)
		c.Next()
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Health check (no auth required)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		// Auth routes (no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
		}

		// Protected routes (auth required)
		protected := api.Group("")
		protected.Use(services.AuthMiddleware())
		{
			// Content generation endpoints
			protected.POST("/generate", h.GenerateContent)
			protected.GET("/content", h.GetContent)
			protected.GET("/content/:id", h.GetContentByID)

			// Subscription plan endpoints
			protected.GET("/subscription-plans", h.GetSubscriptionPlans)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
