package promptselector

import (
	"gsk.com/code-orange/agents/internal/embedding"
	"gsk.com/code-orange/agents/internal/openai"
)

type PromptSelector struct {
	matcher *embedding.EmbeddingMatcher
	client  *openai.Client
}

func NewPromptSelector(client *openai.Client) *PromptSelector {
	return &PromptSelector{
		matcher: embedding.NewEmbeddingMatcher(client),
		client:  client,
	}
}

func (s *PromptSelector) SelectBestPrompt(prompts []string, message string) (string, error) {
	bestMatches, err := s.matcher.FindBestMatches(message, prompts, 1)
	if err != nil {
		return "", err
	}
	return bestMatches[0], nil
}

// Rename the interface to avoid conflict
type PromptSelectorInterface interface {
	SelectBestPrompt(message string, promptPool []string) (string, error)
}
