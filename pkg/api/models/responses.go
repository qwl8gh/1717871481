package models

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	ID             int       `json:"id"`
	SequenceNumber int       `json:"sequence_number"`
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
}
