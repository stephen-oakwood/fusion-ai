package resolver

import (
	"context"
	"fmt"
	"fusion/graph"
	"fusion/graph/model"
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

func (r *Resolver) AgentSendMessage(ctx context.Context, messageInput *model.MessageInput) (<-chan *model.AgentResponse, error) {
	subscriptionChan := make(chan *model.AgentResponse)

	params := createMessageParams(messageInput, 0)

	agentChan, err := r.a2aClient.StreamMessage(ctx, params)
	if err != nil {
		fmt.Printf("Stream Task Request Failed %s", err)
		return nil, err
	}

	go processAgentResponses(ctx, agentChan, subscriptionChan)

	return subscriptionChan, nil
}

func processAgentResponses(ctx context.Context, agentChan <-chan protocol.StreamingMessageEvent, subscriptionChan chan<- *model.AgentResponse) {

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

			switch e := event.Result.(type) {
			case *protocol.TaskStatusUpdateEvent:
				fmt.Printf("  [Status Update: %s (%s)]\n\n", e.Status.State, e.Status.Timestamp)
				if e.Status.Message != nil {

					var parts []model.Part
					for _, part := range e.Status.Message.Parts {
						switch p := part.(type) {
						case *protocol.TextPart:
							parts = append(parts, model.TextPart{Text: p.Text})
						case *protocol.FilePart:
						case *protocol.DataPart:
						default:
							fmt.Printf("%sUnsupported part type: %T\n", p)
						}
					}

					responseType := "TaskStatusUpdateEvent"
					responseState := string(e.Status.State)

					agentMessage := model.AgentResponse{
						MessageID:     &e.Status.Message.MessageID,
						TaskID:        &e.TaskID,
						ContextID:     &e.ContextID,
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

			case *protocol.TaskArtifactUpdateEvent:
				fmt.Printf("  [Artifact Update: %s (Last Chunk %b)]\n\n", *e.Artifact.Name, *e.LastChunk)

				var parts []model.Part
				for _, part := range e.Artifact.Parts {
					switch p := part.(type) {
					case *protocol.TextPart:
						parts = append(parts, model.TextPart{Text: p.Text})
					case *protocol.FilePart:
					case *protocol.DataPart:
					default:
						fmt.Printf("%sUnsupported part type: %T\n", p)
					}
				}

				responseType := "TaskArtifactUpdateEvent"
				agentMessage := model.AgentResponse{
					TaskID:       &e.TaskID,
					ContextID:    &e.ContextID,
					Parts:        parts,
					ResponseType: &responseType,
				}

				subscriptionChan <- &agentMessage

			default:
				fmt.Println("Warning: received unknown event type: %T\n ", event)
			}
		}

	}
}

func createMessageParams(messageInput *model.MessageInput, historyLength int) protocol.SendMessageParams {
	message := protocol.NewMessageWithContext(
		protocol.MessageRoleUser,
		[]protocol.Part{protocol.NewTextPart(messageInput.Text)},
		nil,
		messageInput.ContextID,
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
