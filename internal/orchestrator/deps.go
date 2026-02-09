package orchestrator

import "context"

type RAGService interface {
	Search(ctx context.Context, query string, topK int) ([]string, error)
}

type ToolClient interface {
	Call(ctx context.Context, tool string, payload any) (any, error)
}

type LLMGateway interface {
	Generate(ctx context.Context, req ChatRequest) (ChatResponse, error)
}
