package sources

import (
	"fmt"
	"stewie.com/source/internal"
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

func (source *Source) generateRequest() string {
	request := fmt.Sprintf("sourceId: %d, time: %f", source.id, source.nextGenerationTime)
	source.nextGenerationTime += source.distribution.DetermineInterval()
	return request
}
