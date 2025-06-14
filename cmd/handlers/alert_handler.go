package handlers

import (
	"net/http"

	"github.com/Techeer-Hogwarts/slack-bot/cmd/models"
	"github.com/Techeer-Hogwarts/slack-bot/cmd/services"
	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	service services.AlertService
}

func NewAlertHandler(service services.AlertService) *AlertHandler {
	return &AlertHandler{service: service}
}

// AlertMessageHandler godoc
// @Summary Send alert message
// @Description Send alert message
// @Tags alert
// @Accept json
// @Produce json
// @Security APIKeyAuth
// @Param models.AlertMessageSchema body models.AlertMessageSchema true "AlertMessageSchema"
// @Success 200 {object} map[string]interface{} "Alert message sent"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /alert/message [post]
func (h *AlertHandler) AlertMessageHandler(c *gin.Context) {
	var alertMessage models.AlertMessageSchema
	if err := c.ShouldBindJSON(&alertMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SendAlert(alertMessage.ChannelID, alertMessage.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Alert message Success"})
}

// AlertFindMemberHandler godoc
// @Summary Send Message to Find member
// @Description Send Message to Find member
// @Tags alert
// @Accept json
// @Produce json
// @Security APIKeyAuth
// @Param models.FindMemberSchema body models.FindMemberSchema true "FindMemberSchema"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /alert/find-member [post]
func (h *AlertHandler) AlertFindMemberHandler(c *gin.Context) {
	var alertMessage models.FindMemberSchema
	if err := c.ShouldBindJSON(&alertMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SendAlertToFindMember(alertMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post to find_member channel Success"})
}

// AlertUserMessageHandler godoc
// @Summary Send user message (legacy)
// @Description Send user message (legacy)
// @Tags alert
// @Accept json
// @Produce json
// @Security APIKeyAuth
// @Param models.UserMessageSchema body models.UserMessageSchema true "UserMessageSchema"
// @Success 200 {object} map[string]interface{} "User message sent"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /alert/user [post]
func (h *AlertHandler) AlertUserMessageHandler(c *gin.Context) {
	var alertMessage models.UserMessageSchema
	if err := c.ShouldBindJSON(&alertMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SendAlertToUser(alertMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User message Sent"})
}

// AlertChannelMessageHandler godoc
// @Summary Send channel message to find_member channel (legacy)
// @Description Send channel message to find_member channel (legacy)
// @Tags alert
// @Accept json
// @Produce json
// @Security APIKeyAuth
// @Param models.FindMemberSchema body models.FindMemberSchema true "FindMemberSchema"
// @Success 200 {object} map[string]interface{} "Channel message sent"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /alert/channel [post]
func (h *AlertHandler) AlertChannelMessageHandler(c *gin.Context) {
	var alertMessage models.FindMemberSchema
	if err := c.ShouldBindJSON(&alertMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SendAlertToFindMember(alertMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Channel message Sent"})
}
