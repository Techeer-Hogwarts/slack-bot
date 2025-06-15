package handlers

import (
	"net/http"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
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
// @Security APIKeyAuth
// @Param profilePictureRequest body models.ProfilePictureRequest true "Profile picture request"
// @Success 200 {object} models.ProfilePictureResponse "Profile picture retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /profile/picture [post]
func (h *ProfileHandler) ProfilePictureHandler(c *gin.Context) {
	var req models.ProfilePictureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := h.profileService.GetProfilePicture(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ProfilePictureResponse{
		Email:     req.Email,
		Image:     image,
		IsTecheer: true,
	})
}
