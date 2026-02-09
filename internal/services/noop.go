package services

import (
	"context"
	"log/slog"
	"strings"

	"toolmesh/internal/domain"
	"toolmesh/internal/orchestrator"
)

type MemoryRAG struct {
	logger    *slog.Logger
	documents []string
}

func NewMemoryRAG(logger *slog.Logger, documents []string) *MemoryRAG {
	return &MemoryRAG{
		logger:    logger,
		documents: documents,
	}
}

func (m *MemoryRAG) Search(ctx context.Context, query string, topK int) ([]string, error) {
	log := m.logger.With("module", "rag")
	log.InfoContext(ctx, "rag search", "query", query, "top_k", topK)

	if len(m.documents) == 0 {
		return nil, domain.ErrNoContextFound
	}

	var matches []string
	lower := strings.ToLower(query)
	for _, doc := range m.documents {
		if strings.Contains(strings.ToLower(doc), lower) {
			matches = append(matches, doc)
			if len(matches) >= topK {
				break
			}
		}
	}

	if len(matches) == 0 {
		return nil, domain.ErrNoContextFound
	}

	log.InfoContext(ctx, "rag results", "doc_count", len(matches))
	return matches, nil
}

type NoopTools struct {
	logger *slog.Logger
}

func NewNoopTools(logger *slog.Logger) *NoopTools {
	return &NoopTools{logger: logger}
}

func (n *NoopTools) Call(ctx context.Context, tool string, payload any) (any, error) {
	n.logger.With("module", "mcp").ErrorContext(ctx, "tool call failed", "tool", tool)
	return nil, domain.ErrToolFailed
}

type StubLLM struct {
	logger *slog.Logger
}

func NewStubLLM(logger *slog.Logger) *StubLLM {
	return &StubLLM{logger: logger}
}

func (s *StubLLM) Generate(ctx context.Context, req orchestrator.LLMRequest) (orchestrator.ChatResponse, error) {
	log := s.logger.With("module", "llm")
	log.InfoContext(ctx, "llm generate", "context_count", len(req.Context))

	reply := "LLM gateway not configured"
	if len(req.Context) > 0 {
		reply = reply + " (context available)"
	}

	return orchestrator.ChatResponse{Reply: reply}, nil
}
