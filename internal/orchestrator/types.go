package orchestrator

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

type LLMRequest struct {
	Message string   `json:"message"`
	Context []string `json:"context,omitempty"`
}
