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

	taskID := uuid.New().String()
	sessionID := uuid.New().String()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter text to send to the agent.")
	fmt.Println(strings.Repeat("-", 60))

	for {

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

	fmt.Println("\n<< Agent Response:")
	fmt.Println(strings.Repeat("-", 10))
	fmt.Printf("  State: %s (%s)\n", task.Status.State, task.Status.Timestamp)

	if task.Status.Message != nil {
		fmt.Println("  Message:")
		printParts(task.Status.Message.Parts)
	}

	if len(task.Artifacts) > 0 {
		fmt.Println("  Artifacts:")
		for i, artifact := range task.Artifacts {
			name := fmt.Sprintf("Artifact #%d", i+1)
			if artifact.Name != nil {
				name = *artifact.Name
			}
			fmt.Printf("    [%s]\n", name)
			printParts(artifact.Parts)
		}
	}

	if task.Status.State == protocol.TaskStateInputRequired {
		fmt.Println("  [Additional input required]")
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
