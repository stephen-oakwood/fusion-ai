package a2a

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"log"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

const systemPrompt = `
		"You are an IT Technician, capable of providing detailed answers to the questions that your customers ask regarding their assets." +
		"Think before you reply in <thinking> tags." +
		"First, determine if there are any knowledge articles related to the question that could help with your reply." +
		"Second, using available tools, fetch the schema for a graphQL API that provides the ability to search for assets and return their details." +
		"Third, using the fetched schema, use the graphQL API to collect data from the customer's assets that are relevant to the points raised in the knowledge articles.'" +
		"Finally, provide a detailed summary that includes the relevant points from the knowledge articles with examples of customer's assets that demonstrate these points. Explicitly name the assets when providing examples." +
		"Always fetch the graphQL API Schema first, and construct queries using this schema. Do not construct queries without using the schema."
`

type assetManagementAgent struct {
	InferenceClient *inferenceClient
}

func NewAgent() (*assetManagementAgent, error) {
	inferenceClient, err := NewInferenceClient()
	if err != nil {
		return nil, err
	}

	return &assetManagementAgent{
		InferenceClient: inferenceClient,
	}, nil
}

func (p *assetManagementAgent) Process(ctx context.Context, taskID string, message protocol.Message, handle taskmanager.TaskHandle) error {

	prompt := extractText(message)

	if prompt == "" {
		fmt.Printf("task failed - prompt must contain text")
		failedMessage := protocol.NewMessage(
			protocol.MessageRoleAgent,
			[]protocol.Part{protocol.NewTextPart("input message must contain text")},
		)
		_ = handle.UpdateStatus(protocol.TaskStateFailed, &failedMessage)
		return fmt.Errorf("input message must contain text")
	}

	toolConfig := toolConfig()
	converseInput := &bedrockruntime.ConverseInput{
		System: []types.SystemContentBlock{
			&types.SystemContentBlockMemberText{
				Value: systemPrompt,
			},
		},
		ModelId: &p.InferenceClient.BedrockModel,
		Messages: []types.Message{
			{
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{
						Value: prompt,
					},
				},
				Role: types.ConversationRoleUser,
			},
		},
		ToolConfig: &toolConfig,
	}

	inference, err := p.InferenceClient.MakeInference(ctx, converseInput)
	if err != nil {
		return fmt.Errorf("inference failed: %w", err)
	}

	responseParts := []protocol.Part{protocol.NewTextPart(inference)}

	responseMessage := protocol.NewMessage(
		protocol.MessageRoleAgent,
		responseParts,
	)

	if err := handle.UpdateStatus(protocol.TaskStateCompleted, &responseMessage); err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	artifact := protocol.Artifact{
		Name:        stringPtr("Inference"),
		Description: stringPtr("Inference response from model"),
		Index:       0,
		Parts:       responseParts,
		LastChunk:   boolPtr(true),
	}

	if err := handle.AddArtifact(artifact); err != nil {
		log.Printf("Error adding artifact for task %s: %v", taskID, err)
	}

	return nil
}

func boolPtr(b bool) *bool {
	return &b
}

func extractText(message protocol.Message) string {
	for _, part := range message.Parts {
		if textPart, ok := part.(protocol.TextPart); ok {
			return textPart.Text
		}
	}
	return ""
}
