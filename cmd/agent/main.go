package main

import (
	"fmt"
	"log"
	"os"

	"gsk.com/code-orange/agents/internal/agent"
	"gsk.com/code-orange/agents/internal/openai"
	"gsk.com/code-orange/agents/internal/promptpipeline"
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

	pipeline := promptpipeline.NewPipeline([]promptpipeline.PromptTemplate{
		{Name: "Initial", Template: "Analyze this: {{.Input}}"},
		{Name: "Elaborate", Template: "Provide more details on: {{.Input}}"},
		{Name: "Summarize", Template: "Summarize the key points: {{.Input}}"},
	})

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
		pipeline,
	)

	result, err := agent.Execute("Your initial input here")
	if err != nil {
		fmt.Printf("Error processing message: %v\n", err)
		return
	}

	fmt.Println("AI Response:")
	fmt.Println(result)
}
