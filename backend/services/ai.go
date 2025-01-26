package services

import (
	"ai-content-creation/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AIService struct {
	cloudflareAccountID string
	cloudflareAPIToken  string
}

type CloudflareAIRequest struct {
	Messages []Message `json:"messages"`
	Stream   bool     `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CloudflareAIResponse struct {
	Result struct {
		Response string `json:"response"`
	} `json:"result"`
	Success bool   `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func NewAIService() *AIService {
	return &AIService{
		cloudflareAccountID: os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		cloudflareAPIToken:  os.Getenv("CLOUDFLARE_API_TOKEN"),
	}
}

func (ai *AIService) GenerateContent(c *gin.Context, contentReq *models.ContentRequest) (string, error) {
	var modelEndpoint string
	switch contentReq.AIModel {
	case "mistral-7b":
		modelEndpoint = "@cf/mistralai/mistral-7b-instruct-v0.1"
	case "llama2-7b":
		modelEndpoint = "@cf/meta/llama-2-7b-chat-fp16"
	default:
		modelEndpoint = "@cf/meta/llama-2-7b-chat-fp16"
	}

	aiReq := CloudflareAIRequest{
		Messages: []Message{
			{Role: "system", Content:`
          You are a marketing and sales professional who is looking to increase your sales and the best in the industry for 
        growing local brands.`},
      {Role: "user", Content: contentReq.Prompt},
		},
		Stream: false,
	}

	reqBody, err := json.Marshal(aiReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s",
		ai.cloudflareAccountID, modelEndpoint)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+ai.cloudflareAPIToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var cloudflareResponse CloudflareAIResponse
	if err := json.Unmarshal(body, &cloudflareResponse); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %v", err)
	}

	// Check for API errors
	if !cloudflareResponse.Success && len(cloudflareResponse.Errors) > 0 {
		return "", fmt.Errorf("API error: %s", cloudflareResponse.Errors[0].Message)
	}

	return cloudflareResponse.Result.Response, nil
}
