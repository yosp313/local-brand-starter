package services

import (
	"ai-content-creation/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContentService struct {
	db        *gorm.DB
	aiService *AIService
}

func NewContentService(db *gorm.DB) *ContentService {
	return &ContentService{
		db:        db,
		aiService: NewAIService(),
	}
}

func (s *ContentService) Generate(c *gin.Context, userID string, model string, prompt string) (*models.GeneratedContent, error) {
	// Check user's subscription and credits
	var user models.User
	if err := s.db.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.RemainingCredits <= 0 {
		return nil, fmt.Errorf("no remaining credits")
	}

	// Create content request
	contentReq := models.ContentRequest{
		RequestID: uuid.New().String(),
		UserID:    userID,
		AIModel:   model,
		Prompt:    prompt,
		Status:    "pending",
	}

	if err := s.db.Create(&contentReq).Error; err != nil {
		return nil, fmt.Errorf("failed to create content request: %v", err)
	}

	// Generate content
	response, err := s.aiService.GenerateContent(c, &contentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	// Create generated content
	generatedContent := &models.GeneratedContent{
		ContentID: uuid.New().String(),
		RequestID: contentReq.RequestID,
		Output:    response,
		Version:   1,
	}

	if err := s.db.Create(generatedContent).Error; err != nil {
		return nil, fmt.Errorf("failed to create generated content: %v", err)
	}

	// Deduct credits
	if err := s.db.Model(&user).Update("remaining_credits", user.RemainingCredits-10).Error; err != nil {
		return nil, fmt.Errorf("failed to update credits: %v", err)
	}

	return generatedContent, nil
}

func (s *ContentService) GetUserContent(userID string) ([]models.GeneratedContent, error) {
	var content []models.GeneratedContent
	err := s.db.Joins("JOIN content_requests ON content_requests.request_id = generated_contents.request_id").
		Where("content_requests.user_id = ?", userID).
		Find(&content).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content: %v", err)
	}
	return content, nil
}

func (s *ContentService) GetContentByID(userID string, contentID string) (*models.GeneratedContent, error) {
	var content models.GeneratedContent
	err := s.db.Joins("JOIN content_requests ON content_requests.request_id = generated_contents.request_id").
		Where("content_requests.user_id = ? AND generated_contents.content_id = ?", userID, contentID).
		First(&content).Error
	if err != nil {
		return nil, fmt.Errorf("content not found: %v", err)
	}
	return &content, nil
}
