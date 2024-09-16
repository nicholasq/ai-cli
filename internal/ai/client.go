package ai

import (
	"context"

	"nicholasq.xyz/ai/internal/config"
)

type AIClient interface {
	Query(ctx context.Context, query string, config config.Config) (*AIResponse, error)
	GetCapabilities() []string
	SetContext(context string) error
}

type AIResponse struct {
	Text string
	// Add other fields as needed
}
