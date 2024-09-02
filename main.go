package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/sashabaranov/go-openai"
)

type Tool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  []OpenAIParameter `json:"parameters"`
}

type OpenAIParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Agent struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Prompt        string `json:"prompt"`
	Knowledgebase string `json:"knowledgebase"`
	Tools         []Tool `json:"tools"`
}

// NewAgent creates and returns a new Agent instance
func NewAgent(id, name, description, prompt, knowledgebase string, tools []Tool) *Agent {
	return &Agent{
		ID:            id,
		Name:          name,
		Description:   description,
		Prompt:        prompt,
		Knowledgebase: knowledgebase,
		Tools:         tools,
	}
}

// SubmitConversation sends a conversation to the OpenAI API and returns the response
func (a *Agent) SubmitConversation(messages []openai.ChatCompletionMessage, promptValues map[string]string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	// Render the prompt template if promptValues are provided
	if len(promptValues) > 0 {
		tmpl, err := template.New("prompt").Parse(a.Prompt)
		if err != nil {
			return "", fmt.Errorf("error parsing prompt template: %v", err)
		}

		var renderedPrompt bytes.Buffer
		err = tmpl.Execute(&renderedPrompt, promptValues)
		if err != nil {
			return "", fmt.Errorf("error rendering prompt template: %v", err)
		}

		// Add the rendered prompt as a system message at the beginning of the conversation
		messages = append([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: renderedPrompt.String(),
			},
		}, messages...)
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error creating chat completion: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func main() {
	// Create a sample OpenAIParameter
	param := OpenAIParameter{
		Name: "query",
		Type: "string",
	}

	// Create a sample Tool
	tool := Tool{
		Name:        "search",
		Description: "Search for information",
		Parameters:  []OpenAIParameter{param},
	}

	// Create a new Agent using NewAgent function
	agent := NewAgent(
		"agent001",
		"InfoSeeker",
		"An agent that searches for information",
		"You are an AI assistant named {{.name}}. Your task is to {{.task}}.",
		"general",
		[]Tool{tool},
	)

	// Test the SubmitConversation method
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Hello, I'm looking for information about ranbows.",
		},
	}

	promptValues := map[string]string{
		"name": "InfoSeeker",
		"task": "provide information about various topics",
	}

	response, err := agent.SubmitConversation(messages, promptValues)
	if err != nil {
		fmt.Printf("Error submitting conversation: %v\n", err)
		return
	}

	fmt.Println("AI Response:")
	fmt.Println(response)
}
