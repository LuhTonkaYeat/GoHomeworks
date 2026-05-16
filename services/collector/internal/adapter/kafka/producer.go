package kafka

import (
	"context"
	"encoding/json"
	"log"

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

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		writer: writer,
		topic:  topic,
	}
}

func (p *Producer) SendFetchResponse(ctx context.Context, response RepositoryFetchResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(response.RequestID),
		Value: data,
	})

	if err != nil {
		log.Printf("Failed to send response to Kafka: %v", err)
		return err
	}

	log.Printf("Sent fetch response for %s/%s (request_id: %s)", response.Owner, response.Repo, response.RequestID)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}