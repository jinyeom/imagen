/*


dppn.go implementation of DPPN.

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
	"errors"
	"fmt"
	"github.com/gonum/matrix/mat64"
	"math"
	"sort"
)

// Node implements a node in the DPPN (Differentiable Pattern Producing
// Network).
type Node struct {
	ID      int               // node ID
	Type    string            // node type
	AFunc   *ActivationFunc   // activation function
	Inputs  map[*Node]float64 // connected input nodes and weights
	Outputs map[*Node]float64 // connected output nodes and weights
	Signal  *mat64.Vector     // signal in this node
	Delta   *mat64.Vector     // gradient in this node
}

// NewNode decodes an argument node gene (type NodeGene) and creates a new
// node of the DPPN (Differentiable Pattern Producing Network).
func NewNode(n *NodeGene, batchSize int) *Node {
	return &Node{
		ID:      n.ID,
		Type:    n.Type,
		AFunc:   aFuncSet[n.AFuncType],
		Inputs:  make(map[*Node]float64),
		Outputs: make(map[*Node]float64),
		Signal:  mat64.NewVector(batchSize, nil),
		Delta:   mat64.NewVector(batchSize, nil),
	}
}

func (n *Node) ToString() string {
	str := fmt.Sprintf("[%3d] %8s <- {", n.ID, n.AFunc.Name)
	for input, weight := range n.Inputs {
		if weight >= 0.0 {
			str += fmt.Sprintf(" [%3d](%2.4f)", input.ID, weight)
		} else {
			str += fmt.Sprintf(" [%3d](%2.3f)", input.ID, weight)
		}
	}
	str += " }"
	return str
}

// Activate recursively retreives signals from input nodes, activates this
// node, and returns its output vector. It is used for feedforwarding of the
// network.
func (n *Node) Activate() *mat64.Vector {
	if len(n.Inputs) == 0 {
		return n.Signal
	}

	// accumulate input signals
	n.Signal = mat64.NewVector(n.Signal.Len(), nil)
	for node, weight := range n.Inputs {
		n.Signal.AddScaledVec(n.Signal, weight, node.Activate())
	}

	for i := 0; i < n.Signal.Len(); i++ {
		n.Signal.SetVec(i, n.AFunc.Fn(n.Signal.At(i, 0)))
	}
	return n.Signal
}

// SetDelta recursively retrives its delta value from nodes it's connected to.
// Return its delta value. It is used for backpropagation of the network.
func (n *Node) SetDelta() *mat64.Vector {
	if len(n.Outputs) == 0 {
		return n.Delta
	}

	// accumulate errors
	n.Delta = mat64.NewVector(n.Delta.Len(), nil)
	for node, weight := range n.Outputs {
		n.Delta.AddScaledVec(n.Delta, weight, node.SetDelta())
	}

	for i := 0; i < n.Delta.Len(); i++ {
		n.Delta.SetVec(i, n.Delta.At(i, 0)*n.AFunc.DFn(n.Signal.At(i, 0)))
	}
	return n.Delta
}

// UpdateWeights updates each connection weight, given an argument learning
// rate. It must be called after the node's delta value is already retrieved.
func (n *Node) UpdateWeights(learnRate float64) {
	for node, weight := range n.Inputs {
		update := mat64.Dot(n.Delta, node.Signal)
		updatedWeight := weight - learnRate*update
		n.Inputs[node] = updatedWeight
		node.Outputs[n] = updatedWeight
	}
}

// DPPN (Differentiable Pattern Producing Network) implements the phenotype
// of the genotype (Genome).
type DPPN struct {
	ID         int     // genome ID
	NumInputs  int     // number of inputs
	NumOutputs int     // number of outputs
	Nodes      []*Node // nodes in the network
	BatchSize  int     // size of each batch for training
}

// NewDPPN decodes the argument genome and creates a DPPN.
func NewDPPN(g *Genome, batchSize int) (*DPPN, error) {
	nodes, err := func(g *Genome) ([]*Node, error) {
		nodes := make([]*Node, len(g.NodeGenes))
		for i := range nodes {
			nodes[i] = NewNode(g.NodeGenes[i], batchSize)
		}

		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].ID < nodes[j].ID
		})

		for _, edge := range g.EdgeGenes {
			input := sort.Search(len(nodes), func(j int) bool {
				return nodes[j].ID >= edge.InputNode.ID
			})
			output := sort.Search(len(nodes), func(j int) bool {
				return nodes[j].ID >= edge.OutputNode.ID
			})

			// Panic if the genome contains an edge that connect nodes
			// that do not exist.
			if nodes[input].ID != edge.InputNode.ID ||
				nodes[output].ID != edge.OutputNode.ID {
				err := errors.New("Invalid edge found in the genome")
				return nil, err
			}

			// connect from input layer to output layer
			nodes[output].Inputs[nodes[input]] = edge.Weight
			nodes[input].Outputs[nodes[output]] = edge.Weight
		}
		return nodes, nil
	}(g)

	if err != nil {
		return nil, err
	}

	return &DPPN{
		ID:         g.ID,
		NumInputs:  g.NumInputs,
		NumOutputs: g.NumOutputs,
		Nodes:      nodes,
		BatchSize:  batchSize,
	}, nil
}

// ToString summarizes the network's connections among its nodes.
func (d *DPPN) ToString() string {
	str := fmt.Sprintf("DPPN(%d)\n", d.ID)
	for i := 0; i < len(d.Nodes)-1; i++ {
		str += d.Nodes[i].ToString() + "\n"
	}
	node := d.Nodes[len(d.Nodes)-1]
	str += node.ToString()

	return str
}

// FeedForward feeds a slice of inputs through the network and returns a slice
// of estimated outputs. It returns an error if the input slice has an invalid
// length.
func (d *DPPN) FeedForward(inputs *mat64.Dense) (*mat64.Dense, error) {
	if _, c := inputs.Dims(); c != d.NumInputs {
		return nil, errors.New("Invalid number of inputs")
	}

	// send input signals to input nodes
	for i := 0; i < d.NumInputs; i++ {
		inputVector := mat64.Col(make([]float64, d.BatchSize), i, inputs)
		d.Nodes[i].Signal = mat64.NewVector(d.BatchSize, inputVector)
	}

	outputs := mat64.NewDense(d.BatchSize, d.NumOutputs, nil)
	for i := 0; i < d.NumOutputs; i++ {
		outputVec := d.Nodes[i+d.NumInputs].Activate()
		outputs.SetCol(i, outputVec.RawVector().Data)
	}

	return outputs, nil
}

// Backprop updates the network's weights via Backpropagation, given a slice
// of inputs, a slice of target outputs, and the learning rate. It returns
// a slice of estimated output, its MSE (mean squared error). Return error if
// the input slice has an invalid length.
func (d *DPPN) Backprop(inputs, target *mat64.Dense,
	learningRate float64) (float64, error) {
	if _, c := target.Dims(); c != d.NumOutputs {
		return 0.0, errors.New("Invalid number of outputs")
	}

	// feedforward input vector signal with a side effect of storing each
	// node's signal.
	outputs, err := d.FeedForward(inputs)
	if err != nil {
		return 0.0, err
	}

	// compute mean squared error
	var outputErr mat64.Dense
	outputErr.Sub(outputs, target)

	r, c := outputErr.Dims()
	mse := math.Pow(mat64.Sum(&outputErr)/float64(r*c), 2.0)

	// compute delta vector and assign them
	for i := 0; i < d.NumOutputs; i++ {
		node := d.Nodes[i+d.NumInputs]
		errSlice := mat64.Col(make([]float64, d.BatchSize), i, &outputErr)
		errVec := mat64.NewVector(d.BatchSize, errSlice)

		derivSignal := mat64.NewVector(d.BatchSize, nil)
		for j := 0; j < d.BatchSize; j++ {
			derivSignal.SetVec(j, node.AFunc.DFn(node.Signal.At(j, 0)))
		}
		node.Delta.MulElemVec(errVec, derivSignal)
	}

	// recursively compute error in each node from input nodes to output nodes
	for i := 0; i < d.NumInputs; i++ {
		d.Nodes[i].SetDelta()
	}

	// update all the weights
	for _, node := range d.Nodes {
		node.UpdateWeights(learningRate)
	}

	return mse, nil
}

// Encode encodes the DPPN's weights to the argument genome. Return error if
// the argument genome does not have the same ID as the DPPN.
func (d *DPPN) Encode(g *Genome) error {
	if d.ID != g.ID {
		return fmt.Errorf("Invalid Genome: Must encode into a "+
			"genome with the same ID (%d != %d)", d.ID, g.ID)
	}

	for _, node := range d.Nodes {
		for input, weight := range node.Inputs {
			for i, edge := range g.EdgeGenes {
				// update the argument genome's edge weights
				if edge.InputNode.ID == input.ID &&
					edge.OutputNode.ID == node.ID {
					g.EdgeGenes[i].Weight = weight
				}
			}
		}
	}

	return nil
}
