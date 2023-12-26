package producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(address, topic string) *Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(address),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		BatchSize:              1,
		BatchTimeout:           10 * time.Millisecond,
	}

	return &Producer{writer: w}
}

func (p *Producer) Write(message *kafka.Message) error {
	err := p.writer.WriteMessages(context.Background(), *message)
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
