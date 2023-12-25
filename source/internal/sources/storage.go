package sources

import (
	"errors"
	"stewie.com/source/internal"
)

type Storage struct {
	sources       []Source
	requestAmount int32
}

func NewSourceStorage(sourcesAmount int32, requestAmount int32, pd *internal.PoissonDistribution) *Storage {
	sources := make([]Source, sourcesAmount)
	for i := range sources {
		sources[i] = *NewSource(int32(i), pd)
	}
	return &Storage{sources: sources, requestAmount: requestAmount}
}

func (storage *Storage) GenerateRequest() (string, error) {
	if storage.requestAmount > 0 {
		storage.requestAmount--
		return storage.getYoungestSource().generateRequest(), nil
	}
	return "", errors.New("storage is empty")
}

func (storage *Storage) getYoungestSource() *Source {
	minGenTime := storage.sources[0].nextGenerationTime
	inxOfYoungest := 0
	for inx, value := range storage.sources {
		if value.nextGenerationTime < minGenTime {
			minGenTime = value.nextGenerationTime
			inxOfYoungest = inx
		}
	}
	return &storage.sources[inxOfYoungest]
}
