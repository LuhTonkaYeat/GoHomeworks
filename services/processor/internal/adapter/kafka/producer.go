package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type RepositoryFetchRequest struct {
	RequestID string `json:"request_id"`
	Owner     string `json:"owner"`
	Repo      string `json:"repo"`
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

func (p *Producer) SendFetchRequest(ctx context.Context, requestID, owner, repo string) error {
	msg := RepositoryFetchRequest{
		RequestID: requestID,
		Owner:     owner,
		Repo:      repo,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(requestID),
		Value: data,
	})

	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Sent fetch request for %s/%s (request_id: %s)", owner, repo, requestID)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}