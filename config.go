package main

import (
	"encoding/json"
	"io"
	"os"
)

// Config is a container for all configurations of microbial Genetic
// Algorithm (mGA) and DPPN. It is initialized via importing a JSON file.
type Configuration struct {
	// Random Seed
	Seed int64

	// mGA configurations
	NumInputs      int     // number of inputs
	NumOutputs     int     // number of outputs
	NumInitHidden  int     // number of initial hidden nodes
	PopulationSize int     // population size
	NumTournaments int     // number of tournaments
	MutAddNodeRate float64 // mutation rate for adding an node
	MutAddEdgeRate float64 // mutation rate for adding an edge
	CrossoverRate  float64 // crossover rate

	// DPPN configurations
	NumEpochs    int     // number of training epochs
	BatchSize    int     // size of each training batch
	LearningRate float64 // learning rate (alpha)
}

// NewConfiguration creates a new configuration struct given a JSON filename.
func NewConfiguration(filename string) (*Configuration, error) {
	// import configuration
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	var config Configuration
	if err := dec.Decode(&config); err != nil {
		if err != io.EOF {
			return nil, err
		}
	}

	return &config, nil
}
