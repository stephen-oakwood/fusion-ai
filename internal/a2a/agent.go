package a2a

import (
	"context"
	"fmt"
	"fusion/internal/nable"
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
	ModelClient *modelClient
}

func NewAgent() (*assetManagementAgent, error) {
	modelClient, err := NewModelClient()
	if err != nil {
		return nil, err
	}

	return &assetManagementAgent{
		ModelClient: modelClient,
	}, nil
}

func (p *assetManagementAgent) Process(ctx context.Context, taskID string, message protocol.Message, handle taskmanager.TaskHandle) error {

	prompt := extractText(message)
	var responseParts []protocol.Part

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
		ModelId: &p.ModelClient.BedrockModel,
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

	converseLoop := true
	for converseLoop {
		converseOutput, err := p.ModelClient.Converse(ctx, converseInput)
		if err != nil {
			return fmt.Errorf("model inference failed: %w", err)
		}

		converseMessage, ok := converseOutput.Output.(*types.ConverseOutputMemberMessage)
		if !ok {
			return fmt.Errorf("error casting LLM response")
		}

		switch converseOutput.StopReason {
		case types.StopReasonEndTurn:
			content, ok := converseMessage.Value.Content[0].(*types.ContentBlockMemberText)
			if !ok {
				return fmt.Errorf("error casting content block")
			}
			responseParts = append(responseParts, protocol.NewTextPart(content.Value))
			converseLoop = false

		case types.StopReasonToolUse:
			err := nable.HandleToolUse(converseOutput.Output, &converseInput.Messages)
			if err != nil {
				return fmt.Errorf("tool call failed: %w", err)
			}
			continue

		case types.StopReasonMaxTokens:
		case types.StopReasonContentFiltered:
		case types.StopReasonGuardrailIntervened:
		default:
			return fmt.Errorf("unsupported stop reason")
		}
	}

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
