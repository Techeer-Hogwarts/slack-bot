package handlers

import (
	"net/http"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/services"
	"github.com/gin-gonic/gin"
)

type SlackHandler struct {
	service services.SlackService
}

func NewSlackHandler(service services.SlackService) *SlackHandler {
	return &SlackHandler{service: service}
}

// SlackCommandHandler godoc
// @Summary Handle Slack command
// @Description Handle Slack command
// @Tags slack
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Slack command received"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /slack/command [post]
func (h *SlackHandler) SlackCommandHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Slack command received"})
}

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
	c.JSON(http.StatusOK, gin.H{"message": "Slack interaction received"})
}
