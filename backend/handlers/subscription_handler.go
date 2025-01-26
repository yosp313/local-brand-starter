package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubscriptionPlanResponse struct {
	PlanID          string   `json:"plan_id"`
	Tier            string   `json:"tier"`
	Name            string   `json:"name"`
	Price           float64  `json:"price"`
	TokensPerMonth  int      `json:"tokens_per_month"`
	ModelsAvailable []string `json:"models_available"`
}

func (h *Handler) GetSubscriptionPlans(c *gin.Context) {
	plans, err := h.subscriptionService.GetPlans()
	if err != nil {
		sendError(c, http.StatusInternalServerError, "Failed to fetch subscription plans")
		return
	}

	var response []SubscriptionPlanResponse
	for _, plan := range plans {
		models, err := plan.GetModelsAvailable()
		if err != nil {
			sendError(c, http.StatusInternalServerError, "Failed to parse models available")
			return
		}

		response = append(response, SubscriptionPlanResponse{
			PlanID:          plan.PlanID,
			Tier:            string(plan.Tier),
			Name:            plan.Name,
			Price:           plan.Price,
			TokensPerMonth:  plan.TokensPerMonth,
			ModelsAvailable: models,
		})
	}

	sendSuccess(c, http.StatusOK, response)
} 
