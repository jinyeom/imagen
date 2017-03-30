package mga

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
)

// MGAConfig is a container for all configurations of microbial Genetic
// Algorithm (mGA). It is meant to be imported via a JSON file.
type MGAConfig struct {
	NumInputs      int     // number of inputs
	NumOutputs     int     // number of outputs
	NumInitHidden  int     // number of initial hidden nodes
	PopulationSize int     // population size
	NumTournaments int     // number of tournaments
	MutAddNodeRate float64 // mutation rate for adding an node
	MutAddEdgeRate float64 // mutation rate for adding an edge
	CrossoverRate  float64 // crossover rate
}

func NewMGAConfig(filename string) (*MGAConfig, error) {
	// import configuration
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	var config MGAConfig
	if err := dec.Decode(&config); err != nil {
		if err != io.EOF {
			return nil, err
		}
	}

	return &config, nil
}

// MGA contains an environment of the microbial Genetic Algorithm (mGA).
type MGA struct {
	Config     *MGAConfig // configuration
	Log        *LogBook   // log book
	Population []*Genome  // population of genomes
	Comparison CompFn     // comparison function
	Evaluation EvalFn     // evaluation function
}

// NewMGA creates a new environment for microbial Genetic Algorithm (mGA). It
// returns an error if an invalid configuration file is provided.
func NewMGA(config *MGAConfig, cfn CompFn, efn EvalFn) (*MGA, error) {
	population := make([]*Genome, config.PopulationSize)
	for i := range population {
		population[i] = NewGenome(i, config.NumInputs,
			config.NumInitHidden, config.NumOutputs)
	}

	return &MGA{
		Config:     config,
		Log:        NewLogBook(config.NumTournaments),
		Population: population,
		Comparison: cfn,
		Evaluation: efn,
	}, nil
}

// Run performs microbial Genetic Algorithm (mGA).
func (m *MGA) Run(verbose, exportLog bool) float64 {
	bestScore := 0.0
	if m.Comparison(bestScore, 9999.0) {
		bestScore = 9999.0
	}

	for i := 0; i < m.Config.NumTournaments; i++ {
		ind1 := randGenome(m.Population)
		ind2 := randGenome(m.Population)

		// only evaluate the selected individuals if they were losers
		// in their previous tournament.
		if !ind1.Winner {
			ind1.Fitness = m.Evaluation(ind1)
		}
		if !ind2.Winner {
			ind2.Fitness = m.Evaluation(ind2)
		}

		if m.Comparison(ind1.Fitness, ind2.Fitness) {
			// if score 1 (ind1) is better than score 2 (ind2),
			// perform crossover between the two, and update ind2
			// with the resulting child, and mutate it.

			ind1.Winner = true
			ind2.Winner = false

			if rand.Float64() < m.Config.CrossoverRate {
				ind2.Crossover(ind1)
			}
			ind2.Mutate(m.Config.MutAddNodeRate, m.Config.MutAddEdgeRate)

			if m.Comparison(ind1.Fitness, bestScore) {
				m.Log.Best = ind1
				bestScore = ind1.Fitness
			}
		} else {
			// otherwise, update ind1 (loser) with the resulting
			// child, and mutate it.

			ind1.Winner = false
			ind2.Winner = true

			if rand.Float64() < m.Config.CrossoverRate {
				ind1.Crossover(ind2)
			}
			ind2.Mutate(m.Config.MutAddNodeRate, m.Config.MutAddEdgeRate)

			if m.Comparison(ind2.Fitness, bestScore) {
				m.Log.Best = ind2
				bestScore = ind2.Fitness
			}
		}

		if verbose {
			fmt.Printf("Tournament [%4d] | %3d and %3d | best score: %f\n",
				i, ind1.ID, ind2.ID, bestScore)
		}

		m.Log.Record(ind1.ID, ind2.ID, ind1.Fitness, ind2.Fitness, bestScore)
	}

	if exportLog {
		err := m.Log.Export()
		if err != nil {
			fmt.Println("Log export failed:")
			fmt.Println(err)
		}
	}

	return bestScore
}
