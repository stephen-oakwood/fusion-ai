package a2a

import (
	"context"
	"fmt"
	"fusion/internal/tools"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"time"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

//const systemPrompt = `
//		"You are an IT Technician, capable of providing detailed answers to the questions that your customers ask regarding their assets." +
//		"Think before you reply. Inform the customer of each step you are going to take." +
//		"First, determine if there are any knowledge articles related to the question that could help with your reply." +
//		"Second, using available tools, fetch the schema for a graphQL API that provides the ability to search for assets and return their details." +
//		"Third, using the fetched schema, use the graphQL API to collect data from the customer's assets that are relevant to the points raised in the knowledge articles.'" +
//		"Finally, provide a detailed summary that includes the relevant points from the knowledge articles with examples of customer's assets that demonstrate these points. Explicitly name the assets when providing examples." +
//		"Always fetch the graphQL API Schema first, and construct queries using this schema. Do not construct queries without using the schema."
//`

const systemPrompt = `
	"You are an IT Technician, capable of providing detailed answers to the questions that your customers ask regarding their assets." +
	"Think before you reply. Inform the customer of each step you are going to take." +
	"Tools are available that allow searching knowledge bases, fetching details of an asset," +
	"obtaining the Schema of a GraphQL API that provides the ability to search and manage assets, plus a tool to invoke generated GraphQL Queries for this API."
`

type assetManagementAgent struct {
	ModelClient *modelClient
	Token       string
}

func NewAgent(token string) (*assetManagementAgent, error) {
	modelClient, err := NewModelClient()
	if err != nil {
		return nil, err
	}

	return &assetManagementAgent{
		ModelClient: modelClient,
		Token:       token,
	}, nil
}

func (p *assetManagementAgent) ProcessMessage(ctx context.Context, message protocol.Message, options taskmanager.ProcessOptions, handle taskmanager.TaskHandler) (*taskmanager.MessageProcessingResult, error) {
	inputText := extractText(message)

	if inputText == "" {
		fmt.Printf("process message - input string must contain text")
		errMsg := protocol.NewMessage(
			protocol.MessageRoleAgent,
			[]protocol.Part{protocol.NewTextPart("input message must contain text")},
		)

		return &taskmanager.MessageProcessingResult{
			Result: &errMsg,
		}, nil
	}

	specificTaskID := message.TaskID
	taskID, err := handle.BuildTask(specificTaskID, message.ContextID)
	if err != nil {
		return nil, fmt.Errorf("process message - failed to create task: %w", err)
	}

	if options.Streaming {
		return p.processStreamingMode(ctx, inputText, message.ContextID, taskID, handle)
	}

	return p.processNonStreamingMode(ctx, inputText, message.ContextID, taskID, handle)

}

func (p *assetManagementAgent) processStreamingMode(ctx context.Context, inputText string, contextID *string, taskID string, handle taskmanager.TaskHandler) (*taskmanager.MessageProcessingResult, error) {

	subscriber, err := handle.SubScribeTask(&taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to task: %w", err)
	}

	go p.processRequest(ctx, inputText, contextID, taskID, handle)

	return &taskmanager.MessageProcessingResult{
		StreamingEvents: subscriber,
	}, nil
}

func (p *assetManagementAgent) processNonStreamingMode(ctx context.Context, inputText string, contextID *string, taskID string, handle taskmanager.TaskHandler) (*taskmanager.MessageProcessingResult, error) {

	p.processRequest(ctx, inputText, contextID, taskID, handle)

	cancellable, err := handle.GetTask(&taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &taskmanager.MessageProcessingResult{
		Result: cancellable.Task(),
	}, nil
}

func (p *assetManagementAgent) processRequest(ctx context.Context, inputText string, contextID *string, taskID string, handle taskmanager.TaskHandler) {
	var temperature float32 = 0.0

	toolConfig := tools.ToolConfig(handle, taskID, contextID, p.Token)
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
						Value: inputText,
					},
				},
				Role: types.ConversationRoleUser,
			},
		},
		ToolConfig: &toolConfig,
		InferenceConfig: &types.InferenceConfiguration{
			Temperature: &temperature,
		},
	}

	converseLoop := true

	for converseLoop {
		converseOutput, err := p.ModelClient.Converse(ctx, converseInput)
		if err != nil {
			err = handle.UpdateTaskState(&taskID, protocol.TaskStateFailed, nil)
			if err != nil {
				fmt.Errorf("failed to update task status to failed: %w", err)
			}
			return
		}

		converseMessage, ok := converseOutput.Output.(*types.ConverseOutputMemberMessage)
		if !ok {
			err = handle.UpdateTaskState(&taskID, protocol.TaskStateFailed, nil)
			if err != nil {
				fmt.Errorf("failed to extract converse message: %w", err)
			}
			return
		}

		switch converseOutput.StopReason {
		case types.StopReasonEndTurn:
			content, ok := converseMessage.Value.Content[0].(*types.ContentBlockMemberText)
			if !ok {
				err = handle.UpdateTaskState(&taskID, protocol.TaskStateFailed, nil)
				if err != nil {
					fmt.Errorf("failed to extract end turn converse message: %w", err)
				}
				return
			}

			artifact := protocol.Artifact{
				ArtifactID:  protocol.GenerateArtifactID(),
				Name:        stringPtr("Final Response"),
				Description: stringPtr("Response from model"),
				Parts:       []protocol.Part{protocol.NewTextPart(content.Value)},
				Metadata: map[string]interface{}{
					"processedAt": time.Now().UTC().Format(time.RFC3339),
				},
			}

			err = handle.AddArtifact(&taskID, artifact, true, false)
			if err != nil {
				fmt.Errorf("failed to send artifact event: %w", err)
				return
			}

			err = handle.UpdateTaskState(&taskID, protocol.TaskStateCompleted, nil)
			if err != nil {
				fmt.Errorf("failed to send completed event: %w", err)
				return
			}

			converseLoop = false

		case types.StopReasonToolUse:

			for _, item := range converseMessage.Value.Content {
				switch d := item.(type) {
				case *types.ContentBlockMemberText:

					err = handle.UpdateTaskState(&taskID, protocol.TaskStateWorking, &protocol.Message{
						ContextID: contextID,
						MessageID: protocol.GenerateMessageID(),
						Role:      protocol.MessageRoleAgent,
						Parts:     []protocol.Part{protocol.NewTextPart(d.Value)},
					})
					if err != nil {
						fmt.Errorf("failed to send progress event: %w", err)
						return
					}
				}
			}

			err := tools.HandleToolUse(converseOutput.Output, &converseInput.Messages, p.Token)

			if err != nil {
				err = handle.UpdateTaskState(&taskID, protocol.TaskStateFailed, nil)
				if err != nil {
					fmt.Errorf("failed to send failed event after tool use: %w", err)
				}
				return
			}
			continue

		case types.StopReasonMaxTokens:
		case types.StopReasonContentFiltered:
		case types.StopReasonGuardrailIntervened:
		default:
			fmt.Errorf("unsupported stop reason")
			return
		}
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func extractText(message protocol.Message) string {
	var inputText string
	for _, part := range message.Parts {
		if textPart, ok := part.(*protocol.TextPart); ok {
			inputText += textPart.Text
		}
	}
	return inputText
}
