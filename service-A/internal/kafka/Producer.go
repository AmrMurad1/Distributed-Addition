package kafka

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(defaultHost, topic string) *Producer {
	brokers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if brokers == "" {
		brokers = defaultHost
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}
}

func (p *Producer) SendNumber(number int) error {
	message := kafka.Message{
		Value: []byte(strconv.Itoa(number)),
	}

	err := p.writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Successfully sent number %d to Kafka", number)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
