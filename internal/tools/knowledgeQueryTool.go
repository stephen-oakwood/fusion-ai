package tools

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type KnowledgeQueryTool struct {
	Name        string
	Description string
}

var knowledgeQueryTool = &KnowledgeQueryTool{
	Name:        "knowledge_query",
	Description: "Finds the most relevant knowledge article for a user's questions about managed assets",
}

func NewKnowledgeQueryTool() *KnowledgeQueryTool {
	return knowledgeQueryTool
}

func (t *KnowledgeQueryTool) GenerateToolSchema() *types.ToolMemberToolSpec {

	return &types.ToolMemberToolSpec{
		Value: types.ToolSpecification{
			Name:        aws.String(t.Name),
			Description: aws.String(t.Description),
			InputSchema: &types.ToolInputSchemaMemberJson{
				Value: document.NewLazyDocument(map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"question": map[string]interface{}{
							"type":        "string",
							"description": "User provided question regarded managed assets that is used to find the most relevant knowledge article.",
						},
					},
					"required": []interface{}{"question"},
				}),
			},
		},
	}
}

func (t *KnowledgeQueryTool) Call(toolCall *types.ContentBlockMemberToolUse) (*types.Message, error) {

	content := "There could be several reasons why your devices are running slowly. Here are some common causes and potential solutions:\n\nInsufficient system resources (RAM and CPU):\n\nClose unnecessary applications and browser tabs to free up memory.\nConsider upgrading your device's RAM if it's running low on memory.\nCheck for any resource-intensive processes or programs that may be consuming a lot of CPU power.\nHard disk drive (HDD) issues:\n\nIf your device has a traditional hard disk drive (HDD), it may be slowing down due to fragmentation or lack of free space.\nRun a disk defragmentation tool to optimize the file system.\nDelete unnecessary files and programs to free up disk space.\nConsider upgrading to a solid-state drive (SSD) for faster read/write speeds.\nSoftware issues:\n\nOutdated or bloated software can consume system resources and cause slowdowns.\nUpdate your operating system, drivers, and applications to the latest versions.\nUninstall any unnecessary programs or bloatware that may be running in the background.\nMalware or virus infections:\n\nMalware or viruses can significantly impact system performance.\nRun a full system scan with a reliable anti-virus/anti-malware program to detect and remove any threats.\nOverheating issues:\n\nOverheating can cause your device to throttle its performance to prevent damage.\nClean out any dust buildup and ensure proper ventilation for your device.\nCheck if the cooling fans are working correctly.\nHardware aging:\n\nIf your device is several years old, the hardware components may be reaching the end of their lifespan, resulting in slower performance.\nConsider upgrading to a newer device or replacing specific components, such as RAM or storage drives.\nTo identify the root cause, you can use system monitoring tools, check the Task Manager (Windows) or Activity Monitor (macOS) to see what processes are consuming resources, and perform basic maintenance tasks like disk cleanup and defragmentation."

	return &types.Message{
		Role: types.ConversationRoleUser,
		Content: []types.ContentBlock{
			&types.ContentBlockMemberToolResult{
				Value: types.ToolResultBlock{
					Content: []types.ToolResultContentBlock{
						&types.ToolResultContentBlockMemberText{
							Value: content,
						},
					},
					ToolUseId: toolCall.Value.ToolUseId,
				},
			},
		},
	}, nil
}
