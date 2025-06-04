package main

import (
	"fmt"
	"fusion/graph"
	"fusion/internal/resolver"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"trpc.group/trpc-go/trpc-a2a-go/client"
)

func main() {

	a2aClient, err := client.NewA2AClient("http://localhost:8080", client.WithTimeout(300*time.Second))
	if err != nil {
		fmt.Errorf("failed to create A2A client: %v", err)
	}

	agentResolver, err := resolver.NewResolver(a2aClient)
	if err != nil {
		panic(err)
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: agentResolver}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 300 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Printf("connect to http://localhost:8180/ for GraphQL playground")
	http.ListenAndServe(":8180", nil)
}
