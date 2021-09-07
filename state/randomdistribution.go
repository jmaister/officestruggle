package state

import "math/rand"

type DistributedRandom interface {
	SetNumber(value int, distribution float64)
	GetDistributedRandom() int
}

type distribution struct {
	sum           float64
	distributions map[int]float64
}

func (d *distribution) SetNumber(value int, distribution float64) {
	dist, ok := d.distributions[value]
	if ok {
		d.sum -= dist
	}
	d.distributions[value] = distribution
	d.sum += distribution
}

func (d *distribution) GetDistributedRandom() int {
	rndNum := rand.Float64()
	ratio := 1.0 / d.sum
	temp := 0.0
	for key, distribution := range d.distributions {
		temp += distribution
		if rndNum/ratio <= temp {
			return key
		}
	}
	return 0
}

func NewDistributedRandom() DistributedRandom {
	return &distribution{
		sum:           0,
		distributions: map[int]float64{},
	}
}
