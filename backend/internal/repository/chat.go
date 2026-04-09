package repository

import (
	"context"
	"time"

	"goflow/backend/internal/domain"
)

// AddChatMember is a row for bulk insert into chat_members.
type AddChatMember struct {
	UserID domain.ID
	Role   domain.ChatMemberRole
}

// ChatRepository persists chats and membership.
type ChatRepository interface {
	CreateChat(ctx context.Context, c *domain.Chat) error
	CreateMember(ctx context.Context, m *domain.ChatMember) error
	GetByID(ctx context.Context, id domain.ID) (*domain.Chat, error)
	GetDirectByKey(ctx context.Context, directKey string) (*domain.Chat, error)
	GetUserChats(ctx context.Context, userID domain.ID, page Page) ([]domain.Chat, error)
	GetChatMembers(ctx context.Context, chatID domain.ID) ([]domain.ChatMember, error)
	AddMembers(ctx context.Context, chatID domain.ID, members []AddChatMember) error
	UpdateLastMessage(ctx context.Context, chatID, messageID domain.ID, at time.Time) error
	UpdateMemberRead(ctx context.Context, chatID, userID, readUpToMessageID domain.ID) error
}
