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

// @Summary Login
// @Description Login to the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserLoginRequest true "User Login Request Body"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
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
