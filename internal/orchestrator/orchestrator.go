package orchestrator

import (
	"context"
	"log/slog"

	"toolmesh/internal/domain"
)

type Orchestrator struct {
	rag    RAGService
	tools  ToolClient
	llm    LLMGateway
	logger *slog.Logger
}

func New(rag RAGService, tools ToolClient, llm LLMGateway, logger *slog.Logger) *Orchestrator {
	return &Orchestrator{
		rag:    rag,
		tools:  tools,
		llm:    llm,
		logger: logger,
	}
}

func (o *Orchestrator) HandleChat(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	if o.llm == nil {
		return ChatResponse{}, domain.ErrModelUnavailable
	}

	logger := o.logger.With("module", "orchestrator")
	logger.InfoContext(ctx, "handle chat request")

	resp, err := o.llm.Generate(ctx, req)
	if err != nil {
		logger.ErrorContext(ctx, "llm generation failed", "error", err)
		return ChatResponse{}, err
	}

	return resp, nil
}
