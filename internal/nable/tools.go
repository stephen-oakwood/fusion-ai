package nable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	Token = "eyJhbGciOiJSUzI1NiIsImtpZCI6IkYxOTZFODUxMUI1RjUyMjhCREMxNDRGNkUxMEFFOTlGQ0EzQ0MzRkZSUzI1NiIsIng1dCI6IjhaYm9VUnRmVWlpOXdVVDI0UXJwbjhvOHdfOCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL3Nzby5zdGcuc2hhcmVkc3Zjcy5zeXN0ZW0tbW9uaXRvci5jb20iLCJuYmYiOjE3NDg1MzE3MTYsImlhdCI6MTc0ODUzMTcxNiwiZXhwIjoxNzQ4NTM1MzE2LCJhdWQiOiJodHRwczovL3Nzby5zdGcuc2hhcmVkc3Zjcy5zeXN0ZW0tbW9uaXRvci5jb20vcmVzb3VyY2VzIiwic2NvcGUiOlsib3BlbmlkIiwiZW1haWwiLCJwcm9maWxlIiwib2ZmbGluZV9hY2Nlc3MiXSwiYW1yIjpbInB3ZCJdLCJjbGllbnRfaWQiOiJkOTE3Y2UzYS1hM2Q0LTQ5MzMtODgwNC00NGI5NDY2MTFkMTkiLCJzdWIiOiI0ODc3ZjZmYi0xMDQ2LTQxZDktODU5NS01MmJmMTYxYzI3YTkiLCJhdXRoX3RpbWUiOjE3NDg1MjQ4NDksImlkcCI6ImxvY2FsIiwic2lkIjoiMDQ2M0FDNjUwQjZEOTY5RjVDNTk2MDEzOUNDRDFGMEYiLCJqdGkiOiIwNEZCMTI2OTg4MkNFQkYyNjhBMEEwQzlGQkYxNTM2OSJ9.jR1014QB8wM0WtOexWmGI14bnWIdHr9-BNI-czg8dHucw3o91pD0uZfs4M25S9uY68d5B1jXB0mfwY9-WoEWlNjQz-x9TBJWKjTe5sXJjXJ9kaGj_PdKQk7f6CGI4NyE3OxhmogGW_9lIWR9Tur9TLHPWqBmO2WhoK16sJFSPcbgUF0Q7LeDoJRB9e2BCfOAFyIn-5l6mqa7xJ6KdTR39bmnE3RjVBz28Wav3ypQAHsw--tr4RBJJdqMYf6ea2JhHPN4YJ8qSVgjTO-oc9BPJEBC86X2rGbBDY6r06MU_1nvyZhB_ekSBzNiGzhw4ZGpLJldhd8xz3QQMQAHTF889A"
)

type Tools interface {
	GenerateToolSchema() *types.ToolMemberToolSpec
	GenerateToolResult(string, map[string]interface{}) (*types.Message, error)
	Invoke(string, string, map[string]interface{}) (*types.Message, error)
}

type QuerySchemaTool struct {
	Name        string
	Description string
}

var querySchemaTool = &QuerySchemaTool{
	Name:        "query_schema",
	Description: "Returns a partial GraphQL Schema that can be used to construct queries and mutations for an API that supports searching for managed assets and returning details regarding their operating systems, hardware and more.",
}

type ExecuteQueryTool struct {
	Name        string
	Description string
}

var executeQueryTool = &ExecuteQueryTool{
	Name:        "execute_query",
	Description: "Executes a GraphQL Query using the n-able public API. Provides support for sophisticated searching of assets.",
}

type KnowledgeQueryTool struct {
	Name        string
	Description string
}

var knowledgeQueryTool = &KnowledgeQueryTool{
	Name:        "knowledge_query",
	Description: "Finds the most relevant knowledge article for a user's questions about managed assets",
}

type body struct {
	Query string `json:"query"`
}

func GetQuerySchemaTool() *QuerySchemaTool {
	return querySchemaTool
}

func GetExecuteQueryTool() *ExecuteQueryTool {
	return executeQueryTool
}

func GetKnowledgeQueryTool() *KnowledgeQueryTool {
	return knowledgeQueryTool
}

func (t *QuerySchemaTool) GenerateToolSchema() *types.ToolMemberToolSpec {
	getInputSchema := map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
		"required":   []interface{}{},
	}

	getInputSchemaDoc := document.NewLazyDocument(getInputSchema)
	return &types.ToolMemberToolSpec{
		Value: types.ToolSpecification{
			Name:        aws.String(t.Name),
			Description: aws.String(t.Description),
			InputSchema: &types.ToolInputSchemaMemberJson{
				Value: getInputSchemaDoc,
			},
		},
	}
}

func (t *QuerySchemaTool) GenerateToolResult(toolUseID string, result map[string]interface{}) (*types.Message, error) {

	data := make(map[string]interface{})
	data["schema"] = result
	content := document.NewLazyDocument(data)

	return &types.Message{
		Role: "user",
		Content: []types.ContentBlock{
			&types.ContentBlockMemberToolResult{
				Value: types.ToolResultBlock{
					ToolUseId: &toolUseID,
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

func (t *QuerySchemaTool) Invoke(toolUseId string, toolName string, parameters map[string]interface{}) (*types.Message, error) {

	if toolName != t.Name {
		return nil, fmt.Errorf("%s not match %s", toolName, t.Name)
	}

	result := strings.TrimSpace(AssetsSchema)
	resultValue := map[string]interface{}{"schema": result}

	return t.GenerateToolResult(toolUseId, resultValue)

}

func (t *KnowledgeQueryTool) GenerateToolSchema() *types.ToolMemberToolSpec {
	knowledgeQuestion := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"question": map[string]interface{}{
				"type":        "string",
				"description": "User provided question regarded managed assets that is used to find the most relevant knowledge article.",
			},
		},
		"required": []interface{}{"question"},
	}

	knowledgeQuestionDoc := document.NewLazyDocument(knowledgeQuestion)
	return &types.ToolMemberToolSpec{
		Value: types.ToolSpecification{
			Name:        aws.String(t.Name),
			Description: aws.String(t.Description),
			InputSchema: &types.ToolInputSchemaMemberJson{
				Value: knowledgeQuestionDoc,
			},
		},
	}
}

func (t *KnowledgeQueryTool) GenerateToolResult(toolUseID string, result map[string]interface{}) (*types.Message, error) {

	data := make(map[string]interface{})
	data["knowledge"] = result
	content := document.NewLazyDocument(data)

	return &types.Message{
		Role: "user",
		Content: []types.ContentBlock{
			&types.ContentBlockMemberToolResult{
				Value: types.ToolResultBlock{
					ToolUseId: &toolUseID,
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

func (t *KnowledgeQueryTool) Invoke(toolUseId string, toolName string, parameters map[string]interface{}) (*types.Message, error) {

	if toolName != t.Name {
		return nil, fmt.Errorf("%s not match %s", toolName, t.Name)
	}

	result := "There could be several reasons why your devices are running slowly. Here are some common causes and potential solutions:\n\nInsufficient system resources (RAM and CPU):\n\nClose unnecessary applications and browser tabs to free up memory.\nConsider upgrading your device's RAM if it's running low on memory.\nCheck for any resource-intensive processes or programs that may be consuming a lot of CPU power.\nHard disk drive (HDD) issues:\n\nIf your device has a traditional hard disk drive (HDD), it may be slowing down due to fragmentation or lack of free space.\nRun a disk defragmentation tool to optimize the file system.\nDelete unnecessary files and programs to free up disk space.\nConsider upgrading to a solid-state drive (SSD) for faster read/write speeds.\nSoftware issues:\n\nOutdated or bloated software can consume system resources and cause slowdowns.\nUpdate your operating system, drivers, and applications to the latest versions.\nUninstall any unnecessary programs or bloatware that may be running in the background.\nMalware or virus infections:\n\nMalware or viruses can significantly impact system performance.\nRun a full system scan with a reliable anti-virus/anti-malware program to detect and remove any threats.\nOverheating issues:\n\nOverheating can cause your device to throttle its performance to prevent damage.\nClean out any dust buildup and ensure proper ventilation for your device.\nCheck if the cooling fans are working correctly.\nHardware aging:\n\nIf your device is several years old, the hardware components may be reaching the end of their lifespan, resulting in slower performance.\nConsider upgrading to a newer device or replacing specific components, such as RAM or storage drives.\nTo identify the root cause, you can use system monitoring tools, check the Task Manager (Windows) or Activity Monitor (macOS) to see what processes are consuming resources, and perform basic maintenance tasks like disk cleanup and defragmentation."
	resultValue := map[string]interface{}{"knowledge": result}

	return t.GenerateToolResult(toolUseId, resultValue)

}

func HandleToolUse(output types.ConverseOutput, messages *[]types.Message) error {
	switch v := output.(type) {
	case *types.ConverseOutputMemberMessage:
		*messages = append(*messages, v.Value)

		for _, item := range v.Value.Content {
			switch d := item.(type) {
			case *types.ContentBlockMemberText:
				fmt.Printf("Handle Tool Use: Content Block Member Text: %s\n", d.Value)
			case *types.ContentBlockMemberToolUse:
				fmt.Printf("Handle Tool Use: Content Block Member Tool Use: %s\n", *d.Value.Name)

				if *d.Value.Name == "query_schema" {
					data := make(map[string]interface{})
					err := d.Value.Input.UnmarshalSmithyDocument(&data)

					if err == nil {
						message, err := querySchemaTool.Invoke(*d.Value.ToolUseId, *d.Value.Name, data)
						if err != nil {
							fmt.Printf("Error invoking tool: %v\n", err)
							return err
						} else {
							*messages = append(*messages, *message)
						}
					}
				}

				if *d.Value.Name == "execute_query" {
					data := make(map[string]interface{})
					err := d.Value.Input.UnmarshalSmithyDocument(&data)

					if err == nil {
						message, err := executeQueryTool.Invoke(*d.Value.ToolUseId, *d.Value.Name, data)
						if err != nil {
							fmt.Printf("Error invoking tool: %v\n", err)
							return err
						} else {
							*messages = append(*messages, *message)
						}
					}
				}

				if *d.Value.Name == "knowledge_query" {
					data := make(map[string]interface{})
					err := d.Value.Input.UnmarshalSmithyDocument(&data)

					if err == nil {
						message, err := knowledgeQueryTool.Invoke(*d.Value.ToolUseId, *d.Value.Name, data)
						if err != nil {
							fmt.Printf("Error invoking tool: %v\n", err)
							return err
						} else {
							*messages = append(*messages, *message)
						}
					}
				}
			}
		}
	default:
		fmt.Println("Response is nil or unknown type")
	}

	return nil
}

func (t *ExecuteQueryTool) GenerateToolSchema() *types.ToolMemberToolSpec {
	getInputSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "A GraphQL Query that will be executed to search for managed assets and return their details.",
			},
		},
		"required": []interface{}{"query"},
	}

	getInputSchemaDoc := document.NewLazyDocument(getInputSchema)
	return &types.ToolMemberToolSpec{
		Value: types.ToolSpecification{
			Name:        aws.String(t.Name),
			Description: aws.String(t.Description),
			InputSchema: &types.ToolInputSchemaMemberJson{
				Value: getInputSchemaDoc,
			},
		},
	}
}

func (t *ExecuteQueryTool) GenerateToolResult(toolUseID string, result map[string]interface{}) (*types.Message, error) {

	data := make(map[string]interface{})
	data["assets"] = result
	content := document.NewLazyDocument(data)

	return &types.Message{
		Role: "user",
		Content: []types.ContentBlock{
			&types.ContentBlockMemberToolResult{
				Value: types.ToolResultBlock{
					ToolUseId: &toolUseID,
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

func (t *ExecuteQueryTool) Invoke(toolUseId string, toolName string, parameters map[string]interface{}) (*types.Message, error) {

	if toolName != t.Name {
		return nil, fmt.Errorf("%s not match %s", toolName, t.Name)
	}

	query, hasQuery := parameters["query"].(string)
	if !hasQuery {
		fmt.Println("No query!!!")
		return nil, fmt.Errorf("no query")
	}

	result, err := executeQuery(query)
	if err != nil {
		fmt.Println("Failed to execute query")
		return nil, fmt.Errorf("query execution")
	}

	resultValue := map[string]interface{}{"assets": result}

	return t.GenerateToolResult(toolUseId, resultValue)

}

func PrintJSON(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}
	log.Printf("\nJSON:\n%s\n", string(jsonBytes))
}

func executeQuery(query string) (string, error) {

	payload := &body{Query: query}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal the query %s\n", err)
	}

	bearer := "Bearer " + Token
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
