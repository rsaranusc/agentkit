package main

import (
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
		Data:    fmt.Sprintf(`The temperature in %s is 65 degrees at %s.`, location, time),
	}
}

func sendEmail(args map[string]interface{}, contextVariables map[string]interface{}) agentkit.Result {
	recipient := args["recipient"].(string)
	subject := args["subject"].(string)
	body := args["body"].(string)
	fmt.Println("Sending email...")
	fmt.Printf("To: %s\nSubject: %s\nBody: %s\n", recipient, subject, body)
	return agentkit.Result{
		Success: true,
		Data:    "Sent!",
	}
}

func main() {
	dotenv.Load()

	// Create a new Swarm client using DeepSeek
	client := agentkit.NewSwarm(os.Getenv("DEEPSEEK_API_KEY"), llm.DeepSeek)

	weatherAgent := &agentkit.Agent{
		Name:         "WeatherAgent",
		Instructions: "You are a helpful weather assistant powered by DeepSeek. Always respond in a natural, conversational way. When providing weather information, format it in a friendly manner rather than just returning raw data. For example, instead of showing JSON, say something like 'The temperature in [city] is [temp] degrees.'",
		Functions: []agentkit.AgentFunction{
			{
				Name:        "getWeather",
				Description: "Get the current weather in a given location. Location MUST be a city.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "The city to get the weather for",
						},
						"time": map[string]interface{}{
							"type":        "string",
							"description": "The time to get the weather for",
						},
					},
					"required": []interface{}{"location"},
				},
				Function: getWeather,
			},
			{
				Name:        "sendEmail",
				Description: "Send an email to a recipient.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"recipient": map[string]interface{}{
							"type":        "string",
							"description": "The recipient of the email",
						},
						"subject": map[string]interface{}{
							"type":        "string",
							"description": "The subject of the email",
						},
						"body": map[string]interface{}{
							"type":        "string",
							"description": "The body of the email",
						},
					},
					"required": []interface{}{"recipient", "subject", "body"},
				},
				Function: sendEmail,
			},
		},
		Model: "deepseek-chat", // Using DeepSeek's chat model
	}

	agentkit.RunDemoLoop(client, weatherAgent)
}
