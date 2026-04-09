package repository

import (
	"context"

	"goflow/backend/internal/domain"
)

// Page is offset-based pagination for list/search queries.
type Page struct {
	Limit  int
	Offset int
}

// UpdateUserProfileParams updates profile fields; nil pointer means "leave unchanged".
type UpdateUserProfileParams struct {
	UserID    domain.ID
	Nickname  *string
	FirstName *string
	LastName  *string
	AvatarURL *string
}

// UserRepository persists users.
type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id domain.ID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByNickname(ctx context.Context, nickname string) (*domain.User, error)
	Search(ctx context.Context, query string, page Page) ([]domain.User, error)
	UpdateProfile(ctx context.Context, p UpdateUserProfileParams) error
}
