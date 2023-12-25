package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"stewie.com/source/internal"
	"stewie.com/source/internal/sources"
	"stewie.com/source/util"
	"time"
)

func main() {
	cfg, err := util.LoadConfig("../config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	pd := internal.NewPoissonDistribution(cfg.Lambda)
	sourceStorage := sources.NewSourceStorage(cfg.SourcesAmount, cfg.RequestAmount, pd)

	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
		return
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	for {
		request, err := sourceStorage.GenerateRequest()
		if err != nil {
			break
		}
		requestJSON, _ := json.Marshal(*request)
		_, err = conn.WriteMessages(kafka.Message{Value: requestJSON})
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}

	}
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}
