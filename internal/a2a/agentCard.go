package a2a

import (
	"trpc.group/trpc-go/trpc-a2a-go/server"
)

func GetAgentCard() server.AgentCard {
	agentCard := server.AgentCard{
		Name:        "Asset Management",
		Description: "An agent that can answer questions related to a customer's managed assets.",
		URL:         "http://localhost:8080",
		Version:     "1.0.0",
		Provider: &server.AgentProvider{
			Organization: "n-able",
		},
		Capabilities: server.AgentCapabilities{
			Streaming:              boolPtr(true),
			PushNotifications:      boolPtr(false),
			StateTransitionHistory: boolPtr(true),
		},
		DefaultInputModes:  []string{"text"},
		DefaultOutputModes: []string{"text"},
		Skills: []server.AgentSkill{
			{
				ID:          "query_schema",
				Name:        "Query Schema",
				Description: stringPtr("Returns a partial GraphQL Schema that can be used to construct queries and mutations for an API that supports searching for managed assets and returning details regarding their operating systems, hardware and more."),
				InputModes:  []string{"text"},
				OutputModes: []string{"text"},
			},
			{
				ID:          "execute_query",
				Name:        "Execute Query",
				Description: stringPtr("Executes a GraphQL Query using the n-able public API. Provides support for sophisticated searching of assets."),
				InputModes:  []string{"text"},
				OutputModes: []string{"text"},
			},
			{
				ID:          "knowledge_query",
				Name:        "Knowledge Query",
				Description: stringPtr("Finds the most relevant knowledge article for a user's questions about managed assets"),
				InputModes:  []string{"text"},
				OutputModes: []string{"text"},
			},
		},
	}

	return agentCard
}

func stringPtr(s string) *string {
	return &s
}
