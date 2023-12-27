package internal

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"stewie.com/dispatcher_input/internal/consumer"
	"stewie.com/dispatcher_input/internal/producer"
)

type RequestHandler struct {
	Ctx      *context.Context
	Consumer *consumer.Consumer
	Producer *producer.Producer
}

func (rh *RequestHandler) Process() {
	for {
		request, err := rh.Consumer.Read(rh.Ctx)
		if err != nil {
			logrus.Fatalln(err)
			continue
		}
		logrus.Infoln(request)
		request.CurrentTime = request.GenerationTime + rh.imitateDelay()
		requestJSON, _ := json.Marshal(*request)
		err = rh.Producer.Write(&kafka.Message{Value: requestJSON})
		if err != nil {
			logrus.Fatal("failed to write messages:", err)
			continue
		}
		logrus.Infof("Request %d sent to buffer \n", request.Id)
	}
}

func (rh *RequestHandler) Close() {
	rh.Consumer.Close()
	rh.Producer.Close()
}

func (rh *RequestHandler) imitateDelay() float32 {
	return 0.1
}
