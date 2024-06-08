package models

import "time"

type Message struct {
	ID             int       `json:"id"`
	SequenceNumber int       `json:"sequence_number"`
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
}
