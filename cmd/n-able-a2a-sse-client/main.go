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

		handleStreamingInteraction(a2aClient, params, taskID)

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
	}

	return params
}

func handleStreamingInteraction(a2aClient *client.A2AClient, params protocol.SendTaskParams, taskID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	eventChan, streamErr := a2aClient.StreamTask(ctx, params)
	if streamErr != nil {
		fmt.Printf("Stream Task Request Failed %s", streamErr)
		return
	}

	processStreamResponse(ctx, eventChan)

	fmt.Printf("Stream processing finished for tasl %s", taskID)
	fmt.Println(strings.Repeat("-", 60))

}

func processStreamResponse(ctx context.Context, eventChan <-chan protocol.TaskEvent) (protocol.TaskState, []protocol.Artifact) {
	fmt.Println("\nAgent Response Stream:")
	fmt.Println(strings.Repeat("-", 60))

	var finalTaskState protocol.TaskState
	finalArtifacts := []protocol.Artifact{}

	for {

		select {
		case <-ctx.Done():
			fmt.Printf("Context timeout or cancellatiom while waiting for stream events: %s", ctx.Err())
			return finalTaskState, finalArtifacts

		case event, ok := <-eventChan:
			if !ok {
				fmt.Println("Stream channel closed")
				if ctx.Err() != nil {
					fmt.Printf("Context error after stream close: %s", ctx.Err())
				}
				return finalTaskState, finalArtifacts
			}

			switch e := event.(type) {
			case protocol.TaskStatusUpdateEvent:
				fmt.Printf("  [Status Update: %s (%s)]\n\n", e.Status.State, e.Status.Timestamp)
				if e.Status.Message != nil {
					printMessage(*e.Status.Message)
				}

				finalTaskState = e.Status.State

				if e.Status.State == protocol.TaskStateInputRequired {
					fmt.Println("[Additional Support Required]")
					return finalTaskState, finalArtifacts
				} else if e.IsFinal() {
					fmt.Printf("Final status received: %s", finalTaskState)

					if e.Status.State == protocol.TaskStateCompleted {
						fmt.Println("  [Task completed successfully]")
					} else if e.Status.State == protocol.TaskStateFailed {
						fmt.Println("  [Task failed]")
					} else if e.Status.State == protocol.TaskStateCanceled {
						fmt.Println("  [Task was canceled]")
					}
					return finalTaskState, finalArtifacts

				}

			case protocol.TaskArtifactUpdateEvent:
				name := getArtifactName(e.Artifact)

				if e.Artifact.Append != nil && *e.Artifact.Append {
					fmt.Printf("  [Artifact Update: %s (Appending)]\n\n", name)
				} else {
					fmt.Printf("  [Artifact Update: %s]\n\n", name)
				}

				printParts(e.Artifact.Parts)

				if e.Artifact.Append != nil && *e.Artifact.Append && len(finalArtifacts) > 0 {

					for i, art := range finalArtifacts {
						if art.Index == e.Artifact.Index {
							// Append parts
							combinedParts := append(art.Parts, e.Artifact.Parts...)
							finalArtifacts[i].Parts = combinedParts

							// Update other fields if needed
							if e.Artifact.Name != nil {
								finalArtifacts[i].Name = e.Artifact.Name
							}
							if e.Artifact.Description != nil {
								finalArtifacts[i].Description = e.Artifact.Description
							}
							if e.Artifact.LastChunk != nil {
								finalArtifacts[i].LastChunk = e.Artifact.LastChunk
							}

							break
						}
					}

				} else {
					finalArtifacts = append(finalArtifacts, e.Artifact)
				}

				if e.IsFinal() {
					fmt.Printf("Final artifact recieved for index %d", e.Artifact.Index)
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
	return fmt.Sprintf("Artifact #%d", artifact.Index+1)
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
	case protocol.TextPart:
		fmt.Println(indent + p.Text)
	case protocol.FilePart:
	case protocol.DataPart:
	default:
		fmt.Printf("%sUnsupported part type: %T\n", indent, p)
	}
}
