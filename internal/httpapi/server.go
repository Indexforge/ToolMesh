package httpapi

import (
	"log/slog"
	"net/http"

	"toolmesh/internal/orchestrator"
)

func NewMux(orch *orchestrator.Orchestrator, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/chat", ChatHandler(orch, logger))
	return mux
}
