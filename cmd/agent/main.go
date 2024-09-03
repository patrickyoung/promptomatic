package main

import (
	"fmt"
	"log"
	"os"

	"gsk.com/code-orange/agents/internal/agent"
	"gsk.com/code-orange/agents/internal/openai"
	"gsk.com/code-orange/agents/internal/promptselector"
)

func main() {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	logFile, err := os.Create("llm_interactions.log")
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	agent := agent.New(
		"agent001",
		"InfoSeeker",
		"An agent that searches for information",
		[]string{
			"You are an AI assistant named {{.name}}. Your task is to {{.task}}.",
			"As {{.name}}, your primary function is to {{.task}}. Provide verbose, fully detailed, and accurate information.",
			"You are an AI assistant named {{.name}}. Your task is to {{.task}}. Omly provide a one sentence answer.",
		},
		client,
		*promptselector.NewPromptSelector(client),
		logger,
	)

	response, err := agent.ProcessMessage("provide a detailed essay", map[string]string{
		"name": "an expert in the field of genitic and cell biology",
		"task": "everything you know about munchkin cats",
	})
	if err != nil {
		fmt.Printf("Error processing message: %v\n", err)
		return
	}

	fmt.Println("AI Response:")
	fmt.Println(response)
}
