package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"math"
	"stewie.com/buffer/internal"
	"stewie.com/buffer/internal/consumer"
	"stewie.com/buffer/internal/producer"
	"stewie.com/buffer/internal/types"
	"stewie.com/buffer/util"
	"time"
)

func main() {

	cfg, err := util.LoadConfig("../config.yml")

	if err != nil {
		logrus.Fatal("unable to read config")
		return
	}

	requestConsumer := &consumer.Consumer[types.Request]{
		Topic:   "buffer-topic",
		Address: "localhost:9092",
		GroupID: "buffer-group-id",
	}
	requestConsumer.CreateConnection()

	notificationConsumer := &consumer.Consumer[types.Notification]{
		Topic:   "buffer-notification-topic",
		Address: "localhost:9092",
		GroupID: "notification-group-id",
	}
	notificationConsumer.CreateConnection()

	requestProducer := producer.NewProducer("localhost:9092", "device-topic")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	buffer := internal.NewBuffer(cfg.BufferCapacity)
	fmt.Println(buffer)

	notification, err := notificationConsumer.Read(&ctx)
	if err != nil {
		logrus.Warnf("No need to send notifications")
	}
	for {
		request, err := requestConsumer.Read(&ctx)
		if err != nil {
			logrus.Fatal("No more input requests")
		}

		for !buffer.IsEmpty() && notification.ReleasingTime < request.CurrentTime {
			bufferRequest := buffer.Poll()
			bufferRequest.DeviceId = notification.DeviceId
			bufferRequest.CurrentTime = float32(math.Max(float64(bufferRequest.CurrentTime), float64(notification.ReleasingTime)))
			requestJSON, _ := json.Marshal(bufferRequest)
			requestProducer.Write(&kafka.Message{Value: requestJSON})
			logrus.Infof("Device %d took for processing request %d at %f", bufferRequest.DeviceId, bufferRequest.Id, bufferRequest.CurrentTime)
			notification, err = notificationConsumer.Read(&ctx)
			if err != nil {
				break
			}
		}

		if request != nil {
			if buffer.IsFilled() {
				refusedRequest := buffer.RefuseRequest(request)
				logrus.Warnf("Request %d was refused due to request %d at %f", refusedRequest.Id, request.Id, request.CurrentTime)
				continue
			}
			inx := buffer.AddRequest(request)
			logrus.Infof("Request with %d id was added into buffer on index %d at %f", request.Id, inx, request.CurrentTime)
		}
	}

}
