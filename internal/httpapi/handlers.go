package httpapi

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"toolmesh/internal/domain"
	"toolmesh/internal/orchestrator"
)

func ChatHandler(orch *orchestrator.Orchestrator, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = newRequestID()
		}

		ctx := withRequestID(r.Context(), requestID)
		log := logger.With("module", "http", "request_id", requestID)
		log.InfoContext(ctx, "request started", "method", r.Method, "path", r.URL.Path)
		defer func() {
			log.InfoContext(ctx, "request completed", "duration_ms", time.Since(start).Milliseconds())
		}()

		var req orchestrator.ChatRequest
		if err := decodeJSON(r.Body, &req); err != nil {
			log.ErrorContext(ctx, "decode json failed", "error", err)
			writeError(w, domain.ErrBadRequest)
			return
		}

		if strings.TrimSpace(req.Message) == "" {
			writeError(w, domain.ErrBadRequest)
			return
		}

		resp, err := orch.HandleChat(ctx, req)
		if err != nil {
			log.ErrorContext(ctx, "handle chat failed", "error", err)
			writeError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, resp)
	}
}
