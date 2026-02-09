package services

import (
	"context"

	"toolmesh/internal/domain"
	"toolmesh/internal/orchestrator"
)

type NoopRAG struct{}

func (n NoopRAG) Search(ctx context.Context, query string, topK int) ([]string, error) {
	return nil, domain.ErrNoContextFound
}

type NoopTools struct{}

func (n NoopTools) Call(ctx context.Context, tool string, payload any) (any, error) {
	return nil, domain.ErrToolFailed
}

type NoopLLM struct{}

func (n NoopLLM) Generate(ctx context.Context, req orchestrator.ChatRequest) (orchestrator.ChatResponse, error) {
	return orchestrator.ChatResponse{Reply: "LLM gateway not configured"}, nil
}
