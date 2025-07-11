package tools

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type AssetDetailsTool struct {
	Name        string
	Description string
	Token       string
}

var assetDetailsTool = &AssetDetailsTool{
	Name:        "asset_details",
	Description: "Provides the details of an asset e.g. name, owner, operating system, hardware details, etc",
}

//go:embed asset_details.gql
var assetDetailsQuery string

func NewAssetDetailsTool(token string) *AssetDetailsTool {
	assetDetailsTool.Token = token
	return assetDetailsTool
}

func (t *AssetDetailsTool) GenerateToolSchema() *types.ToolMemberToolSpec {

	return &types.ToolMemberToolSpec{
		Value: types.ToolSpecification{
			Name:        aws.String(t.Name),
			Description: aws.String(t.Description),
			InputSchema: &types.ToolInputSchemaMemberJson{
				Value: document.NewLazyDocument(map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"assetId": map[string]interface{}{
							"type":        "string",
							"description": "User provided identifier for the asset.",
						},
					},
					"required": []interface{}{"assetId"},
				}),
			},
		},
	}
}

func (t *AssetDetailsTool) Call(toolCall *types.ContentBlockMemberToolUse) (*types.Message, error) {
	var parameters map[string]interface{}

	if toolCall.Value.Input != nil {
		err := toolCall.Value.Input.UnmarshalSmithyDocument(&parameters)
		if err != nil {
			fmt.Errorf("tool call failed. unable to marshal parameters: %w", err)
			return nil, err
		}
	}

	if parameters == nil || parameters["assetId"] == nil {
		err := errors.New("tool call failed. missing input parameter")
		fmt.Errorf("tool call failed. no asset ID provided by the model: %w", err)
		return nil, err
	}

	assetId, ok := parameters["assetId"].(string)
	if !ok || assetId == "" {
		err := errors.New("tool call failed. invalid parameter")
		fmt.Errorf("tool call failed. invalid parameter: %w", err)
		return nil, err
	}

	variables := map[string]interface{}{
		"id": assetId,
	}

	result, err := executeQuery(graphQLBody{Query: assetDetailsQuery, Variables: variables}, t.Token)
	if err != nil {
		fmt.Println("Failed to execute query")
		return nil, fmt.Errorf("query execution")
	}

	content := document.NewLazyDocument(map[string]interface{}{"asset": result})

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
