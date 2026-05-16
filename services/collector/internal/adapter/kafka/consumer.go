package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type RepositoryFetchRequest struct {
	RequestID string `json:"request_id"`
	Owner     string `json:"owner"`
	Repo      string `json:"repo"`
}

type RequestHandler func(RepositoryFetchRequest) error

type Consumer struct {
	reader  *kafka.Reader
	handler RequestHandler
}

func NewConsumer(brokers []string, topic string, groupID string, handler RequestHandler) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping Kafka consumer...")
				return
			default:
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					log.Printf("Error reading message: %v", err)
					continue
				}

				var request RepositoryFetchRequest
				if err := json.Unmarshal(msg.Value, &request); err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					continue
				}

				log.Printf("Received fetch request for %s/%s (request_id: %s)", request.Owner, request.Repo, request.RequestID)

				if err := c.handler(request); err != nil {
					log.Printf("Error handling request: %v", err)
				}
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}