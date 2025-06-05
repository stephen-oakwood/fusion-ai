package resolver

import (
	"context"
	"fmt"
	"fusion/graph"
	"fusion/graph/model"
	"reflect"
	"trpc.group/trpc-go/trpc-a2a-go/client"
	"trpc.group/trpc-go/trpc-a2a-go/protocol"
)

type Resolver struct {
	graph.ResolverRoot
	a2aClient *client.A2AClient
}

type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Subscription() graph.SubscriptionResolver {
	return &subscriptionResolver{r}
}

func NewResolver(a2aClient *client.A2AClient) (*Resolver, error) {
	return &Resolver{a2aClient: a2aClient}, nil
}

func (r *Resolver) Placeholder(ctx context.Context) (*string, error) {
	str := "Hello World"
	return &str, nil
}

func (r *Resolver) AgentTaskExecute(ctx context.Context, task *model.TaskInput) (<-chan *model.AgentResponse, error) {
	subscriptionChan := make(chan *model.AgentResponse)

	params := createTaskParams(task.ID, task.SessionID, task.Message)

	agentChan, err := r.a2aClient.StreamTask(ctx, params)
	if err != nil {
		fmt.Printf("Stream Task Request Failed %s", err)
		return nil, err
	}

	go processAgentResponses(ctx, agentChan, subscriptionChan)

	return subscriptionChan, nil
}

func processAgentResponses(ctx context.Context, agentChan <-chan protocol.TaskEvent, subscriptionChan chan<- *model.AgentResponse) {

	defer close(subscriptionChan)

	for {

		select {
		case <-ctx.Done():
			fmt.Printf("Context timeout or cancellatiom while waiting for events: %s", ctx.Err())
			return

		case event, ok := <-agentChan:
			if !ok {
				fmt.Println("\nStream channel closed")
				if ctx.Err() != nil {
					fmt.Printf("Context error after stream close: %s", ctx.Err())
				}
				return
			}

			switch e := event.(type) {
			case protocol.TaskStatusUpdateEvent:
				fmt.Printf("  [Status Update: %s (%s)]\n\n", e.Status.State, e.Status.Timestamp)
				if e.Status.Message != nil {

					var parts []model.Part
					for _, part := range e.Status.Message.Parts {
						switch p := part.(type) {
						case protocol.TextPart:
							parts = append(parts, model.TextPart{Text: p.Text})
						case protocol.FilePart:
							parts = append(parts, model.FilePart{Name: *p.File.Name, MimeType: *p.File.MimeType, Bytes: p.File.Bytes, URI: p.File.URI})
						case protocol.DataPart:
						default:
							fmt.Printf("%sUnsupported part type: %T\n", p)
						}
					}

					responseType := reflect.TypeOf(e).Name()
					responseState := string(e.Status.State)

					agentMessage := model.AgentResponse{
						Parts:         parts,
						ResponseType:  &responseType,
						ResponseState: &responseState,
					}

					subscriptionChan <- &agentMessage
				}

				if e.IsFinal() {
					if e.Status.State == protocol.TaskStateCompleted ||
						e.Status.State == protocol.TaskStateFailed ||
						e.Status.State == protocol.TaskStateCanceled {
						return
					}
				}

			case protocol.TaskArtifactUpdateEvent:
				fmt.Printf("  [Artifact Update: %s (Last Chunk %b)]\n\n", *e.Artifact.Name, *e.Artifact.LastChunk)

				var parts []model.Part
				for _, part := range e.Artifact.Parts {
					switch p := part.(type) {
					case protocol.TextPart:
						parts = append(parts, model.TextPart{Text: p.Text})
					case protocol.FilePart:
						parts = append(parts, model.FilePart{Name: *p.File.Name, MimeType: *p.File.MimeType, Bytes: p.File.Bytes, URI: p.File.URI})
					case protocol.DataPart:
					default:
						fmt.Printf("%sUnsupported part type: %T\n", p)
					}
				}

				messageType := reflect.TypeOf(e).Name()
				agentMessage := model.AgentResponse{
					Parts:        parts,
					ResponseType: &messageType,
				}

				subscriptionChan <- &agentMessage

			default:
				fmt.Println("Warning: received unknown event type: %T\n ", event)
			}
		}

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
