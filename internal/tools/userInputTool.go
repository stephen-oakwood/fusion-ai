package tools

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

type UserInputTool struct {
	Name        string
	Description string
	Handle      taskmanager.TaskHandler
	TaskID      string
	ContextID   *string
}

var userInputTool = &UserInputTool{
	Name:        "user_input_required",
	Description: "This tool can be used when you require input from the user before proceeding with a task.",
}

func NewUserInputTool(handle taskmanager.TaskHandler, taskID string, contextID *string) *UserInputTool {
	userInputTool.Handle = handle
	userInputTool.TaskID = taskID
	userInputTool.ContextID = contextID
	return userInputTool
}

func (t *UserInputTool) GenerateToolSchema() *types.ToolMemberToolSpec {
	return &types.ToolMemberToolSpec{
		Value: types.ToolSpecification{
			Name:        &t.Name,
			Description: &t.Description,
			InputSchema: &types.ToolInputSchemaMemberJson{
				Value: document.NewLazyDocument(map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "The input required from the user to proceed with the task.",
						},
					},
					"required": []interface{}{
						"reason",
					},
				}),
			},
		},
	}
}

func (t *UserInputTool) Call(toolCall *types.ContentBlockMemberToolUse) (*types.Message, error) {
	var parameters map[string]interface{}

	if toolCall.Value.Input != nil {
		err := toolCall.Value.Input.UnmarshalSmithyDocument(&parameters)
		if err != nil {
			fmt.Errorf("tool call failed. unable to marshal parameters: %w", err)
			return nil, err
		}
	}

	if parameters == nil || parameters["reason"] == nil {
		err := errors.New("tool call failed. missing input parameter")
		fmt.Errorf("tool call failed. no reason provided by the model: %w", err)
		return nil, err
	}

	reason, ok := parameters["reason"].(string)
	if !ok || reason == "" {
		err := errors.New("tool call failed. invalid parameter")
		fmt.Errorf("tool call failed. invalid parameter: %w", err)
		return nil, err
	}

	err := t.Handle.UpdateTaskState(&t.TaskID, protocol.TaskStateInputRequired, &protocol.Message{
		ContextID: t.ContextID,
		MessageID: protocol.GenerateMessageID(),
		Role:      protocol.MessageRoleAgent,
		Parts:     []protocol.Part{protocol.NewTextPart(reason)},
	})
	if err != nil {
		fmt.Errorf("tool call failed. unable to update task: %w", err)
		return nil, err
	}

	return &types.Message{
		Role: types.ConversationRoleUser,
		Content: []types.ContentBlock{
			&types.ContentBlockMemberToolResult{
				Value: types.ToolResultBlock{
					Content: []types.ToolResultContentBlock{
						&types.ToolResultContentBlockMemberText{
							Value: "Successfully set task status to input required",
						},
					},
					ToolUseId: toolCall.Value.ToolUseId,
				},
			},
		},
	}, nil
}
