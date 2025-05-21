package main

import (
	"bufio"
	"context"
	"fmt"
	"fusion/internal/nable"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"os"
	"strings"
)

const systemPrompt = `
		"You are an IT Technician, capable of providing detailed answers to the questions that your customers ask regarding their assets." +
		"Think before you reply in <thinking> tags." +
		"First, determine if there are any knowledge articles related to the question that could help with your reply." +
		"Second, using available tools, fetch the schema for a graphQL API that provides the ability to search for assets and return their details." +
		"Third, using the fetched schema, use the graphQL API to collect data from the customer's assets that are relevant to the points raised in the knowledge articles.'" +
		"Finally, provide a detailed summary that includes the relevant points from the knowledge articles with examples of customer's assets that demonstrate these points. Explicitly name the assets when providing examples." +
		"Always fetch the graphQL API Schema first, and construct queries using this schema. Do not construct queries without using the schema."`

func main() {

	reader := bufio.NewReader(os.Stdin)

	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-1"), config.WithSharedConfigProfile("archpoc_dev-developer"))
	if err != nil {
		panic(err)
	}
	client := bedrockruntime.NewFromConfig(awsConfig)
	modelId := "anthropic.claude-3-haiku-20240307-v1:0"

	querySchemaToolSchema := nable.GetQuerySchemaTool().GenerateToolSchema()
	executeQueryToolSchema := nable.GetExecuteQueryTool().GenerateToolSchema()
	knowledgeQueryToolSchema := nable.GetKnowledgeQueryTool().GenerateToolSchema()

	toolConfig := types.ToolConfiguration{
		Tools: []types.Tool{
			querySchemaToolSchema,
			executeQueryToolSchema,
			knowledgeQueryToolSchema,
		},
	}

	fmt.Println("\nEnter your message:")
	inputText, _ := reader.ReadString('\n')
	inputText = strings.TrimSpace(inputText)

	initialMessage := types.Message{
		Role: "user",
		Content: []types.ContentBlock{
			&types.ContentBlockMemberText{
				Value: inputText,
			},
		},
	}
	messages := []types.Message{initialMessage}

	system := []types.SystemContentBlock{
		&types.SystemContentBlockMemberText{
			Value: systemPrompt,
		},
	}

	nable.PrintJSON(messages)

	for {
		response, err := client.Converse(context.Background(), &bedrockruntime.ConverseInput{
			ModelId:    &modelId,
			Messages:   messages,
			ToolConfig: &toolConfig,
			System:     system,
		})
		if err != nil {
			panic(err)
		}

		nable.PrintJSON(response.Output)

		if response.StopReason == types.StopReasonToolUse {
			nable.HandleToolUse(response.Output, &messages)
		}

		if response.StopReason == types.StopReasonEndTurn {
			switch v := response.Output.(type) {
			case *types.ConverseOutputMemberMessage:
				for _, content := range v.Value.Content {
					switch c := content.(type) {
					case *types.ContentBlockMemberText:
						fmt.Printf("%s\n", c.Value)
					}
				}
			}
			return
		}

	}

}
