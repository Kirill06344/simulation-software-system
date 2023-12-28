package main

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"stewie.com/source/internal"
	prod "stewie.com/source/internal/producer"
	"stewie.com/source/internal/sources"
	"stewie.com/source/util"
)

func main() {
	cfg, err := util.LoadConfig("config.yml")
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
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}

		log.Infof("Request %d generated on source %d at %f", request.Id, request.SourceId, request.GenerationTime)
	}

}
