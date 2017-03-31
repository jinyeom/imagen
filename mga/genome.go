/*


genome.go implementation of the genome in mGA.

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

package mga

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

// NodeGene represents a node in the graph representation of a genome.
type NodeGene struct {
	ID        int    // node ID
	Type      string // node type
	AFuncType string // name of activation function
}

// NewNodeGene creates a new node gene, given a node ID, a node type, and its
// activation function type.
func NewNodeGene(id int, nodeType, afuncType string) *NodeGene {
	return &NodeGene{
		ID:        id,
		Type:      nodeType,
		AFuncType: afuncType,
	}
}

// EdgeGene represents an edge in the graph representation of a genome.
type EdgeGene struct {
	InputNode  *NodeGene // input node
	OutputNode *NodeGene // output node
	Weight     float64   // connection weight
	Disabled   bool      // true if disabled
}

// NewEdgeGene creates a new edge gene, given pointers to input node gene and
// output node gene.
func NewEdgeGene(inputNode, outputNode *NodeGene) *EdgeGene {
	return &EdgeGene{
		InputNode:  inputNode,
		OutputNode: outputNode,
		Weight:     randWeight(),
		Disabled:   false,
	}
}

// Genome is a graph G = {N, E}, where N is a list of node genes and E is a
// list of edge genes.
type Genome struct {
	ID         int         // genome ID
	NumInputs  int         // number of inputs
	NumOutputs int         // number of outputs
	NumHidden  int         // number of hidden nodes
	NodeGenes  []*NodeGene // list of node genes
	EdgeGenes  []*EdgeGene // list of edge genes
	Fitness    float64     // fitness score
	Winner     bool        // won the most recent tournament
}

// NewGenome creates a new genome, given a genome ID, number of inputs, number
// of initial hidden nodes, and number of outputs.
func NewGenome(id, numInputs, numHidden, numOutputs int) *Genome {
	// initialize fully connected input and output nodes
	nodes := make([]*NodeGene, 0, numInputs+numHidden+numOutputs)
	edges := make([]*EdgeGene, 0, numInputs*numHidden+numHidden*numOutputs)

	// input nodes
	for i := 0; i < numInputs; i++ {
		nodes = append(nodes, NewNodeGene(i, "input", "identity"))
	}

	// output nodes
	for i := numInputs; i < numInputs+numOutputs; i++ {
		output := NewNodeGene(i, "output", "sigmoid")
		nodes = append(nodes, output)
	}

	// hidden nodes
	iter := numInputs + numOutputs
	for i := iter; i < iter+numHidden; i++ {
		hidden := NewNodeGene(i, "hidden", randAFuncName())
		nodes = append(nodes, hidden)
		// connect to all input nodes
		for j := 0; j < numInputs; j++ {
			edges = append(edges, NewEdgeGene(nodes[j], hidden))
		}

		// connect to all output nodes
		for j := numInputs; j < numInputs+numOutputs; j++ {
			edges = append(edges, NewEdgeGene(hidden, nodes[j]))
		}
	}

	return &Genome{
		ID:         id,
		NumInputs:  numInputs,
		NumOutputs: numOutputs,
		NumHidden:  numHidden,
		NodeGenes:  nodes,
		EdgeGenes:  edges,
		Fitness:    0.0,
		Winner:     false,
	}
}

// ToString summarizes the genome's connectivity in a string.
func (g *Genome) ToString() string {
	str := fmt.Sprintf("Genome(%d):\n", g.ID)
	for i := 0; i < len(g.EdgeGenes)-1; i++ {
		var conn string
		if g.EdgeGenes[i].Disabled {
			conn = "~~      ~~"
		} else {
			weight := g.EdgeGenes[i].Weight
			if g.EdgeGenes[i].Weight > 0.0 {
				conn = fmt.Sprintf("( %3.4f )", weight)
			} else {
				conn = fmt.Sprintf("( %3.3f )", weight)
			}
		}
		str += fmt.Sprintf("%6s(%3d, %8s) --%s--> %6s(%3d, %8s)\n",
			g.EdgeGenes[i].InputNode.Type,
			g.EdgeGenes[i].InputNode.ID,
			g.EdgeGenes[i].InputNode.AFuncType, conn,
			g.EdgeGenes[i].OutputNode.Type,
			g.EdgeGenes[i].OutputNode.ID,
			g.EdgeGenes[i].OutputNode.AFuncType)
	}

	// for the last edge gene
	edge := g.EdgeGenes[len(g.EdgeGenes)-1]
	var conn string
	if edge.Disabled {
		conn = "~~      ~~"
	} else {
		if edge.Weight > 0.0 {
			conn = fmt.Sprintf("( %3.4f )", edge.Weight)
		} else {
			conn = fmt.Sprintf("( %3.3f )", edge.Weight)
		}
	}
	str += fmt.Sprintf("%6s(%3d, %8s) --%s--> %6s(%3d, %8s)",
		edge.InputNode.Type,
		edge.InputNode.ID,
		edge.InputNode.AFuncType, conn,
		edge.OutputNode.Type,
		edge.OutputNode.ID,
		edge.OutputNode.AFuncType)

	return str
}

func (g *Genome) Export() error {
	// genome_[id]_[exported time].txt
	f, err := os.Create("genome_%d_%d", g.ID, time.Now().UnixNano())
	if err != nil {
		return err
	}

	// node data
	for _, node := range g.NodeGene {
		dat := fmt.Sprintf("n %d %s %s", node.ID, node.Type, node.AFuncType)
		_, err := f.WriteString(dat + "\n")
		if err != nil {
			return err
		}
	}

	// edge data
	for _, edge := range g.EdgeGene {
		dat := fmt.Sprintf("e %d %d %f", edge.InputNode, node.Type, node.AFuncType)
		_, err := f.WriteString(dat + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// pathSearch checks if there is a path from the start node to the goal node
// in the genome, and therefore the network is recurrent, when the start node
// is connected from the goal node.
func (g *Genome) pathSearch(from, to *NodeGene) bool {
	// if reached the goal
	if from.ID == to.ID {
		return true
	}

	for _, edge := range g.EdgeGenes {
		if edge.InputNode.ID == from.ID {
			if g.pathSearch(edge.OutputNode, to) {
				return true
			}
		}
	}

	return false
}

// Mutate mutates this genome given the rate of mutation by adding a node and
// by adding an edge. Return the ID of a newly added node and the IDs of the
// nodes are connected by the newly added edge.
func (g *Genome) Mutate(addNodeRate, addEdgeRate float64) (int, int, int) {
	nid, from, to := -1, -1, -1
	if rand.Float64() < addNodeRate {
		nid = g.AddNode()
	}
	if rand.Float64() < addEdgeRate {
		from, to = g.AddEdge()
	}
	return nid, from, to
}

// AddNode randomly selects an edge in the genome and adds a node that
// separates the connection by the edge, disabling the edge. Return the ID
// of the newly added node.
func (g *Genome) AddNode() int {
	newNode := NewNodeGene(len(g.NodeGenes), "hidden", randAFuncName())
	edgeIndex := rand.Intn(len(g.EdgeGenes))
	edge := g.EdgeGenes[edgeIndex]

	g.NodeGenes = append(g.NodeGenes, newNode)
	g.NumHidden++

	g.EdgeGenes[edgeIndex].Disabled = true
	g.EdgeGenes = append(g.EdgeGenes, NewEdgeGene(edge.InputNode, newNode))
	g.EdgeGenes = append(g.EdgeGenes, NewEdgeGene(newNode, edge.OutputNode))

	return newNode.ID
}

// AddEdge randomly selects two nodes in the genome and connect the two with
// a new edge, while ensuring that the graph keeps its feedforwarding
// property. Return the IDs of the nodes that are connected by the newly added
// edge (in order of input node and output node), -1 and -1 otherwise.
func (g *Genome) AddEdge() (int, int) {
	input := g.NodeGenes[rand.Intn(len(g.NodeGenes))]  // input node
	output := g.NodeGenes[rand.Intn(len(g.NodeGenes))] // output node

	// check if there is already an edge with these two nodes
	for _, edge := range g.EdgeGenes {
		if edge.InputNode.ID == input.ID && edge.OutputNode.ID == output.ID {
			return -1, -1
		}
	}

	// Output node of the edge cannot be an input type node.
	if output.Type == "input" {
		return -1, -1
	}

	// If there is a path from the output node to the input node, then adding
	// an edge from the input node to the output node causes the network to be
	// recurrent.
	if g.pathSearch(output, input) {
		return -1, -1
	}

	g.EdgeGenes = append(g.EdgeGenes, NewEdgeGene(input, output))
	return input.ID, output.ID
}

// Crossover takes another genome, performs crossover, then replace this
// genome with the resulting child.
func (g *Genome) Crossover(g0 *Genome) error {
	if g.NumInputs != g0.NumInputs || g.NumOutputs != g0.NumOutputs {
		return errors.New("Invalid number of inputs/outputs provided")
	}

	nodeCopies := make([]*NodeGene, 0, len(g0.NodeGenes))
	for _, node := range g0.NodeGenes {
		nodeCopies = append(nodeCopies,
			NewNodeGene(node.ID, node.Type, node.AFuncType))
	}

	edgeCopies := make([]*EdgeGene, 0, len(g0.EdgeGenes))

	sort.Slice(nodeCopies, func(i, j int) bool {
		return nodeCopies[i].ID < nodeCopies[j].ID
	})

	var inputNode, outputNode *NodeGene
	for _, edge := range g0.EdgeGenes {
		// search if any of the input/output nodes is added yet
		inputID, outputID := edge.InputNode.ID, edge.OutputNode.ID

		// search for input node
		index := sort.Search(len(nodeCopies), func(i int) bool {
			return nodeCopies[i].ID >= inputID
		})
		if index < len(nodeCopies) && nodeCopies[index].ID == inputID {
			if nodeCopies[index].Type == "input" {
				inputNode = g.NodeGenes[index]
			} else {
				inputNode = nodeCopies[index]
			}
		}

		// search for output node
		index = sort.Search(len(nodeCopies), func(i int) bool {
			return nodeCopies[i].ID >= outputID
		})
		if index < len(nodeCopies) && nodeCopies[index].ID == outputID {
			if nodeCopies[index].Type == "output" {
				outputNode = g.NodeGenes[index]
			} else {
				outputNode = nodeCopies[index]
			}
		}

		// create a new copy of the edge
		edgeCopy := NewEdgeGene(inputNode, outputNode)
		edgeCopy.Weight = edge.Weight
		edgeCopy.Disabled = edge.Disabled

		// check if this genome already has an edge with the same input node
		// and output node, when the new edge copy has an input node as an
		// input and an output node as an output.
		if edgeCopy.InputNode.Type == "input" &&
			edgeCopy.OutputNode.Type == "output" {
			for _, edge0 := range g.EdgeGenes {
				if !(edge0.InputNode.ID == edgeCopy.InputNode.ID &&
					edge0.OutputNode.ID == edgeCopy.OutputNode.ID) {
					edgeCopies = append(edgeCopies, edgeCopy)
				}
			}
		} else {
			edgeCopies = append(edgeCopies, edgeCopy)
		}
	}

	// update node copies' IDs
	for _, node := range nodeCopies {
		node.ID += len(g.NodeGenes) - (g0.NumInputs + g0.NumOutputs)
	}

	nodeCopies = nodeCopies[g0.NumInputs+g0.NumOutputs:]
	g.NodeGenes = append(g.NodeGenes, nodeCopies...)
	g.EdgeGenes = append(g.EdgeGenes, edgeCopies...)

	g.NumHidden += g0.NumHidden

	return nil
}
