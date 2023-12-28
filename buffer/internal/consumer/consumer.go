package consumer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Consumer[T comparable] struct {
	Topic   string
	Address string
	GroupID string
	reader  *kafka.Reader
}

func (c *Consumer[T]) CreateConnection() {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{c.Address},
		Topic:     c.Topic,
		GroupID:   c.GroupID,
		Partition: 0,
	})
}

func (c *Consumer[T]) Read(ctx *context.Context) (*T, error) {
	message, err := c.reader.ReadMessage(*ctx)
	if err != nil {
		return nil, err
	}
	var model T
	err = json.Unmarshal(message.Value, &model)
	if err != nil {
		return nil, err
	}
	return &model, err
}

func (c *Consumer[T]) Close() {
	if err := c.reader.Close(); err != nil {
		logrus.Fatal("failed to close writer:", err)
	}
}
