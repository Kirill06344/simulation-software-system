package internal

import "math/rand"

type UniformDistribution struct {
	alpha int32
	beta  int32
}

func NewUniformDistribution(alpha, beta int32) *UniformDistribution {
	return &UniformDistribution{
		alpha: alpha,
		beta:  beta,
	}
}

func (ud *UniformDistribution) DetermineInterval() float32 {
	return float32(float64(ud.alpha+(ud.beta-ud.alpha)) * rand.Float64())
}
