package domain

import (
	"time"
)

type Notification struct {
	NotificationID string    `json:"notification_id"`
	UserID         string    `json:"user_id"`
	SenderID       string    `json:"sender_id"`
	Read           bool      `json:"read"`
	CreatedAt      time.Time `json:"created_at"`
	Type           string    `json:"type"`
}
