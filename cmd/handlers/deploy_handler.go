package handlers

import (
	"log"
	"net/http"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/services"
	"github.com/gin-gonic/gin"
)

type DeployHandler struct {
	deployService services.DeployService
}

func NewDeployHandler(service services.DeployService) *DeployHandler {
	return &DeployHandler{
		deployService: service,
	}
}

// DeployImageHandler godoc
// @Summary Deploy image
// @Description Deploy image
// @Tags deploy
// @Accept json
// @Produce json
// @Security APIKeyAuth
// @Param deployRequest body models.ImageDeployRequest true "Deployment request"
// @Success 200 {object} map[string]interface{} "Deployment request received"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /deploy/image [post]
func (h *DeployHandler) DeployImageHandler(c *gin.Context) {
	apiKey, _ := c.Get("valid_api_key")
	if apiKey != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.ImageDeployRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.deployService.SendDeploymentMessage(req); err != nil {
		log.Printf("Failed to send deployment message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process deployment request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment request received"})
}

// DeployStatusHandler godoc
// @Summary Deploy status
// @Description Deploy status
// @Tags deploy
// @Accept json
// @Produce json
// @Security APIKeyAuth
// @Param statusRequest body models.StatusRequest true "Status request"
// @Success 200 {object} map[string]interface{} "Status request received"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /deploy/status [post]
func (h *DeployHandler) DeployStatusHandler(c *gin.Context) {
	apiKey, _ := c.Get("valid_api_key")
	if apiKey != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var status models.StatusRequest
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.deployService.SendDeploymentStatus(status); err != nil {
		log.Printf("Failed to send deployment status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process status update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status update received"})
}
