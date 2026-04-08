package domain

import "time"

// Message is a single item in a chat timeline.
type Message struct {
	ID         ID
	ChatID     ID
	SenderID   ID
	Type       MessageType
	Text       *string
	ReplyToID  *ID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
