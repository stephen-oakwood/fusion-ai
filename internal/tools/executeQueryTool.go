package tools

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type ExecuteQueryTool struct {
	Name        string
	Description string
	Token       string
}

var executeQueryTool = &ExecuteQueryTool{
	Name:        "execute_query",
	Description: "Executes a GraphQL Query using the n-able public API. Provides support for sophisticated searching of assets.",
}

func NewExecuteQueryTool(token string) *ExecuteQueryTool {
	executeQueryTool.Token = token
	return executeQueryTool
}

func (t *ExecuteQueryTool) GenerateToolSchema() *types.ToolMemberToolSpec {

	return &types.ToolMemberToolSpec{
		Value: types.ToolSpecification{
			Name:        aws.String(t.Name),
			Description: aws.String(t.Description),
			InputSchema: &types.ToolInputSchemaMemberJson{
				Value: document.NewLazyDocument(map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "A GraphQL Query that will be executed to search for managed assets and return their details.",
						},
					},
					"required": []interface{}{"query"},
				}),
			},
		},
	}
}

func (t *ExecuteQueryTool) Call(toolCall *types.ContentBlockMemberToolUse) (*types.Message, error) {
	var parameters map[string]interface{}

	if toolCall.Value.Input != nil {
		err := toolCall.Value.Input.UnmarshalSmithyDocument(&parameters)
		if err != nil {
			fmt.Errorf("tool call failed. unable to marshal parameters: %w", err)
			return nil, err
		}
	}

	if parameters == nil || parameters["query"] == nil {
		err := errors.New("tool call failed. missing input parameter")
		fmt.Errorf("tool call failed. no query provided by the model: %w", err)
		return nil, err
	}

	query, ok := parameters["query"].(string)
	if !ok || query == "" {
		err := errors.New("tool call failed. invalid parameter")
		fmt.Errorf("tool call failed. invalid parameter: %w", err)
		return nil, err
	}

	result, err := executeQuery(graphQLBody{Query: query}, t.Token)
	if err != nil {
		fmt.Println("Failed to execute query")
		return nil, fmt.Errorf("query execution")
	}

	content := document.NewLazyDocument(map[string]interface{}{"assets": result})

	return &types.Message{
		Role: "user",
		Content: []types.ContentBlock{
			&types.ContentBlockMemberToolResult{
				Value: types.ToolResultBlock{
					ToolUseId: toolCall.Value.ToolUseId,
					Content: []types.ToolResultContentBlock{
						&types.ToolResultContentBlockMemberJson{
							Value: content,
						},
					},
				},
			},
		},
	}, nil
}
