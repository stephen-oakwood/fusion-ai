package a2a

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
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
