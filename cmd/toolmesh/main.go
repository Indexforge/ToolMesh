package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"toolmesh/internal/httpapi"
	"toolmesh/internal/orchestrator"
	"toolmesh/internal/services"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	rag := services.NewMemoryRAG(logger, nil)
	tools := services.NewNoopTools(logger)
	llm := services.NewStubLLM(logger)
	orch := orchestrator.New(rag, tools, llm, logger)

	mux := httpapi.NewMux(orch, logger)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("http server starting", "module", "http", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server failed", "module", "http", "error", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("http server shutting down", "module", "http")
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("http server shutdown failed", "module", "http", "error", err)
	}
}
