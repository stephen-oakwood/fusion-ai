package a2a

import (
	"context"
	"errors"
	"fmt"
	"fusion/internal/nable"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type modelClient struct {
	BedrockClient *bedrockruntime.Client
	BedrockModel  string
}

func NewModelClient() (*modelClient, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-1"), config.WithSharedConfigProfile("archpoc_dev-developer"))
	if err != nil {
		fmt.Printf("Couldn't load AWS Config %s", err)
		return nil, err
	}
	client := bedrockruntime.NewFromConfig(awsConfig)

	return &modelClient{
		BedrockClient: client,
		BedrockModel:  "anthropic.claude-3-haiku-20240307-v1:0",
	}, nil
}

func (i *modelClient) Converse(ctx context.Context, converseInput *bedrockruntime.ConverseInput) (*bedrockruntime.ConverseOutput, error) {
	converseOutput, err := i.BedrockClient.Converse(ctx, converseInput)
	if err != nil {
		return nil, err
	}
	return converseOutput, nil
}

func (i *modelClient) MakeInference(ctx context.Context, converseInput *bedrockruntime.ConverseInput) (string, error) {

	converseOutput, err := i.BedrockClient.Converse(ctx, converseInput)
	if err != nil {
		return "", err
	}

	converseMessage, ok := converseOutput.Output.(*types.ConverseOutputMemberMessage)
	if !ok {
		return "", errors.New("error casting LLM response")
	}

	switch converseOutput.StopReason {
	case types.StopReasonEndTurn:
		content, ok := converseMessage.Value.Content[0].(*types.ContentBlockMemberText)
		if !ok {
			return "", errors.New("error casting content block")
		}
		return content.Value, nil

	case types.StopReasonToolUse:
		err := nable.HandleToolUse(converseOutput.Output, &converseInput.Messages)
		if err != nil {
			return "", err
		}
		return i.MakeInference(ctx, converseInput)

	case types.StopReasonMaxTokens:
	case types.StopReasonContentFiltered:
	case types.StopReasonGuardrailIntervened:
	default:
		return "", errors.New("unknown stop reason")
	}

	return "", errors.New("unexpected return")
}

func toolConfig() types.ToolConfiguration {
	querySchemaToolSchema := nable.GetQuerySchemaTool().GenerateToolSchema()
	executeQueryToolSchema := nable.GetExecuteQueryTool().GenerateToolSchema()
	knowledgeQueryToolSchema := nable.GetKnowledgeQueryTool().GenerateToolSchema()

	return types.ToolConfiguration{
		Tools: []types.Tool{
			querySchemaToolSchema,
			executeQueryToolSchema,
			knowledgeQueryToolSchema,
		},
	}
}
