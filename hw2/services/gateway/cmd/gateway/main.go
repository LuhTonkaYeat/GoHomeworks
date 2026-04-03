package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/docs"
	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/internal/adapter/grpc"
	httpHandler "github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/internal/delivery/http"
	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/internal/usecase"
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

	grpcClient, err := grpc.NewClient(processorAddr)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer grpcClient.Close()

	repoUseCase := usecase.NewRepositoryUseCase(grpcClient)
	httpHandler := httpHandler.NewHandler(repoUseCase)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/repositories/info", httpHandler.GetRepository)
	mux.HandleFunc("/api/ping", httpHandler.Ping)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	port := ":8080"
	log.Printf("Gateway server started on port %s", port)
	log.Printf("Processor address: %s", processorAddr)
	log.Printf("Available endpoints:")
	log.Printf("  GET /api/repositories/info?url=<github_url> - get repository info")
	log.Printf("  GET /api/ping - check services status")
	log.Printf("  GET /swagger/ - Swagger UI")

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
