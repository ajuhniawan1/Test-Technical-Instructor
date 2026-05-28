package handler

import (
	"net/http"

	"assignment-platform/internal/dto"
	"assignment-platform/internal/service"
	"assignment-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler menangani endpoint auth seperti login dan me.
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login menerima email/password lalu mengembalikan JWT token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// ShouldBindJSON membaca JSON request dan menjalankan validasi tag binding.
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login success", result)
}

// Me mengembalikan data user dari token yang sedang login.
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint64("user_id")

	user, err := h.authService.Me(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Current user profile", user)
}
