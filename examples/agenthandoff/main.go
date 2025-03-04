package main

import (
	"context"
	"fmt"
	"log"
	"os"

	dotenv "github.com/joho/godotenv"
	agentkit "github.com/rsaranusc/agentkit"
	"github.com/rsaranusc/agentkit/llm"
)

func transferToSpanishAgent(args map[string]interface{}, contextVariables map[string]interface{}) agentkit.Result {
	spanishAgent := &agentkit.Agent{
		Name:         "SpanishAgent",
		Instructions: "You only speak Spanish.",
		Model:        "gpt-4",
	}
	return agentkit.Result{
		Agent: spanishAgent,
		Data:  "Transferring to Spanish Agent.",
	}
}
func main() {
	dotenv.Load()

	client := agentkit.NewSwarm(os.Getenv("OPENAI_API_KEY"), llm.OpenAI)

	englishAgent := &agentkit.Agent{
		Name:         "EnglishAgent",
		Instructions: "You only speak English.",
		Functions: []agentkit.AgentFunction{
			{
				Name:        "transferToSpanishAgent",
				Description: "Transfer Spanish-speaking users immediately.",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
				Function: transferToSpanishAgent,
			},
		},
		Model: "gpt-4",
	}

	messages := []llm.Message{
		{Role: "user", Content: "Hola. ¿Cómo estás?"},
	}

	ctx := context.Background()
	response, err := client.Run(ctx, englishAgent, messages, nil, "", false, false, 5, true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("%s: %s\n", response.Agent.Name, response.Messages[len(response.Messages)-1].Content)
}
