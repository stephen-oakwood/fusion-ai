package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"trpc.group/trpc-go/trpc-a2a-go/taskmanager"
)

type Tool interface {
	GenerateToolSchema() *types.ToolMemberToolSpec
	Call(toolCall *types.ContentBlockMemberToolUse) (*types.Message, error)
}

var Tools map[string]Tool

func ToolConfig(handle taskmanager.TaskHandler, taskID string, contextID *string, token string) types.ToolConfiguration {
	querySchemaTool := NewQuerySchemaTool()
	querySchemaToolSchema := querySchemaTool.GenerateToolSchema()

	executeQueryTool := NewExecuteQueryTool(token)
	executeQueryToolSchema := executeQueryTool.GenerateToolSchema()

	knowledgeQueryTool := NewKnowledgeQueryTool()
	knowledgeQueryToolSchema := knowledgeQueryTool.GenerateToolSchema()

	assetDetailsTool := NewAssetDetailsTool(token)
	assetDetailsToolSchema := assetDetailsTool.GenerateToolSchema()

	userInputTool := NewUserInputTool(handle, taskID, contextID)
	userInputToolSchema := userInputTool.GenerateToolSchema()

	Tools = map[string]Tool{
		querySchemaTool.Name:    querySchemaTool,
		executeQueryTool.Name:   executeQueryTool,
		knowledgeQueryTool.Name: knowledgeQueryTool,
		assetDetailsTool.Name:   assetDetailsTool,
		userInputTool.Name:      userInputTool,
	}

	return types.ToolConfiguration{
		Tools: []types.Tool{
			querySchemaToolSchema,
			executeQueryToolSchema,
			knowledgeQueryToolSchema,
			assetDetailsToolSchema,
			userInputToolSchema,
		},
	}
}

type graphQLBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

func executeQuery(payload graphQLBody, token string) (string, error) {

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal the query %s\n", err)
	}

	bearer := "Bearer " + token
	request, err := http.NewRequest(http.MethodPost, "https://stg.api.n-able.com/graphql", bytes.NewBuffer(payloadJSON))
	request.Header.Add("Authorization", bearer)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The N-Query Client Request failed with error %s\n", err)
		return "", err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("The N-Query reading response failed with error %s\n", err)
		return "", err
	}

	return string(data), nil
}

func HandleToolUse(output types.ConverseOutput, messages *[]types.Message, token string) error {
	switch v := output.(type) {
	case *types.ConverseOutputMemberMessage:
		*messages = append(*messages, v.Value)

		for _, item := range v.Value.Content {
			switch contentBlock := item.(type) {
			case *types.ContentBlockMemberText:
				fmt.Printf("Handle Tool Use: Content Block Member Text: %s\n", contentBlock.Value)
			case *types.ContentBlockMemberToolUse:
				fmt.Printf("Handle Tool Use: Content Block Member Tool Use: %s\n", *contentBlock.Value.Name)

				if *contentBlock.Value.Name == "query_schema" {
					message, err := Tools["query_schema"].Call(contentBlock)
					if err != nil {
						fmt.Printf("Error invoking tool: %v\n", err)
						return err
					} else {
						*messages = append(*messages, *message)
					}
				}

				if *contentBlock.Value.Name == "execute_query" {
					message, err := Tools["execute_query"].Call(contentBlock)
					if err != nil {
						fmt.Printf("Error invoking tool: %v\n", err)
						return err
					} else {
						*messages = append(*messages, *message)
					}
				}

				if *contentBlock.Value.Name == "knowledge_query" {
					message, err := Tools["knowledge_query"].Call(contentBlock)
					if err != nil {
						fmt.Printf("Error invoking tool: %v\n", err)
						return err
					} else {
						*messages = append(*messages, *message)
					}
				}

				if *contentBlock.Value.Name == "asset_details" {
					message, err := Tools["asset_details"].Call(contentBlock)
					if err != nil {
						fmt.Printf("Error invoking tool: %v\n", err)
						return err
					} else {
						*messages = append(*messages, *message)
					}
				}

				if *contentBlock.Value.Name == "user_input_required" {
					message, err := Tools["user_input_required"].Call(contentBlock)
					if err != nil {
						fmt.Printf("Error invoking tool: %v\n", err)
						return err
					} else {
						*messages = append(*messages, *message)
					}
				}
			}
		}
	default:
		fmt.Println("Response is nil or unknown type")
	}

	return nil
}

func PrintJSON(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}
	log.Printf("\nJSON:\n%s\n", string(jsonBytes))
}
