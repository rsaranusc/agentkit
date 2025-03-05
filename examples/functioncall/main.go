package main

import (
	"context"
	"fmt"
	"os"

	dotenv "github.com/joho/godotenv"
	agentkit "github.com/rsaranusc/agentkit"
	"github.com/rsaranusc/agentkit/llm"
)

func getWeather(args map[string]interface{}, contextVariables map[string]interface{}) agentkit.Result {
	location := args["location"].(string)
	time := "now"
	if t, ok := args["time"].(string); ok {
		time = t
	}
	return agentkit.Result{
		Success: true,
		Data:    fmt.Sprintf(`{"location": "%s", "temperature": "65", "time": "%s"}`, location, time),
	}
}

func main() {
	dotenv.Load()

	client := agentkit.NewSwarm(os.Getenv("OPENAI_API_KEY"), llm.OpenAI)

	agent := &agentkit.Agent{
		Name:         "Agent",
		Instructions: "You are a helpful agent.",
		Functions: []agentkit.AgentFunction{
			{
				Name:        "getWeather",
				Description: "Get the current weather in a given location.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "The city to get the weather for",
						},
					},
					"required": []string{"location"},
				},
				Function: getWeather,
			},
		},
		Model: "meta-llama/Llama-3.3-70B-Instruct",
	}

	messages := []llm.Message{
		{Role: llm.RoleUser, Content: "What's the weather in NYC?"},
	}

	ctx := context.Background()
	response, err := client.Run(ctx, agent, messages, nil, "", false, false, 5, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Messages[len(response.Messages)-1].Content)
}
