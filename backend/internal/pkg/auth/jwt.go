package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessClaims is embedded in signed access JWTs.
// UserID is stored both as a dedicated claim and as RegisteredClaims.Subject.
//
// Example (login handler):
//
//	token, err := auth.GenerateAccessToken([]byte(cfg.JWT.Secret), string(user.ID), time.Duration(cfg.JWT.AccessTTLSeconds)*time.Second)
//	if err != nil { ... }
//
// Example (middleware):
//
//	claims, err := auth.ParseAccessToken([]byte(cfg.JWT.Secret), raw)
//	if err != nil { return apperr.Unauthorized("invalid token") }
//	uid := claims.UserID
type AccessClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken issues an HS256 access token with user_id and standard time claims.
func GenerateAccessToken(secret []byte, userID string, ttl time.Duration) (string, error) {
	if len(secret) == 0 {
		return "", fmt.Errorf("auth: empty jwt secret")
	}
	if userID == "" {
		return "", fmt.Errorf("auth: empty user id")
	}
	if ttl <= 0 {
		return "", fmt.Errorf("auth: non-positive ttl")
	}
	now := time.Now().UTC()
	claims := AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

// ParseAccessToken validates signature and expiry and returns typed claims.
func ParseAccessToken(secret []byte, tokenString string) (*AccessClaims, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("auth: empty jwt secret")
	}
	if tokenString == "" {
		return nil, fmt.Errorf("auth: empty token")
	}
	claims := &AccessClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("auth: unexpected signing method %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("auth: parse access token: %w", err)
	}
	if claims.UserID == "" && claims.Subject != "" {
		claims.UserID = claims.Subject
	}
	if claims.UserID == "" {
		return nil, fmt.Errorf("auth: missing user_id claim")
	}
	return claims, nil
}
