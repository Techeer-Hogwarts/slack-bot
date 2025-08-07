package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/services"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/slackmessages"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
)

type SlackHandler struct {
	slackService  services.SlackService
	deployService services.DeployService
}

func NewSlackHandler(slackService services.SlackService, deployService services.DeployService) *SlackHandler {
	return &SlackHandler{slackService: slackService, deployService: deployService}
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
// @Security SlackSigningSecret
// @Success 200 {object} map[string]interface{} "Slack interaction received"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /slack/interactions [post]
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
		action := payload.ActionCallback.BlockActions[0]

		switch action.ActionID {
		case "deploy_button":
			err := h.deployService.TriggerDeployment(action.Value, payload)
			if err != nil {
				log.Printf("Failed to trigger deployment: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to trigger deployment"})
				return
			}
		case "no_deploy_button":
			err := h.slackService.DeleteMessage(payload.Channel.ID, payload.Message.Timestamp)
			if err != nil {
				log.Printf("Failed to delete message: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
				return
			}
		case "delete_button":
			userID := payload.User.ID
			userIDButtonValue := action.Value
			var profileIDs []string
			json.Unmarshal([]byte(userIDButtonValue), &profileIDs)
			if !slackmessages.CheckUserIsAllowed(profileIDs, userID) {
				log.Printf("User is not allowed to delete message")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not allowed to delete message"})
				return
			}
			log.Printf("User is allowed to delete message")
			err := h.slackService.DeleteMessage(payload.Channel.ID, payload.Message.Timestamp)
			if err != nil {
				log.Printf("Failed to delete message: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
				return
			}
		case "redirect_button":
			log.Printf("Redirect action: %s", action.ActionID)
			return
		default:
			log.Printf("Unknown action: %s", action.ActionID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown action"})
			return
		}
	}
	log.Printf("Action: %+v success", payload.ActionCallback.BlockActions[0].ActionID)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
