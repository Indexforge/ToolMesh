package orchestrator

import (
	"context"
	"errors"
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

	var contextDocs []string
	if o.rag != nil {
		docs, err := o.rag.Search(ctx, req.Message, 3)
		if err != nil {
			if !errors.Is(err, domain.ErrNoContextFound) {
				logger.ErrorContext(ctx, "rag search failed", "error", err)
				return ChatResponse{}, err
			}
			logger.InfoContext(ctx, "rag returned no context")
		} else {
			contextDocs = docs
			logger.InfoContext(ctx, "rag context selected", "doc_count", len(docs))
		}
	}

	resp, err := o.llm.Generate(ctx, LLMRequest{
		Message: req.Message,
		Context: contextDocs,
	})
	if err != nil {
		logger.ErrorContext(ctx, "llm generation failed", "error", err)
		return ChatResponse{}, err
	}

	return resp, nil
}
