package main

import (
	"fmt"
	"os"

	dotenv "github.com/joho/godotenv"
	agentkit "github.com/rsaranusc/agentkit"
	"github.com/rsaranusc/agentkit/llm"
)

func processRefund(args map[string]interface{}, contextVariables map[string]interface{}) agentkit.Result {
	itemID := args["item_id"].(string)
	reason := "NOT SPECIFIED"
	if val, ok := args["reason"].(string); ok {
		reason = val
	}
	fmt.Printf("[mock] Refunding item %s because %s...\n", itemID, reason)
	return agentkit.Result{
		Data: fmt.Sprintf("Refunded item %s because %s.", itemID, reason),
	}
}

func applyDiscount(args map[string]interface{}, contextVariables map[string]interface{}) agentkit.Result {
	fmt.Println("[mock] Applying discount...")
	return agentkit.Result{
		Data: "Applied discount of 11%",
	}
}

func transferBackToTriage(args map[string]interface{}, contextVariables map[string]interface{}) agentkit.Result {
	return agentkit.Result{
		Agent: triageAgent,
		Data:  "Transferring back to TriageAgent.",
	}
}

func transferToSales(args map[string]interface{}, contextVariables map[string]interface{}) agentkit.Result {
	return agentkit.Result{
		Agent: salesAgent,
		Data:  "Transferring to SalesAgent.",
	}
}

func transferToRefunds(args map[string]interface{}, contextVariables map[string]interface{}) agentkit.Result {
	return agentkit.Result{
		Agent: refundsAgent,
		Data:  "Transferring to RefundsAgent.",
	}
}

var triageAgent *agentkit.Agent
var salesAgent *agentkit.Agent
var refundsAgent *agentkit.Agent

func initAgents() {
	triageAgent = &agentkit.Agent{
		Name:         "TriageAgent",
		Instructions: "Determine which agent is best suited to handle the user's request, and transfer the conversation to that agent.",
		Model:        "meta-llama/Llama-3.3-70B-Instruct",
	}

	salesAgent = &agentkit.Agent{
		Name:         "SalesAgent",
		Instructions: "Be super enthusiastic about selling bees. If the user's request is unrelated to sales and related to discount or refund, call the 'transferBackToTriage' function to transfer the conversation back to the triage agent.",
		Model:        "meta-llama/Llama-3.3-70B-Instruct",
	}

	refundsAgent = &agentkit.Agent{
		Name:         "RefundsAgent",
		Instructions: "Assist the user with refund inquiries. If the reason is that it was too expensive, offer the user a discount code. If they insist, acknowledge their request and inform them that the refund process will be initiated through the appropriate channels.",
		Functions: []agentkit.AgentFunction{
			{
				Name:        "processRefund",
				Description: "Process a refund request. Confirm with the user that they wish to proceed with the refund without asking for personal details.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the item to refund.",
						},
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "The reason for the refund.",
						},
					},
					"required": []interface{}{"item_id"},
				},
				Function: processRefund,
			},
			{
				Name:        "applyDiscount",
				Description: "Apply a discount to the user's cart.",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
				Function: applyDiscount,
			},
		},
		Model: "meta-llama/Llama-3.3-70B-Instruct",
	}

	// Assign functions to agents
	triageAgent.Functions = []agentkit.AgentFunction{
		{
			Name:        "transferToSales",
			Description: "Transfer the conversation to the SalesAgent.",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
			Function: transferToSales,
		},
		{
			Name:        "transferToRefunds",
			Description: "Transfer the conversation to the RefundsAgent.",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
			Function: transferToRefunds,
		},
	}

	salesAgent.Functions = []agentkit.AgentFunction{
		{
			Name:        "transferBackToTriage",
			Description: "If you are unable to assist the user or if the user's request is outside your expertise, call this function to transfer the conversation back to the triage agent.",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
			Function: transferBackToTriage,
		},
	}

	refundsAgent.Functions = append(refundsAgent.Functions, agentkit.AgentFunction{
		Name:        "transferBackToTriage",
		Description: "If you are unable to assist the user or if the user's request is outside your expertise, call this function to transfer the conversation back to the triage agent.",
		Parameters: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		Function: transferBackToTriage,
	})

	refundsAgent.Functions = append(refundsAgent.Functions, agentkit.AgentFunction{
		Name:        "transferToSales",
		Description: "Transfer the conversation to the SalesAgent.",
		Parameters: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		Function: transferToSales,
	})
}

func main() {
	dotenv.Load()

	client := agentkit.NewSwarm(os.Getenv("OPENAI_API_KEY"), llm.OpenAI)

	initAgents() // Initialize agents and their functions

	agentkit.RunDemoLoop(client, triageAgent)
}
