package services

import (
	"ai-content-creation/models"
	"fmt"

	"gorm.io/gorm"
)

type SubscriptionService struct {
	db *gorm.DB
}

func NewSubscriptionService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}

func (s *SubscriptionService) GetPlans() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	if err := s.db.Find(&plans).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch subscription plans: %v", err)
	}
	return plans, nil
}

func (s *SubscriptionService) GetPlanByTier(tier models.SubscriptionTier) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	if err := s.db.Where("tier = ?", tier).First(&plan).Error; err != nil {
		return nil, fmt.Errorf("plan not found: %v", err)
	}
	return &plan, nil
} 