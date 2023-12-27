package sources

import (
	"stewie.com/source/internal"
	"stewie.com/source/internal/types"
)

var requestId uint32 = 1

type Source struct {
	id                 int32
	nextGenerationTime float32
	distribution       *internal.PoissonDistribution
}

func NewSource(id int32, pd *internal.PoissonDistribution) *Source {
	return &Source{
		id:                 id,
		nextGenerationTime: pd.DetermineInterval(),
		distribution:       pd,
	}
}

func (source *Source) generateRequest() types.Request {
	request := types.Request{
		Id:             requestId,
		GenerationTime: source.nextGenerationTime,
		SourceId:       source.id,
	}
	requestId++
	source.nextGenerationTime += source.distribution.DetermineInterval()
	return request
}
