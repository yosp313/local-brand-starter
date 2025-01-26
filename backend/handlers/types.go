package handlers

import (
	"ai-content-creation/services"

	"github.com/gin-gonic/gin"
)

// Response is a generic response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string     `json:"error,omitempty"`
}

// Handler wraps all services needed by handlers
type Handler struct {
	authService         *services.AuthService
	userService        *services.UserService
	contentService     *services.ContentService
	subscriptionService *services.SubscriptionService
}

// NewHandler creates a new handler instance
func NewHandler(
	authService *services.AuthService,
	userService *services.UserService,
	contentService *services.ContentService,
	subscriptionService *services.SubscriptionService,
) *Handler {
	return &Handler{
		authService:         authService,
		userService:        userService,
		contentService:     contentService,
		subscriptionService: subscriptionService,
	}
}

// sendError sends an error response
func sendError(c *gin.Context, code int, err string) {
	c.JSON(code, Response{
		Success: false,
		Error:   err,
	})
}

// sendSuccess sends a success response
func sendSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Success: true,
		Data:    data,
	})
} 
