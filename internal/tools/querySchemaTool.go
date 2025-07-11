package tools

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"strings"
)

type QuerySchemaTool struct {
	Name        string
	Description string
}

var querySchemaTool = &QuerySchemaTool{
	Name:        "query_schema",
	Description: "Returns a partial GraphQL Schema that can be used to construct queries and mutations for an API that supports searching for managed assets and returning details regarding their operating systems, hardware and more.",
}

func NewQuerySchemaTool() *QuerySchemaTool {
	return querySchemaTool
}

func (t *QuerySchemaTool) GenerateToolSchema() *types.ToolMemberToolSpec {

	return &types.ToolMemberToolSpec{
		Value: types.ToolSpecification{
			Name:        aws.String(t.Name),
			Description: aws.String(t.Description),
			InputSchema: &types.ToolInputSchemaMemberJson{
				Value: document.NewLazyDocument(map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []interface{}{},
				}),
			},
		},
	}
}

func (t *QuerySchemaTool) Call(toolCall *types.ContentBlockMemberToolUse) (*types.Message, error) {

	schema := strings.TrimSpace(AssetsSchema)
	content := document.NewLazyDocument(map[string]interface{}{"schema": schema})

	return &types.Message{
		Role: types.ConversationRoleUser,
		Content: []types.ContentBlock{
			&types.ContentBlockMemberToolResult{
				Value: types.ToolResultBlock{
					Content: []types.ToolResultContentBlock{
						&types.ToolResultContentBlockMemberJson{
							Value: content,
						},
					},
					ToolUseId: toolCall.Value.ToolUseId,
				},
			},
		},
	}, nil
}
