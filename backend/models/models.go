package models

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionTier represents the available subscription tiers
type SubscriptionTier string

const (
	FreeTier       SubscriptionTier = "free"
	ProTier        SubscriptionTier = "pro"
	EnterpriseTier SubscriptionTier = "enterprise"
)

// SubscriptionPlan represents the available plans and their features
type SubscriptionPlan struct {
	gorm.Model
	PlanID          string           `gorm:"type:string;uniqueIndex" json:"plan_id"`
	Tier            SubscriptionTier `gorm:"type:string;uniqueIndex" json:"tier"`
	Name            string           `json:"name"`
	Price           float64          `json:"price"`
	TokensPerMonth  int              `json:"tokens_per_month"`
	ModelsAvailable string           `json:"models_available"` // JSON string array
}

// SetModelsAvailable converts string slice to JSON string for storage
func (sp *SubscriptionPlan) SetModelsAvailable(models []string) error {
	data, err := json.Marshal(models)
	if err != nil {
		return err
	}
	sp.ModelsAvailable = string(data)
	return nil
}

// GetModelsAvailable converts stored JSON string to string slice
func (sp *SubscriptionPlan) GetModelsAvailable() ([]string, error) {
	var models []string
	if err := json.Unmarshal([]byte(sp.ModelsAvailable), &models); err != nil {
		return nil, err
	}
	return models, nil
}

type User struct {
	gorm.Model
	UserID           string           `gorm:"type:string;uniqueIndex" json:"user_id"`
	Name             string           `json:"name"`
	Email            string           `gorm:"uniqueIndex" json:"email"`
	Password         string           `json:"-"` // "-" means this field won't be included in JSON
	SubscriptionTier SubscriptionTier `gorm:"type:string;default:'free'" json:"subscription_tier"`
	StripeCustomerID string           `json:"stripe_customer_id,omitempty"`
	RemainingCredits int              `gorm:"default:1000" json:"remaining_credits"`
}

type ContentRequest struct {
	gorm.Model
	RequestID string `gorm:"type:string;uniqueIndex" json:"request_id"`
	UserID    string `gorm:"type:string" json:"user_id"`
	AIModel   string `json:"model"` // mistral-7b or llama2-7b
	Prompt    string `json:"prompt"`
	Status    string `gorm:"default:'pending'" json:"status"`
}

type GeneratedContent struct {
	gorm.Model
	ContentID string `gorm:"type:string;uniqueIndex" json:"content_id"`
	RequestID string `gorm:"type:string" json:"request_id"`
	Output    string `json:"output"`
	Version   int    `gorm:"default:1" json:"version"`
	CacheKey  string `json:"cache_key"`
}

func InitDB(db *gorm.DB) error {
	// Auto-migrate the schemas
	if err := db.AutoMigrate(&User{}, &ContentRequest{}, &GeneratedContent{}, &SubscriptionPlan{}); err != nil {
		return err
	}

	// Initialize default subscription plans if they don't exist
	plans := []SubscriptionPlan{
		{
			PlanID:         uuid.New().String(),
			Tier:           FreeTier,
			Name:           "Free",
			Price:          0,
			TokensPerMonth: 10000,
		},
		{
			PlanID:         uuid.New().String(),
			Tier:           ProTier,
			Name:           "Pro",
			Price:          29.99,
			TokensPerMonth: 500000,
		},
		{
			PlanID:         uuid.New().String(),
			Tier:           EnterpriseTier,
			Name:           "Enterprise",
			Price:          999.99,
			TokensPerMonth: -1, // Unlimited
		},
	}

	// Set available models for each plan
	if err := plans[0].SetModelsAvailable([]string{"llama2-7b"}); err != nil {
		return err
	}
	if err := plans[1].SetModelsAvailable([]string{"llama2-7b", "mistral-7b"}); err != nil {
		return err
	}
	if err := plans[2].SetModelsAvailable([]string{"llama2-7b", "mistral-7b"}); err != nil {
		return err
	}

	for _, plan := range plans {
		// Check if plan exists
		var existingPlan SubscriptionPlan
		if err := db.Where("tier = ?", plan.Tier).First(&existingPlan).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new plan if it doesn't exist
				if err := db.Create(&plan).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
