package handlers

import (
	"net/http"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/services"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService services.ProfileService
}

func NewProfileHandler(service services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: service}
}

// ProfilePictureHandler godoc
// @Summary Get profile picture
// @Description Get profile picture
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Profile picture retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /profile/picture [post]
func (h *ProfileHandler) ProfilePictureHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Profile picture retrieved successfully"})
}
