package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/docs"
	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/internal/adapter/grpc"
	httpHandler "github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/internal/adapter/http"
	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/internal/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title GitHub Repository API
// @version 1.0
// @description API for getting information about GitHub repositories
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @schemes http
func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:50051"
		log.Printf("COLLECTOR_ADDR not set, using default: %s", collectorAddr)
	}

	grpcClient, err := grpc.NewClient(collectorAddr)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer grpcClient.Close()

	repoUseCase := usecase.NewRepositoryUseCase(grpcClient)
	httpHandler := httpHandler.NewHandler(repoUseCase)

	mux := http.NewServeMux()

	mux.HandleFunc("/repo", httpHandler.GetRepository)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	port := ":8080"
	log.Printf("Gateway server is running on port %s", port)
	log.Printf("Collector address: %s", collectorAddr)
	log.Printf("Available endpoints:")
	log.Printf("  GET /repo?owner={owner}&repo={repo} - get repository info")
	log.Printf("  GET /swagger/ - Swagger UI")

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
