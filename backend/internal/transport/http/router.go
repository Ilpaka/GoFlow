package httptransport

import (
	"log/slog"
	"net/http"

	"goflow/backend/internal/config"
)

// Deps are transport-level dependencies supplied by the app layer (typically from app.Container).
type Deps struct {
	Config *config.Config
	Logger *slog.Logger
}

// Register wires HTTP routes. Handlers stay thin; domain logic lives outside handlers.
func Register(mux *http.ServeMux, deps *Deps) {
	mux.Handle("GET /health", healthHandler(deps))
}

func healthHandler(_ *Deps) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}
