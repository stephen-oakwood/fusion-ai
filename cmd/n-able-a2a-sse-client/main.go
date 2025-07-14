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

	var currentTaskID *string

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

		params := createMessageParams(input, contextID, currentTaskID, 0)

		currentTaskID = handleStreamingInteraction(a2aClient, params)

	}

}

func createMessageParams(input string, contextID string, taskID *string, historyLength int) protocol.SendMessageParams {
	message := protocol.NewMessageWithContext(
		protocol.MessageRoleUser,
		[]protocol.Part{protocol.NewTextPart(input)},
		taskID,
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

	if message.TaskID != nil {
		fmt.Printf("Sending Message ID %s for task %s", message.MessageID, *message.TaskID)
	} else {
		fmt.Printf("Sending Message ID %s with no task", message.MessageID)
	}

	return params
}

func handleStreamingInteraction(a2aClient *client.A2AClient, params protocol.SendMessageParams) *string {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	eventChan, streamErr := a2aClient.StreamMessage(ctx, params)
	if streamErr != nil {
		fmt.Printf("Stream Message Request Failed %s", streamErr)
		return nil
	}

	taskID := processStreamResponse(ctx, eventChan)

	fmt.Printf("Stream processing finished for message %s", params.Message.MessageID)
	fmt.Println(strings.Repeat("-", 60))

	return taskID
}

func processStreamResponse(ctx context.Context, eventChan <-chan protocol.StreamingMessageEvent) *string {
	fmt.Println("\nAgent Response Stream:")
	fmt.Println(strings.Repeat("-", 60))

	var taskID string

	for {

		select {
		case <-ctx.Done():
			fmt.Printf("Context timeout or cancellation while waiting for stream events: %s", ctx.Err())
			return nil

		case event, ok := <-eventChan:
			if !ok {
				fmt.Println("Stream channel closed")
				if ctx.Err() != nil {
					fmt.Printf("Context error after stream close: %s", ctx.Err())
				}
				return nil
			}

			switch e := event.Result.(type) {
			case *protocol.Message:
				fmt.Println("[Message Response:]")
				printMessage(*e)
			case *protocol.Task:
				taskID = e.ID
				fmt.Printf("[Task %s State: %s]\n", e.ID, e.Status.State)
				if e.Status.Message != nil {
					printMessage(*e.Status.Message)
				}
			case *protocol.TaskStatusUpdateEvent:
				taskID = e.TaskID
				fmt.Printf("[Status Update: Task %s, State %s]\n", e.TaskID, e.Status.State)
				if e.Status.Message != nil {
					printMessage(*e.Status.Message)
				}

				if e.Status.State == protocol.TaskStateInputRequired {
					fmt.Println("[Additional Input Required]")
					return &taskID
				} else if e.IsFinal() {
					fmt.Printf("Final status received: %s", e.Status.State)

					if e.Status.State == protocol.TaskStateCompleted {
						fmt.Println("  [Task completed successfully]")
					} else if e.Status.State == protocol.TaskStateFailed {
						fmt.Println("  [Task failed]")
					} else if e.Status.State == protocol.TaskStateCanceled {
						fmt.Println("  [Task was canceled]")
					}
					return nil

				}
			case *protocol.TaskArtifactUpdateEvent:
				taskID = e.TaskID
				name := getArtifactName(e.Artifact)
				fmt.Printf("[Artifact Update: Task %s, Name %s]\n", e.TaskID, name)

				printParts(e.Artifact.Parts)

				if e.LastChunk != nil && *e.LastChunk {
					fmt.Printf("Final artefact received with ID %s", e.Artifact.ArtifactID)
				}

			default:
				fmt.Println("Warning: received unknown event type: %T\n ", event)
			}
		}

	}
}

func getArtifactName(artifact protocol.Artifact) string {
	if artifact.Name != nil {
		return *artifact.Name
	}
	return fmt.Sprintf("Artifact %s", artifact.ArtifactID)
}

func printMessage(message protocol.Message) {
	printParts(message.Parts)
}

func printParts(parts []protocol.Part) {
	for _, part := range parts {
		printPart(part)
	}
}

func printPart(part interface{}) {
	indent := ""

	switch p := part.(type) {
	case *protocol.TextPart:
		fmt.Println(indent + p.Text)
	case *protocol.FilePart:
	case *protocol.DataPart:
	default:
		fmt.Printf("%sUnsupported part type: %T\n", indent, p)
	}
}
