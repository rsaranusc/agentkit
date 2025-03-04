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

// WeatherRequest represents the parameters for the getWeather function
type WeatherRequest struct {
	Location string `json:"location"`
}

func main() {
	dotenv.Load()
	// Initialize Gemini client with API key from environment
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_API_KEY environment variable is required")
	}

	swarm := agentkit.NewSwarm(apiKey, llm.Gemini)

	// Example 1: Basic chat completion
	fmt.Println("Example 1: Basic Chat Completion")
	agent := &agentkit.Agent{
		Name:         "Agent",
		Instructions: "You are a helpful agent.",
		Model:        "gemini-pro",
	}

	messages := []llm.Message{
		{Role: llm.RoleUser, Content: "Hi!"},
	}

	ctx := context.Background()
	response, err := swarm.Run(ctx, agent, messages, nil, "", false, false, 5, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Messages[len(response.Messages)-1].Content)

/* 	// Example 2: Function calling
	fmt.Println("\nExample 2: Function Calling")

	weatherAgent := &agentkit.Agent{
		Name:         "WeatherAgent",
		Instructions: "You are a weather assistant. When asked about weather, always use the getWeather function.",
		Model:        "gemini-pro",
		Functions: []agentkit.AgentFunction{
			{
				Name:        "getWeather",
				Description: "Get the current weather for a location",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "The city and state/country",
						},
					},
					"required": []interface{}{"location"},
				},
				Function: func(args map[string]interface{}, contextVariables map[string]interface{}) agentkit.Result {
					location := args["location"].(string)
					return agentkit.Result{
						Success: true,
						Data:    fmt.Sprintf(`{"location": "%s", "temperature": "65"}`, location),
					}
				},
			},
		},
	}

	agentkit.RunDemoLoop(swarm, weatherAgent) */

	// Example 3: Streaming completion
	fmt.Println("\nExample 3: Streaming Completion")

	streamAgent := &agentkit.Agent{
		Name:         "StreamAgent",
		Instructions: "You are a helpful agent that responds in a streaming fashion.",
		Model:        "gemini-pro",
	}

	streamMessages := []llm.Message{
		{Role: llm.RoleUser, Content: "Count from 1 to 5 slowly"},
	}

	streamResponse, err := swarm.Run(ctx, streamAgent, streamMessages, nil, "", true, false, 5, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(streamResponse.Messages[len(streamResponse.Messages)-1].Content)

}
