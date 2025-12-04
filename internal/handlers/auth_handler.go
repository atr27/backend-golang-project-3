package handlers

import (
	"hr-backend/internal/models"
	"hr-backend/internal/services"
	"hr-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, 401, "AUTHENTICATION_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Login successful", response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	utils.SuccessResponse(c, 200, "Logged out successfully", nil)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	err := h.authService.ChangePassword(userID.(uint), &req)
	if err != nil {
		utils.ErrorResponse(c, 400, "CHANGE_PASSWORD_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Password changed successfully", nil)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	profile := gin.H{
		"id":    userID,
		"email": email,
		"role":  role,
	}

	utils.SuccessResponse(c, 200, "Profile retrieved successfully", profile)
}
