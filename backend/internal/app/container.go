package app

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"goflow/backend/internal/config"
	"goflow/backend/internal/repository"
	"goflow/backend/internal/repository/postgres"
)

// Container wires application dependencies.
type Container struct {
	Config *config.Config
	Logger *slog.Logger
	Pool   *pgxpool.Pool

	Users    repository.UserRepository
	Chats    repository.ChatRepository
	Messages repository.MessageRepository
	Sessions repository.SessionRepository
}

// NewContainer returns a container. When pool is non-nil, PostgreSQL repositories are constructed.
func NewContainer(cfg *config.Config, log *slog.Logger, pool *pgxpool.Pool) *Container {
	c := &Container{
		Config: cfg,
		Logger: log,
		Pool:   pool,
	}
	if pool != nil {
		c.Users = postgres.NewUserRepository(pool)
		c.Chats = postgres.NewChatRepository(pool)
		c.Messages = postgres.NewMessageRepository(pool)
		c.Sessions = postgres.NewSessionRepository(pool)
	}
	return c
}
