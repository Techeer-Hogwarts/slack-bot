package handlers

import (
	"net/http"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

// LoginHandler godoc
// @Summary Login user
// @Description Login to the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserLoginRequest true "User Login Request Body"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /auth/login [post]
func (h *UserHandler) LoginHandler(c *gin.Context) {
	var req models.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("access_token", token, 6*3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// ResetPasswordHandler godoc
// @Summary Reset user password
// @Description Allows authenticated users to change their password by providing the current and new password.
// @Tags auth
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param request body models.PasswordResetRequest true "Password Reset Payload"
// @Success 200 {object} map[string]string "Password updated successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /auth/reset [post]
func (h *UserHandler) ResetPasswordHandler(c *gin.Context) {
	var req models.PasswordResetRequest

	// Extract user_id from access_token (assumed to be set in middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.ResetPassword(userID.(int), req.CurrentPass, req.NewPass)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
