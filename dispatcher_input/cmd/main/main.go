package main

import (
	"context"
	"stewie.com/dispatcher_input/internal"
	consumer "stewie.com/dispatcher_input/internal/consumer"
	producer "stewie.com/dispatcher_input/internal/producer"
	"time"
)

func main() {
	c := &consumer.Consumer{
		Topic:   "source-topic",
		Address: "localhost:9092",
	}
	c.CreateConnection()

	p := producer.NewProducer("localhost:9092", "buffer-topic")
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	requestHandler := internal.RequestHandler{
		Ctx:      &ctx,
		Consumer: c,
		Producer: p,
	}

	requestHandler.Process()
	requestHandler.Close()
}
