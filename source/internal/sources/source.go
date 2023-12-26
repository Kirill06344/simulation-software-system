package sources

import (
	"github.com/google/uuid"
	"stewie.com/source/internal"
	"stewie.com/source/internal/types"
)

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
		Id:          uuid.New().String(),
		CurrentTime: source.nextGenerationTime,
		SourceId:    source.id,
	}
	source.nextGenerationTime += source.distribution.DetermineInterval()
	return request
}
