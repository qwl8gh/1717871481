package controllers

import (
	"log"
	"net/http"
	"time"
	"web-messaging-service/pkg/api/models"
	"web-messaging-service/pkg/services"
	"web-messaging-service/pkg/ws"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	messageService *services.MessageService
}

func NewMessageController(ms *services.MessageService) *MessageController {
	return &MessageController{
		messageService: ms,
	}
}

// @Summary Send message
// @Description Send a new message
// @Tags message
// @Accept json
// @Produce json
// @Param message body models.Message true "Message"
// @Success 200 {object} controllers.MessageResponse
// @Failure 400 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /api/message [post]
func (mc *MessageController) SendMessage(c *gin.Context) {
	var message models.Message
	if err := c.BindJSON(&message); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Received message:", message)
	if err := mc.messageService.ProcessMessage(&message); err != nil {
		log.Println("Error processing message:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not insert message"})
		return
	}

	ws.NotifyClients(message)

	c.JSON(http.StatusOK, message)
}

// @Summary Get messages
// @Description Get messages within a date range
// @Tags message
// @Accept json
// @Produce json
// @Param start_time query string true "Start date"
// @Param end_time query string true "End date"
// @Success 200 {array} controllers.MessageResponse
// @Failure 400 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /api/messages [get]
func (mc *MessageController) GetMessagesByTimeRange(c *gin.Context) {
	startTime, err := time.Parse(time.RFC3339, c.Query("start_time"))
	if err != nil {
		log.Println("Invalid start time:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, c.Query("end_time"))
	if err != nil {
		log.Println("Invalid end time:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time"})
		return
	}

	messages, err := mc.messageService.GetMessagesByTimeRange(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve messages"})
		return
	}

	log.Printf("Retrieved %d messages between %s and %s", len(messages), startTime, endTime)
	c.JSON(http.StatusOK, messages)
}
