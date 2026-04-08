package domain

import "time"

// RefreshSession stores a hashed refresh token for a user device.
type RefreshSession struct {
	ID         ID
	UserID     ID
	TokenHash  string
	UserAgent  *string
	IPAddress  *string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	RevokedAt  *time.Time
}
