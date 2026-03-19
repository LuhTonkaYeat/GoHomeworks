package main

import (
	"log"
	"net/http"
	"os"

	"github.com/LuhTonkaYeat/GoHomeworks/hw2/internal/gateway/client"
	"github.com/LuhTonkaYeat/GoHomeworks/hw2/internal/gateway/handler"

	_ "github.com/LuhTonkaYeat/GoHomeworks/hw2/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:50051"
		log.Printf("COLLECTOR_ADDR not set, using default: %s", collectorAddr)
	}

	collectorClient, err := client.NewCollectorClient(collectorAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Collector: %v", err)
	}
	defer collectorClient.Close()

	log.Printf("Connected to Collector at %s", collectorAddr)

	repoHandler := handler.NewRepoHandler(collectorClient)

	http.HandleFunc("/repo", repoHandler.GetRepository)

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	port := ":8080"
	log.Printf("Gateway server is running on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET /repo?owner={owner}&repo={repo} - get repository info")
	log.Printf("  GET /swagger/ - Swagger UI")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
