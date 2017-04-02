/*


mga.go implementation of mGA.

@licstart   The following is the entire license notice for
the Go code in this page.

Copyright (C) 2017 jin yeom

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

As additional permission under GNU GPL version 3 section 7, you
may distribute non-source (e.g., minimized or compacted) forms of
that code without the copy of the GNU GPL normally required by
section 4, provided you include this license notice and a URL
through which recipients can access the Corresponding Source.

@licend    The above is the entire license notice
for the Go code in this page.


*/

package main

import (
	"fmt"
	"math/rand"
)

// MGA contains an environment of the microbial Genetic Algorithm (mGA).
type MGA struct {
	Config     *Configuration // configuration
	Log        *LogBook       // log book
	Population []*Genome      // population of genomes
	Comparison ComparisonFunc // comparison function
	Evaluation EvaluationFunc // evaluation function
}

// NewMGA creates a new environment for microbial Genetic Algorithm (mGA). It
// returns an error if an invalid configuration file is provided.
func NewMGA(config *Configuration, comparison ComparisonFunc,
	evaluation EvaluationFunc) (*MGA, error) {
	population := make([]*Genome, config.PopulationSize)
	for i := range population {
		population[i] = NewGenome(i, config.NumInputs,
			config.NumInitHidden, config.NumOutputs)
	}

	return &MGA{
		Config:     config,
		Log:        NewLogBook(config.NumTournaments),
		Population: population,
		Comparison: comparison,
		Evaluation: evaluation,
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

		ind1.Fitness = m.Evaluation(ind1)
		ind2.Fitness = m.Evaluation(ind2)

		if m.Comparison(ind1.Fitness, ind2.Fitness) {
			// if score 1 (ind1) is better than score 2 (ind2),
			// perform crossover between the two, and update ind2
			// with the resulting child, and mutate it.
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
		if err := m.Log.Export(); err != nil {
			fmt.Println("Log export failed:")
			fmt.Println(err)
		}
	}

	return bestScore
}
