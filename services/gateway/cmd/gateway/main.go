package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/LuhTonkaYeat/GoHomeworks/services/gateway/docs"
	"github.com/LuhTonkaYeat/GoHomeworks/services/gateway/internal/adapter/grpc"
	httpHandler "github.com/LuhTonkaYeat/GoHomeworks/services/gateway/internal/delivery/http"
	"github.com/LuhTonkaYeat/GoHomeworks/services/gateway/internal/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title GitHub Repository API
// @version 2.0
// @description API for getting information about GitHub repositories
// @host localhost:8080
// @BasePath /

func main() {
	processorAddr := os.Getenv("PROCESSOR_ADDR")
	if processorAddr == "" {
		processorAddr = "processor:50053"
		log.Printf("PROCESSOR_ADDR not set, using default: %s", processorAddr)
	}

	subscribeAddr := os.Getenv("SUBSCRIBE_ADDR")
	if subscribeAddr == "" {
		subscribeAddr = "subscribe:50053"
		log.Printf("SUBSCRIBE_ADDR not set, using default: %s", subscribeAddr)
	}

	processorClient, err := grpc.NewClient(processorAddr)
	if err != nil {
		log.Fatalf("Failed to create gRPC client for Processor: %v", err)
	}
	defer processorClient.Close()

	subscribeClient, err := grpc.NewSubscribeClient(subscribeAddr)
	if err != nil {
		log.Fatalf("Failed to create gRPC client for Subscribe: %v", err)
	}
	defer subscribeClient.Close()

	repoUseCase := usecase.NewRepositoryUseCase(processorClient)
	subscriptionUseCase := usecase.NewSubscriptionUseCase(subscribeClient, processorClient)

	httpHandler := httpHandler.NewHandler(repoUseCase, subscriptionUseCase)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/repositories/info", httpHandler.GetRepository)
	mux.HandleFunc("/api/ping", httpHandler.Ping)

	mux.HandleFunc("POST /subscriptions", httpHandler.CreateSubscription)
	mux.HandleFunc("DELETE /subscriptions/{owner}/{repo}", httpHandler.DeleteSubscription)
	mux.HandleFunc("GET /subscriptions", httpHandler.GetSubscriptions)
	mux.HandleFunc("GET /subscriptions/info", httpHandler.GetSubscriptionsInfo)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	port := ":8080"
	log.Printf("Gateway server started on port %s", port)
	log.Printf("Processor address: %s", processorAddr)
	log.Printf("Subscribe address: %s", subscribeAddr)
	log.Printf("Available endpoints:")
	log.Printf("  GET /api/repositories/info?url=<github_url> - get repository info")
	log.Printf("  GET /api/ping - check services status")
	log.Printf("  POST /subscriptions - subscribe to repository")
	log.Printf("  DELETE /subscriptions/{owner}/{repo} - unsubscribe")
	log.Printf("  GET /subscriptions - get all subscriptions")
	log.Printf("  GET /subscriptions/info - get info for all subscriptions")
	log.Printf("  GET /swagger/ - Swagger UI")

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
