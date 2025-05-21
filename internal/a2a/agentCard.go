package a2a

import (
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/server"
)

func GetAgentCard() server.AgentCard {
	agentCard := server.AgentCard{
		Name:        "Asset Management",
		Description: stringPtr("An agent that can answer questions related to a customer's managed assets."),
		URL:         "http://localhost:8080",
		Version:     "1.0.0",
		Provider: &server.AgentProvider{
			Organization: "n-able",
		},
		Capabilities: server.AgentCapabilities{
			Streaming:              false,
			PushNotifications:      false,
			StateTransitionHistory: true,
		},
		DefaultInputModes:  []string{string(protocol.PartTypeText)},
		DefaultOutputModes: []string{string(protocol.PartTypeText)},
		Skills: []server.AgentSkill{
			{
				ID:          "query_schema",
				Name:        "Query Schema",
				Description: stringPtr("Returns a partial GraphQL Schema that can be used to construct queries and mutations for an API that supports searching for managed assets and returning details regarding their operating systems, hardware and more."),
				InputModes:  []string{string(protocol.PartTypeText)},
				OutputModes: []string{string(protocol.PartTypeText)},
			},
			{
				ID:          "execute_query",
				Name:        "Execute Query",
				Description: stringPtr("Executes a GraphQL Query using the n-able public API. Provides support for sophisticated searching of assets."),
				InputModes:  []string{string(protocol.PartTypeText)},
				OutputModes: []string{string(protocol.PartTypeText)},
			},
			{
				ID:          "knowledge_query",
				Name:        "Knowledge Query",
				Description: stringPtr("Finds the most relevant knowledge article for a user's questions about managed assets"),
				InputModes:  []string{string(protocol.PartTypeText)},
				OutputModes: []string{string(protocol.PartTypeText)},
			},
		},
	}

	return agentCard
}

func stringPtr(s string) *string {
	return &s
}
