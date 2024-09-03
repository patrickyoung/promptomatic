package agent

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"gsk.com/code-orange/agents/internal/openai"
	"gsk.com/code-orange/agents/internal/promptpipeline"
	"gsk.com/code-orange/agents/internal/promptselector"
)

type Agent struct {
	id             string
	name           string
	description    string
	promptPool     []string
	client         *openai.Client
	promptSelector promptselector.PromptSelector
	logger         *log.Logger
	Pipeline       *promptpipeline.Pipeline
}

func New(id, name, description string, promptPool []string, client *openai.Client, selector promptselector.PromptSelector, logger *log.Logger, pipeline *promptpipeline.Pipeline) *Agent {
	return &Agent{
		id:             id,
		name:           name,
		description:    description,
		promptPool:     promptPool,
		client:         client,
		promptSelector: selector,
		logger:         logger,
		Pipeline:       pipeline,
	}
}

func (a *Agent) ProcessMessage(message string, promptValues map[string]string) (string, error) {
	bestPrompt, err := a.promptSelector.SelectBestPrompt(a.promptPool, message)
	if err != nil {
		return "", fmt.Errorf("error selecting best prompt: %w", err)
	}

	a.logger.Printf("Selected prompt: %s", bestPrompt)

	renderedPrompt, err := a.renderPrompt(bestPrompt, promptValues)
	if err != nil {
		return "", fmt.Errorf("error rendering prompt: %w", err)
	}

	response, err := a.client.CreateChatCompletion(renderedPrompt, message)
	if err != nil {
		return "", fmt.Errorf("error creating chat completion: %w", err)
	}

	return response, nil
}

func (a *Agent) renderPrompt(prompt string, values map[string]string) (string, error) {
	tmpl, err := template.New("prompt").Parse(prompt)
	if err != nil {
		return "", fmt.Errorf("error parsing prompt template: %w", err)
	}

	var renderedPrompt bytes.Buffer
	err = tmpl.Execute(&renderedPrompt, values)
	if err != nil {
		return "", fmt.Errorf("error rendering prompt template: %w", err)
	}

	return renderedPrompt.String(), nil
}

func (a *Agent) Execute(input string) (string, error) {
	return a.Pipeline.Execute(a.client, input, map[string]interface{}{})
}
