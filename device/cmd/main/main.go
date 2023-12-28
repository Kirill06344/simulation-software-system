package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"stewie.com/device/internal"
	"stewie.com/device/internal/consumer"
	"stewie.com/device/internal/devices"
	"stewie.com/device/internal/producer"
	"stewie.com/device/internal/types"
	"stewie.com/device/internal/util"
)

func main() {
	cfg, err := util.LoadConfig("../config.yml")
	if err != nil {
		logrus.Fatal("unable to read config")
	}

	p := producer.NewProducer("localhost:9092", "buffer-notification-topic")
	distribution := internal.NewUniformDistribution(cfg.Alpha, cfg.Beta)

	deviceStorage := devices.NewDeviceStorage(cfg.DevicesAmount, distribution, p)
	deviceStorage.Start()

	requestConsumer := consumer.Consumer{
		Topic:   "device-topic",
		Address: "localhost:9092",
	}
	requestConsumer.CreateConnection()

	var receivedReplies int32 = 0
	requests := make([]types.Request, cfg.DevicesAmount)
	ctx := context.Background()
	for {
		request, _ := requestConsumer.Read(&ctx)
		requests[receivedReplies] = *request
		receivedReplies++
		if receivedReplies == cfg.DevicesAmount {
			receivedReplies = 0
			deviceStorage.ProcessRequest(requests)
		}
	}

}
