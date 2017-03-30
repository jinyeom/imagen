package mga

import (
	"math/rand"
)

// randWeight returns a random connection weight.
func randWeight() float64 {
	return rand.NormFloat64()
}

// randAFuncName returns a random activation function name.
func randAFuncName() string {
	options := []string{
		"identity",
		"sigmoid",
		"tanh",
		"relu",
		"sine",
		"gaussian",
	}
	return options[rand.Intn(len(options))]
}

func randGenome(population []*Genome) *Genome {
	return population[rand.Intn(len(population))]
}
