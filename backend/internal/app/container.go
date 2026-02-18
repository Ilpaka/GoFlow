package app

import (
	"log/slog"

	"goflow/backend/internal/config"
)

type Container struct {
	Config *config.Config
	Logger *slog.Logger
}

func NewContainer(cfg *config.Config, log *slog.Logger) *Container {
	return &Container{
		Config: cfg,
		Logger: log,
	}
}
