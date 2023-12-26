package main

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"stewie.com/source/internal"
	prod "stewie.com/source/internal/producer"
	"stewie.com/source/internal/sources"
	"stewie.com/source/util"
)

func main() {
	cfg, err := util.LoadConfig("../config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	pd := internal.NewPoissonDistribution(cfg.Lambda)
	sourceStorage := sources.NewSourceStorage(cfg.SourcesAmount, cfg.RequestAmount, pd)
	producer := prod.NewProducer("localhost:9092", "source-topic")

	for {
		request, err := sourceStorage.GenerateRequest()
		if err != nil {
			break
		}
		requestJSON, _ := json.Marshal(*request)
		err = producer.Write(&kafka.Message{Value: requestJSON})
		log.Print(request)
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
	}

}
