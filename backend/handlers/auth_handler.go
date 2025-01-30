package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, token, err := h.authService.RegisterUser(req.Name, req.Email, req.Password)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, http.StatusCreated, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.UserID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, token, err := h.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		sendError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	sendSuccess(c, http.StatusOK, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.UserID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

