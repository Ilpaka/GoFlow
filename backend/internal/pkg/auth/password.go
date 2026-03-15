package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// DefaultBcryptCost is used by HashPassword when cost is not tuned elsewhere.
const DefaultBcryptCost = bcrypt.DefaultCost

// HashPassword returns a bcrypt hash suitable for storage in users.password_hash.
//
// Example:
//
//	hash, err := auth.HashPassword(plain)
//	if err != nil { ... }
//	user.PasswordHash = hash
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("auth: empty password")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(password), DefaultBcryptCost)
	if err != nil {
		return "", fmt.Errorf("auth: hash password: %w", err)
	}
	return string(b), nil
}

// ComparePassword returns nil if the password matches the bcrypt hash.
//
// Example:
//
//	if err := auth.ComparePassword(user.PasswordHash, plain); err != nil {
//	    return apperr.Unauthorized("invalid credentials")
//	}
func ComparePassword(hash, password string) error {
	if hash == "" {
		return fmt.Errorf("auth: empty hash")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("auth: compare password: %w", err)
	}
	return nil
}
