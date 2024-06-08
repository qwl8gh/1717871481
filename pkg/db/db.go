package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"web-messaging-service/pkg/api/models"
	"web-messaging-service/pkg/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
		return err
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("could not ping database: %v", err)
		return err
	}

	log.Println("Connected to database")
	return nil
}

func SendMessage(msg *models.Message) error {
	_, err := DB.Exec("INSERT INTO messages (sequence_number, content, timestamp) VALUES ($1, $2, $3)",
		msg.SequenceNumber, msg.Content, msg.Timestamp)
	if err != nil {
		log.Printf("Error inserting message into database: %v\n", err)
		return err
	}
	log.Printf("Message inserted successfully: %v\n", msg)
	return nil
}

func GetMessagesByTimeRange(startTime, endTime time.Time) ([]models.Message, error) {
	rows, err := DB.Query("SELECT sequence_number, content, timestamp FROM messages WHERE timestamp BETWEEN $1 AND $2",
		startTime, endTime)
	if err != nil {
		log.Printf("Error retrieving messages from database: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.SequenceNumber, &msg.Content, &msg.Timestamp); err != nil {
			log.Printf("Error scanning message row: %v\n", err)
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating through message rows: %v\n", err)
		return nil, err
	}

	log.Printf("Retrieved %d messages from database\n", len(messages))
	return messages, nil
}
