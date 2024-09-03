package embedding

import (
	"fmt"
	"math"
	"sort"

	"gsk.com/code-orange/agents/internal/openai"
)

type EmbeddingMatcher struct {
	client     *openai.Client
	embeddings map[string][]float32
}

func NewEmbeddingMatcher(client *openai.Client) *EmbeddingMatcher {
	return &EmbeddingMatcher{
		client:     client,
		embeddings: make(map[string][]float32),
	}
}

func (m *EmbeddingMatcher) FindBestMatches(query string, candidates []string, n int) ([]string, error) {
	queryEmbedding, err := m.getEmbedding(query)
	if err != nil {
		return nil, err
	}

	similarities := make([]struct {
		candidate  string
		similarity float32
	}, len(candidates))

	for i, candidate := range candidates {
		candidateEmbedding, err := m.getEmbedding(candidate)
		if err != nil {
			return nil, err
		}
		similarities[i] = struct {
			candidate  string
			similarity float32
		}{candidate, cosineSimilarity(queryEmbedding, candidateEmbedding)}
	}

	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].similarity > similarities[j].similarity
	})

	n = min(n, len(similarities))
	bestMatches := make([]string, n)
	for i := 0; i < n; i++ {
		bestMatches[i] = similarities[i].candidate
	}

	return bestMatches, nil
}

func (m *EmbeddingMatcher) getEmbedding(text string) ([]float32, error) {
	if embedding, ok := m.embeddings[text]; ok {
		return embedding, nil
	}

	embedding, err := m.client.CreateEmbedding(text)
	if err != nil {
		return nil, fmt.Errorf("error creating embedding: %w", err)
	}

	m.embeddings[text] = embedding
	return embedding, nil
}

func cosineSimilarity(a, b []float32) float32 {
	var dotProduct, magnitudeA, magnitudeB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}
	return dotProduct / (sqrt(magnitudeA) * sqrt(magnitudeB))
}

func sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *EmbeddingMatcher) SelectBestPrompt(message string, promptPool []string) (string, error) {
	bestMatches, err := m.FindBestMatches(message, promptPool, 1)
	if err != nil {
		return "", fmt.Errorf("error finding best match: %w", err)
	}
	if len(bestMatches) == 0 {
		return "", fmt.Errorf("no matching prompt found")
	}
	return bestMatches[0], nil
}
