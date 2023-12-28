package devices

import (
	"container/heap"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"stewie.com/device/internal"
	"stewie.com/device/internal/producer"
	"stewie.com/device/internal/types"
)

type DeviceHeap []Device

func (h DeviceHeap) Len() int {
	return len(h)
}

func (h DeviceHeap) Less(i, j int) bool {
	if h[i].endTime == h[j].endTime {
		return h[i].id < h[j].id
	}
	return h[i].endTime < h[j].endTime
}

func (h DeviceHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *DeviceHeap) Push(x interface{}) {
	*h = append(*h, x.(Device))
}

func (h *DeviceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Storage struct {
	devices      DeviceHeap
	size         int32
	distribution *internal.UniformDistribution
	producer     *producer.Producer
}

func NewDeviceStorage(deviceAmount int32, pd *internal.UniformDistribution, prod *producer.Producer) *Storage {
	return &Storage{
		devices:      DeviceHeap{},
		size:         deviceAmount,
		distribution: pd,
		producer:     prod,
	}
}

func (ds *Storage) Start() {
	for i := 1; int32(i) <= ds.size; i++ {
		var notification = types.Notification{
			DeviceId:      int32(i),
			ReleasingTime: 0.0,
		}
		notificationJSON, err := json.Marshal(&notification)
		if err != nil {
			log.Fatal("unable to marshal")
		}
		ds.producer.Write(&kafka.Message{Value: notificationJSON})
	}
}

func (ds *Storage) ProcessRequest(requests []types.Request) {
	for _, request := range requests {
		device := NewDevice(request.DeviceId, ds.distribution)
		device.ProcessRequest(&request)
		heap.Push(&ds.devices, *device)
	}
	for i := 0; int32(i) < ds.size; i++ {
		youngestDevice := heap.Pop(&ds.devices).(Device)
		var notification = types.Notification{
			DeviceId:      youngestDevice.id,
			ReleasingTime: youngestDevice.endTime,
		}
		notificationJSON, err := json.Marshal(&notification)
		if err != nil {
			log.Fatal("unable to marshal")
		}
		ds.producer.Write(&kafka.Message{Value: notificationJSON})
	}
}
