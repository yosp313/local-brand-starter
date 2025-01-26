package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenerateContentRequest struct {
	Model  string `json:"model" binding:"required"`
	Prompt string `json:"prompt" binding:"required"`
}

type ContentResponse struct {
	ContentID string `json:"content_id"`
	RequestID string `json:"request_id"`
	Output    string `json:"output"`
	Version   int    `json:"version"`
}

func (h *Handler) GenerateContent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req GenerateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	content, err := h.contentService.Generate(c, userID.(string), req.Model, req.Prompt)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, http.StatusOK, ContentResponse{
		ContentID: content.ContentID,
		RequestID: content.RequestID,
		Output:    content.Output,
		Version:   content.Version,
	})
}

func (h *Handler) GetContent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	content, err := h.contentService.GetUserContent(userID.(string))
	if err != nil {
		sendError(c, http.StatusInternalServerError, "Failed to fetch content")
		return
	}

	var response []ContentResponse
	for _, c := range content {
		response = append(response, ContentResponse{
			ContentID: c.ContentID,
			RequestID: c.RequestID,
			Output:    c.Output,
			Version:   c.Version,
		})
	}

	sendSuccess(c, http.StatusOK, response)
}

func (h *Handler) GetContentByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	contentID := c.Param("id")
	content, err := h.contentService.GetContentByID(userID.(string), contentID)
	if err != nil {
		sendError(c, http.StatusNotFound, "Content not found")
		return
	}

	sendSuccess(c, http.StatusOK, ContentResponse{
		ContentID: content.ContentID,
		RequestID: content.RequestID,
		Output:    content.Output,
		Version:   content.Version,
	})
} 
