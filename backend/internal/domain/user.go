package domain

import "time"

// User is a registered account.
type User struct {
	ID           ID
	Email        string
	PasswordHash string
	Nickname     string
	FirstName    *string
	LastName     *string
	AvatarURL    *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastSeenAt   *time.Time
	IsActive     bool
}
