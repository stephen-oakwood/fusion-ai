package main

import (
	"fusion/internal/a2a"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"trpc.group/trpc-go/trpc-a2a-go/server"
	redismgr "trpc.group/trpc-go/trpc-a2a-go/taskmanager/redis"
)

func main() {

	agentCard := a2a.GetAgentCard()
	processor, err := a2a.NewAgent()
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	redisOptions := redis.UniversalOptions{
		Addrs:    []string{"localhost:6379"},
		Password: "",
		DB:       0,
	}

	taskManager, err := redismgr.NewRedisTaskManager(redis.NewUniversalClient(&redisOptions), processor)
	if err != nil {
		log.Fatalf("Failed to create task manager: %v", err)
	}

	options := []server.Option{
		server.WithIdleTimeout(120 * time.Second),
		server.WithReadTimeout(60 * time.Second),
		server.WithWriteTimeout(60 * time.Second),
	}
	srv, err := server.NewA2AServer(agentCard, taskManager, options...)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start("localhost:8080"); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	sig := <-sigChan
	log.Printf("Received signal %v, shutting down...", sig)
}

func stringPtr(s string) *string {
	return &s
}
