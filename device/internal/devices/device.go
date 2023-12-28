package devices

import (
	"stewie.com/device/internal"
	"stewie.com/device/internal/types"
)

type Device struct {
	id      int32
	endTime float32
	ud      *internal.UniformDistribution
}

func NewDevice(id int32, ud *internal.UniformDistribution) *Device {
	return &Device{
		id:      id,
		endTime: 0.0,
		ud:      ud,
	}
}

func (d *Device) ProcessRequest(request *types.Request) {
	d.endTime = request.CurrentTime + d.ud.DetermineInterval()
}
