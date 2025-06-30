package main

import (
	"bufio"
	"context"
	"fmt"
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

	contextID := protocol.GenerateContextID()
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

		params := createMessageParams(input, contextID, 0)

		handleStandardInteraction(a2aClient, params)

	}

}

func createMessageParams(input string, contextID string, historyLength int) protocol.SendMessageParams {
	message := protocol.NewMessageWithContext(
		protocol.MessageRoleUser,
		[]protocol.Part{protocol.NewTextPart(input)},
		nil,
		&contextID,
	)

	params := protocol.SendMessageParams{
		Message: message,
	}

	if historyLength > 0 {
		params.Configuration = &protocol.SendMessageConfiguration{
			HistoryLength: &historyLength,
		}
	}

	return params
}

func handleStandardInteraction(a2aClient *client.A2AClient, params protocol.SendMessageParams) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messageResult, err := a2aClient.SendMessage(ctx, params)
	if err != nil {
		fmt.Printf("failed to send message %v", err)
		return
	}

	switch result := messageResult.Result.(type) {
	case *protocol.Message:
		fmt.Println("[Message Response:]")
		printMessage(*result)
	case *protocol.Task:
		fmt.Printf("[Task %s State %s]\n", result.ID, result.Status.State)
		printTaskResult(result)
	}

	fmt.Println(strings.Repeat("-", 60))

}

func printMessage(message protocol.Message) {
	fmt.Printf("Message ID: %s", message.MessageID)
	if message.ContextID != nil {
		fmt.Printf("Context ID: %s", *message.ContextID)
	}
	fmt.Printf("Role: %s", message.Role)

	fmt.Printf("Message parts:")
	for i, part := range message.Parts {
		switch p := part.(type) {
		case *protocol.TextPart:
			fmt.Printf("  Part %d (text): %s", i+1, p.Text)
		case *protocol.FilePart:
			fmt.Printf("  Part %d (file): [file content]", i+1)
		case *protocol.DataPart:
			fmt.Printf("  Part %d (data): %+v", i+1, p.Data)
		default:
			fmt.Printf("  Part %d (unknown): %+v", i+1, part)
		}
	}
}

func printTaskResult(task *protocol.Task) {
	if task.Status.Message != nil {
		fmt.Printf("Task result message:")
		printMessage(*task.Status.Message)
	}

	// Print artifacts if any
	if len(task.Artifacts) > 0 {
		fmt.Printf("Task artifacts:")
		for i, artifact := range task.Artifacts {
			name := "Unnamed"
			if artifact.Name != nil {
				name = *artifact.Name
			}
			fmt.Printf("  Artifact %d: %s", i+1, name)
			for j, part := range artifact.Parts {
				switch p := part.(type) {
				case *protocol.TextPart:
					fmt.Printf("    Part %d (text): %s", j+1, p.Text)
				case *protocol.FilePart:
					fmt.Printf("    Part %d (file): [file content]", j+1)
				case *protocol.DataPart:
					fmt.Printf("    Part %d (data): %+v", j+1, p.Data)
				default:
					fmt.Printf("    Part %d (unknown): %+v", j+1, part)
				}
			}
		}
	}
}
