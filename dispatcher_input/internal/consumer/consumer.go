package consumer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"stewie.com/dispatcher_input/internal/types"
)

type Consumer struct {
	Topic   string
	Address string
	reader  *kafka.Reader
}

func (c *Consumer) CreateConnection() {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{c.Address},
		Topic:     c.Topic,
		GroupID:   "source-group-id",
		Partition: 0,
	})
}

func (c *Consumer) Read(ctx *context.Context) (*types.Request, error) {
	message, err := c.reader.ReadMessage(*ctx)
	if err != nil {
		return nil, err
	}
	var request types.Request
	err = json.Unmarshal(message.Value, &request)
	if err != nil {
		return nil, err
	}
	return &request, err
}

func (c *Consumer) Close() {
	if err := c.reader.Close(); err != nil {
		logrus.Fatal("failed to close writer:", err)
	}
}
