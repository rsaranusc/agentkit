package main

import (
	"context"
	"fmt"
	"os"

	dotenv "github.com/joho/godotenv"
	agentkit "github.com/rsaranusc/agentkit"
	"github.com/rsaranusc/agentkit/llm"
	"github.com/spf13/cast"
)

func calculateSum(args map[string]any, contextVariables map[string]any) agentkit.Result {
	num1 := cast.ToInt(args["num1"])
	num2 := cast.ToInt(args["num2"])
	sum := num1 + num2
	return agentkit.Result{
		Success: true,
		Data:    fmt.Sprintf("The sum of %d and %d is %d", num1, num2, sum),
	}
}

func calculateProduct(args map[string]any, contextVariables map[string]any) agentkit.Result {
	num1 := cast.ToInt(args["num1"])
	num2 := cast.ToInt(args["num2"])
	product := num1 * num2
	return agentkit.Result{
		Success: true,
		Data:    fmt.Sprintf("The product of %d and %d is %d", num1, num2, product),
	}
}

func main() {
	dotenv.Load()

	client := agentkit.NewSwarm(os.Getenv("OPENAI_API_KEY"), llm.OpenAI)

	mathAgent := &agentkit.Agent{
		Name:         "MathAgent",
		Instructions: "You are a math assistant. When given two numbers, calculate both their sum and product.",
		Functions: []agentkit.AgentFunction{
			{
				Name:        "calculateSum",
				Description: "Calculate the sum of two numbers",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"num1": map[string]any{
							"type":        "number",
							"description": "First number",
						},
						"num2": map[string]any{
							"type":        "number",
							"description": "Second number",
						},
					},
					"required": []any{"num1", "num2"},
				},
				Function: calculateSum,
			},
			{
				Name:        "calculateProduct",
				Description: "Calculate the product of two numbers",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"num1": map[string]any{
							"type":        "number",
							"description": "First number",
						},
						"num2": map[string]any{
							"type":        "number",
							"description": "Second number",
						},
					},
					"required": []any{"num1", "num2"},
				},
				Function: calculateProduct,
			},
		},
		Model: "Qwen/Qwen2.5-32B-Instruct",
	}

	// Create context
	ctx := context.Background()

	// Example message asking to perform calculations
	messages := []llm.Message{
		{
			Role:    llm.RoleUser,
			Content: "Calculate the sum and product of 5 and 3",
		},
	}

	// Run the agent with tool execution enabled
	response, err := client.Run(ctx, mathAgent, messages, nil, "", false, true, 1, true)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the final response from the agent
	fmt.Println("\nAgent's Final Response:")
	for _, msg := range response.Messages {
		if msg.Role == llm.RoleAssistant {
			fmt.Printf("Assistant: %s\n", msg.Content)
		}
	}

	// Print detailed information about tool calls
	fmt.Println("\nTool Call Results:")
	for _, result := range response.ToolResults {
		fmt.Printf("\nTool: %s\n", result.ToolName)
		fmt.Printf("Arguments: %v\n", result.Args)
		fmt.Printf("Result: %v\n", result.Result.Data)

		// You can also check if the tool call was successful
		if result.Result.Success {
			fmt.Printf("Status: Success\n")
		} else {
			fmt.Printf("Status: Failed\nError: %v\n", result.Result.Error)
		}
	}
}
