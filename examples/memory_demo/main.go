package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rsaranusc/agentkit"
	"github.com/rsaranusc/agentkit/llm"
	"github.com/spf13/cast"
)

// createMemoryAgent creates an agent with memory capabilities and custom functions
func createMemoryAgent() *agentkit.Agent {
	agent := agentkit.NewAgent("MemoryAgent", "meta-llama/Llama-3.3-70B-Instruct", llm.OpenAI)
	agent.Instructions = `You are a helpful assistant with memory capabilities. 
	You can remember our conversations and use that information in future responses.
	When asked about past interactions, search your memories and provide relevant information.`

	agent.Functions = []agentkit.AgentFunction{
		{
			Name:        "store_fact",
			Description: "Store an important fact in memory",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"content": map[string]any{
						"type":        "string",
						"description": "The fact to remember",
					},
					"importance": map[string]any{
						"type":        "number",
						"description": "Importance score (0-1)",
					},
				},
				"required": []any{"content", "importance"},
			},
			Function: func(args map[string]any, contextVars map[string]any) agentkit.Result {
				content := args["content"].(string)
				importance := cast.ToFloat64(args["importance"])

				memory := agentkit.Memory{
					Content:    content,
					Type:       "fact",
					Context:    contextVars,
					Timestamp:  time.Now(),
					Importance: importance,
				}

				// Add the memory to the agent's memory store
				agent.Memory.AddMemory(memory)

				return agentkit.Result{
					Data: fmt.Sprintf("Stored fact: %s (importance: %.2f)", content, importance),
				}
			},
		},
		{
			Name:        "recall_memories",
			Description: "Recall recent memories or search for specific types of memories",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"memory_type": map[string]any{
						"type":        "string",
						"description": "Type of memories to recall (conversation, fact, tool_result)",
					},
					"count": map[string]any{
						"type":        "number",
						"description": "Number of recent memories to recall",
					},
				},
				"required": []any{"memory_type", "count"},
			},
			Function: func(args map[string]any, contextVars map[string]any) agentkit.Result {
				memoryType := args["memory_type"].(string)
				count := cast.ToInt(args["count"])

				var memories []agentkit.Memory
				if memoryType == "recent" {
					memories = agent.Memory.GetRecentMemories(count)
				} else {
					memories = agent.Memory.SearchMemories(memoryType, nil)
					if len(memories) > count {
						memories = memories[len(memories)-count:]
					}
				}
				// Format memories nicely
				var result string
				for i, mem := range memories {
					result += fmt.Sprintf("\n%d. [%s] %s (Importance: %.2f)",
						i+1, mem.Timestamp.Format("15:04:05"), mem.Content, mem.Importance)
				}

				if result == "" {
					result = "No memories found."
				}

				return agentkit.Result{Data: result}
			},
		},
	}

	return agent
}

func main() {
	// Load environment variables
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	godotenv.Load()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create a new swarm and memory-enabled agent
	client := agentkit.NewSwarm(apiKey, llm.OpenAI)
	agent := createMemoryAgent()

	// Example conversation demonstrating memory capabilities
	conversations := []string{
		"Hi! My name is Alice.",
		"Could you store that my favorite color is blue?",
		"What do you remember about me?",
		"I also like cats. Please remember that.",
		"What are all the facts you remember about me?",
	}

	ctx := context.Background()

	fmt.Println("Starting memory demonstration...")
	fmt.Println("=================================")

	// Run through the conversation
	for _, userInput := range conversations {
		fmt.Printf("\n👤 User: %s\n", userInput)

		// Create message for this turn
		messages := []llm.Message{
			{Role: "user", Content: userInput},
		}

		// Get response from agent
		response, err := client.Run(ctx, agent, messages, nil, "", false, false, 5, true)
		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}

		// Print agent's response
		if len(response.Messages) > 0 {
			lastMessage := response.Messages[len(response.Messages)-1]
			fmt.Printf("🤖 Agent: %s\n", lastMessage.Content)
		}

		// Optional: Save memories to file after each interaction
		if data, err := agent.Memory.SerializeMemories(); err == nil {
			os.WriteFile("memories.json", data, 0644)
		}
	}
}
