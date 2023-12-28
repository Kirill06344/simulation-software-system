package producer

import (
	"context"
	"errors"
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
		BatchTimeout:           500 * time.Millisecond,
	}

	return &Producer{writer: w}
}

func (p *Producer) Write(message *kafka.Message) error {
	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = p.writer.WriteMessages(ctx, *message)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
