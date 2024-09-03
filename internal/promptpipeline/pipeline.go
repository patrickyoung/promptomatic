package promptpipeline

import (
	"bytes"
	"text/template"

	"gsk.com/code-orange/agents/internal/openai"
)

type PromptTemplate struct {
	Name     string
	Template string
}

type Pipeline struct {
	Templates []PromptTemplate
}

func NewPipeline(templates []PromptTemplate) *Pipeline {
	return &Pipeline{Templates: templates}
}

func (p *Pipeline) Execute(client *openai.Client, initialInput string, variables map[string]interface{}) (string, error) {
	result := initialInput
	for _, promptTemplate := range p.Templates {
		// Parse the template
		tmpl, err := template.New(promptTemplate.Name).Parse(promptTemplate.Template)
		if err != nil {
			return "", err
		}

		// Execute the template with variables
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, variables)
		if err != nil {
			return "", err
		}

		// Use the client to get a response from the API
		response, err := client.CreateChatCompletion(buf.String(), "gpt-4o-mini")
		if err != nil {
			return "", err
		}

		// Update the result and variables for the next iteration
		result = response
		variables["Input"] = result
	}
	return result, nil
}
