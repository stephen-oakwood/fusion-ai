package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
	"trpc.group/trpc-go/trpc-a2a-go/client"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
)

func main() {
	a2aClient, err := client.NewA2AClient("http://localhost:8080", client.WithTimeout(300*time.Second))
	if err != nil {
		fmt.Errorf("failed to create A2A client: %v", err)
	}
	
	sessionID := uuid.New().String()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter text to send to the agent.")
	fmt.Println(strings.Repeat("-", 60))

	for {

		taskID := uuid.New().String()

		fmt.Println("> ")
		input, readErr := reader.ReadString('\n')
		if readErr != nil {
			fmt.Printf("failed to read input %v", readErr)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		params := createTaskParams(taskID, sessionID, input)

		handleStandardInteraction(a2aClient, params, taskID)

	}

}

func createTaskParams(taskID, sessionID, input string) protocol.SendTaskParams {
	message := protocol.NewMessage(
		protocol.MessageRoleUser,
		[]protocol.Part{protocol.NewTextPart(input)},
	)

	params := protocol.SendTaskParams{
		ID:        taskID,
		SessionID: &sessionID,
		Message:   message,
		Metadata:  map[string]any{"invocationKey": "1111-2222-3333-4444"},
	}

	return params
}

func handleStandardInteraction(a2aClient *client.A2AClient, params protocol.SendTaskParams, taskID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	task, err := a2aClient.SendTasks(ctx, params)
	if err != nil {
		fmt.Printf("failed to send tasks %v", err)
		return
	}

	fmt.Printf("  State: %s (%s)\n", task.Status.State, task.Status.Timestamp)

	historyLength := 10

	if task.Status.State == protocol.TaskStateCompleted {
		queryParams := protocol.TaskQueryParams{
			ID:            taskID,
			HistoryLength: &historyLength,
		}

		task, err = a2aClient.GetTasks(ctx, queryParams)
		if err != nil {
			fmt.Printf("failed to get tasks %v", err)
			return
		}

		for i := 0; i < len(task.History)-1; i++ {
			message := task.History[i]
			if message.Role == protocol.MessageRoleAgent {
				fmt.Printf("\033[1;32m%s", "Agent Response")
				fmt.Println(strings.Repeat("-", 60))
				printParts(message.Parts)
				fmt.Printf("\033[0:30m%s", "\n")
			} else {
				fmt.Printf("\033[1;30m%s", "User Response")
				fmt.Println(strings.Repeat("-", 60))
				printParts(message.Parts)
				fmt.Printf("\033[0:30m%s", "\n")
			}
		}

		if len(task.Artifacts) > 0 {
			for _, artifact := range task.Artifacts {
				fmt.Printf("Artifact [%s]\n", *artifact.Name)
				printParts(artifact.Parts)
			}
		}
	}

	fmt.Println(strings.Repeat("-", 60))

}

func printParts(parts []protocol.Part) {
	for _, part := range parts {
		printPart(part)
	}
}

func printPart(part interface{}) {
	indent := "  "

	switch p := part.(type) {
	case protocol.TextPart:
		fmt.Println(indent + p.Text)
	case protocol.FilePart:
	case protocol.DataPart:
	default:
		fmt.Printf("%sUnsupported part type: %T\n", indent, p)
	}
}
