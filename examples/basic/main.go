package main

import (
	"context"
	"fmt"
	"os"

	dotenv "github.com/joho/godotenv"
	agentkit "github.com/rsaranusc/agentkit"
	"github.com/rsaranusc/agentkit/llm"
)

func main() {
	dotenv.Load()

	client := agentkit.NewSwarm(os.Getenv("OPENAI_API_KEY"), llm.OpenAI)

	agent := &agentkit.Agent{
		Name:         "Agent",
		Instructions: "You are a helpful agent.",
		Model:        "meta-llama/Llama-3.3-70B-Instruct",
	}

	messages := []llm.Message{
		{Role: llm.RoleUser, Content: "Hi!"},
	}

	ctx := context.Background()
	response, err := client.Run(ctx, agent, messages, nil, "", false, false, 5, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Messages[len(response.Messages)-1].Content)
}
