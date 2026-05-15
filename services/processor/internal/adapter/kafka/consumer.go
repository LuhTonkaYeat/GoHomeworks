package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type RepositoryFetchResponse struct {
	RequestID   string `json:"request_id"`
	Owner       string `json:"owner"`
	Repo        string `json:"repo"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	CreatedAt   string `json:"created_at"`
	Error       string `json:"error,omitempty"`
}

type ResponseHandler func(RepositoryFetchResponse) error

type Consumer struct {
	reader  *kafka.Reader
	handler ResponseHandler
}

func NewConsumer(brokers []string, topic string, groupID string, handler ResponseHandler) *Consumer {
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

				var response RepositoryFetchResponse
				if err := json.Unmarshal(msg.Value, &response); err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					continue
				}

				log.Printf("Received response for %s/%s (request_id: %s)", response.Owner, response.Repo, response.RequestID)

				if err := c.handler(response); err != nil {
					log.Printf("Error handling response: %v", err)
				}
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}