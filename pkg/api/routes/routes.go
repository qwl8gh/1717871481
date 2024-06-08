package routes

import (
	"web-messaging-service/pkg/api/controllers"
	"web-messaging-service/pkg/services"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, messageService *services.MessageService) {
	messageController := controllers.NewMessageController(messageService)
	api := router.Group("/api")
	{
		api.POST("/message", messageController.SendMessage)
		api.GET("/messages", messageController.GetMessagesByTimeRange)
	}
}
