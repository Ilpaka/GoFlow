package domain

import "time"

// Chat is a direct or group conversation.
type Chat struct {
	ID            ID
	Type          ChatType
	Title         *string
	AvatarURL     *string
	CreatedBy     *ID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastMessageID *ID
	LastMessageAt *time.Time
	IsDeleted     bool
	DirectKey     *string
}

// ChatMember links a user to a chat with per-chat preferences.
type ChatMember struct {
	ChatID              ID
	UserID              ID
	Role                ChatMemberRole
	JoinedAt            time.Time
	LastReadMessageID   *ID
	LastReadAt          *time.Time
	IsMuted             bool
	IsArchived          bool
	IsPinned            bool
}
