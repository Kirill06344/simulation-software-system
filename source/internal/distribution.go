package internal

import (
	"math"
	"math/rand"
)

type PoissonDistribution struct {
	lambda int32
}

func NewPoissonDistribution(lambda int32) *PoissonDistribution {
	return &PoissonDistribution{lambda: lambda}
}

func (p *PoissonDistribution) DetermineInterval() float32 {
	return -1.0 / float32(p.lambda) * float32(math.Log(rand.Float64()))
}
