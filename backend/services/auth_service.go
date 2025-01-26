package services

import (
	"ai-content-creation/models"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) RegisterUser(name, email, password string) (*models.User, string, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, "", errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user
	user := models.User{
		UserID:           strings.Replace(uuid.New().String(), "-", "", -1),
		Name:             name,
		Email:            email,
		Password:         string(hashedPassword),
		SubscriptionTier: models.FreeTier,
		RemainingCredits: 200,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create user: %v", err)
	}

	// Generate token
	token, err := s.generateToken(user.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return &user, token, nil
}

func (s *AuthService) LoginUser(email, password string) (*models.User, string, error) {
	// Find user
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate token
	token, err := s.generateToken(user.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return &user, token, nil
}

func (s *AuthService) generateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat": time.Now().Unix(),
	})

	// Get JWT secret from environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-256-bit-secret" // Default secret for development
	}

	return token.SignedString([]byte(secret))
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "your-256-bit-secret" // Default secret for development
		}

		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}

	return "", errors.New("invalid token claims")
}

// AuthMiddleware is a Gin middleware to validate JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		authService := c.MustGet("authService").(*AuthService)
		userID, err := authService.ValidateToken(tokenParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set user ID in context
		c.Set("userID", userID)
		c.Next()
	}
} 
