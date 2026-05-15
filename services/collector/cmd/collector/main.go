package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/adapter/github"
	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/adapter/kafka"
	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/adapter/subscribe"
	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/usecase"
)

func main() {
	subscribeAddr := os.Getenv("SUBSCRIBE_ADDR")
	if subscribeAddr == "" {
		subscribeAddr = "subscribe:50053"
	}
	
	subscribeClient, err := subscribe.NewClient(subscribeAddr)
	if err != nil {
		log.Fatalf("Failed to create Subscribe client: %v", err)
	}
	defer subscribeClient.Close()
	
	githubClient := github.NewClient()
	
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"kafka:9092"}
	}
	
	requestsTopic := os.Getenv("KAFKA_REQUESTS_TOPIC")
	if requestsTopic == "" {
		requestsTopic = "repository.fetch.requests"
	}
	
	responsesTopic := os.Getenv("KAFKA_RESPONSES_TOPIC")
	if responsesTopic == "" {
		responsesTopic = "repository.fetch.responses"
	}
	
	kafkaProducer := kafka.NewProducer(kafkaBrokers, responsesTopic)
	defer kafkaProducer.Close()
	
	collectorUseCase := usecase.NewCollectorUseCase(githubClient, subscribeClient, kafkaProducer)
	
	kafkaConsumer := kafka.NewConsumer(kafkaBrokers, requestsTopic, "collector-group", func(request kafka.RepositoryFetchRequest) error {
		log.Printf("Received request: %+v", request)
		return collectorUseCase.ProcessFetchRequest(context.Background(), request)
	})
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	kafkaConsumer.Start(ctx)
	defer kafkaConsumer.Close()
	
	log.Println("Collector service started")
	log.Printf("Subscribe address: %s", subscribeAddr)
	log.Printf("Kafka brokers: %v", kafkaBrokers)
	log.Printf("Requests topic: %s", requestsTopic)
	log.Printf("Responses topic: %s", responsesTopic)
	log.Println("Waiting for Kafka messages...")
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	
	log.Println("Shutting down...")
	cancel()
}