package handlers

import (
	"log"
	"net/http"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/services"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
)

type SlackHandler struct {
	service       services.SlackService
	deployService services.DeployService
}

func NewSlackHandler(service services.SlackService, deployService services.DeployService) *SlackHandler {
	return &SlackHandler{service: service, deployService: deployService}
}

// // SlackCommandHandler godoc
// // @Summary Handle Slack command
// // @Description Handle Slack command
// // @Tags slack
// // @Accept json
// // @Produce json
// // @Success 200 {object} map[string]interface{} "Slack command received"
// // @Failure 400 {object} map[string]interface{} "Bad request"
// // @Router /slack/command [post]
// func (h *SlackHandler) SlackCommandHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "Slack command received"})
// }

// SlackInteractionHandler godoc
// @Summary Handle Slack interaction
// @Description Handle Slack interaction
// @Tags slack
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Slack interaction received"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /slack/interaction [post]
func (h *SlackHandler) SlackInteractionHandler(c *gin.Context) {
	payloadStr := c.PostForm("payload")
	if payloadStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload is required"})
		return
	}

	var payload slack.InteractionCallback
	err := payload.UnmarshalJSON([]byte(payloadStr))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode interaction payload"})
		return
	}

	if payload.Type == slack.InteractionTypeBlockActions {
		log.Printf("payload: %+v", payload.BlockActionState.Values)
		// h.deployService.TriggerDeployment(payload.BlockActionState.Values["deploy_button"]["image_name"].Value, payload)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slack interaction received"})
}
