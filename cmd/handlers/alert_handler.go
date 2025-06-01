package handlers

import (
	"net/http"

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
// @Success 200 {object} map[string]interface{} "Alert message sent"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /alert/message [post]
func (h *AlertHandler) AlertMessageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Alert message received"})
}

// AlertFindMemberHandler godoc
// @Summary Find member
// @Description Find member
// @Tags alert
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Member found"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /alert/find-member [post]
func (h *AlertHandler) AlertFindMemberHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Member found"})
}
