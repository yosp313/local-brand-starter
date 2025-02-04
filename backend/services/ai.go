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
	Stream   bool      `json:"stream"`
}

type ImageRequest struct {
	Prompt string `json:"prompt"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CloudflareAIResponse struct {
	Result struct {
		Response string `json:"response"`
	} `json:"result"`
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ImageResponse struct {
	Result struct {
		Image []byte `json:"image"`
	}
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	}
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
			{Role: "system", Content: `
          You are a marketing and sales professional who is looking to increase your sales and the best in the industry for 
        growing local brands.`},
			{Role: "user", Content: contentReq.Prompt},
		},
		Stream: false,
	}

	reqBody, err := json.Marshal(aiReq)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal the request into json: %s", err)
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

func (ai *AIService) GenerateImage(c *gin.Context, contentReq *models.ContentRequest) (string, error) {
	modelEndpoint := "@cf/black-forest-labs/flux-1-schnell"

	imageReq := ImageRequest{
		Prompt: "You are a senior in digital marketing and your task is to Generate a professional digital illustration of the following prompt: " + contentReq.Prompt,
	}

	reqBody, err := json.Marshal(imageReq)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal the request into json: %s", err)
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s",
		ai.cloudflareAccountID, modelEndpoint)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+ai.cloudflareAPIToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var imageResponse ImageResponse
	if err := json.Unmarshal(respBody, &imageResponse); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %v", err)
	}

	if !imageResponse.Success && len(imageResponse.Errors) > 0 {
		return "", fmt.Errorf("Image API error: %s", imageResponse.Errors[0].Message)
	}

	if err := UploadImageToS3(imageResponse.Result.Image, contentReq); err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	imageURL := GetImageURL(contentReq.RequestID)

	return imageURL, nil
}
