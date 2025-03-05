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

func instructions(contextVariables map[string]interface{}) string {
	name, ok := contextVariables["name"].(string)
	if !ok {
		name = "User"
	}
	return fmt.Sprintf("You are a helpful agent. Greet the user by name (%s).", name)
}

func printAccountDetails(args map[string]interface{}, contextVariables map[string]interface{}) agentkit.Result {
	userID := contextVariables["user_id"]
	name := contextVariables["name"]
	return agentkit.Result{
		Data: fmt.Sprintf("Account Details: %v %v", name, userID),
	}
}
func main() {
	dotenv.Load()

	client := agentkit.NewSwarm(os.Getenv("OPENAI_API_KEY"), llm.OpenAI)

	agent := &agentkit.Agent{
		Name:             "Agent",
		InstructionsFunc: instructions,
		Functions: []agentkit.AgentFunction{
			{
				Name:        "printAccountDetails",
				Description: "Print the account details of the user.",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
				Function: printAccountDetails,
			},
		},
		Model: "meta-llama/Llama-3.3-70B-Instruct",
	}

	contextVariables := map[string]interface{}{
		"name":    "James",
		"user_id": 123,
	}

	ctx := context.Background()

	// First interaction
	response, err := client.Run(ctx, agent, []llm.Message{
		{Role: "user", Content: "Hi!"},
	}, contextVariables, "", false, false, 5, true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(response.Messages[len(response.Messages)-1].Content)

	// Second interaction
	response, err = client.Run(ctx, agent, []llm.Message{
		{Role: "user", Content: "Print my account details!"},
	}, contextVariables, "", false, false, 5, true)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(response.Messages[len(response.Messages)-1].Content)
}
