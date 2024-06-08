// pkg/services/message_service.go

package services

import (
	"log"
	"time"
	"web-messaging-service/pkg/api/models"
	"web-messaging-service/pkg/db"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (ms *MessageService) ProcessMessage(msg *models.Message) error {
	log.Println("Processing message:", msg)
	// Add any business logic here, such as validating message content or formatting timestamps
	msg.Timestamp = time.Now()

	// Persist the message to the database
	if err := db.SendMessage(msg); err != nil {
		log.Printf("Error processing message: %v\n", err)
		return err
	}

	log.Println("Message processed successfully")
	return nil
}

func (ms *MessageService) GetMessagesByTimeRange(startTime, endTime time.Time) ([]models.Message, error) {
	log.Printf("Retrieving messages between %v and %v from database\n", startTime, endTime)
	return db.GetMessagesByTimeRange(startTime, endTime)
}
