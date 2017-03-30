package dppn

import (
	"fmt"
	"github.com/gonum/matrix/mat64"
	. "github.com/jinyeom/ean/mga"
	//"image"
	//"image/color"
	//"image/png"
	"log"
	"math/rand"
	//"os"
	"testing"
)

// DPPNUnitTest performs a series of white box testing.
func DPPNUnitTest() {
	fmt.Println("=== DPPN Unit Test ===")

	// Test NewDPPN()
	fmt.Println("=== Test NewDPPN() ===")
	n, err := NewDPPN(NewGenome(0, 2, 3, 1), 4)
	if err != nil {
		log.Fatal("Error: %s\n", err)
	}

	fmt.Println(n.ToString())

	fmt.Printf("Expected number of inputs: %d\n", 2)
	fmt.Printf("Actual number of inputs: %d\n", n.NumInputs)

	fmt.Printf("Expected number of outputs: %d\n", 1)
	fmt.Printf("Actual number of outputs: %d\n", n.NumOutputs)

	fmt.Printf("Expected number of nodes: %d\n", 6)
	fmt.Printf("Actual number of nodes: %d\n", len(n.Nodes))

	numEdges := 0
	for _, node := range n.Nodes {
		numEdges += len(node.Inputs)
	}

	fmt.Printf("Expected number of edges: %d\n", 9)
	fmt.Printf("Actual number of edges: %d\n", numEdges)

	if n.NumInputs != 2 || n.NumOutputs != 1 ||
		len(n.Nodes) != 6 || numEdges != 9 {
		fmt.Println("TEST \033[1;31mFAILED\033[0m")
	} else {
		fmt.Println("TEST \033[1;32mPASSED\033[0m")
	}
}

// DPPNAcceptanceTest performs a series of black box testing.
func DPPNAcceptanceTest() {
	fmt.Println("=== DPPN Acceptance Test ===")

	// DPPN XOR Test
	fmt.Println("=== DPPN XOR Test ===")
	n1, err := NewDPPN(NewGenome(1, 3, 4, 1), 4)
	if err != nil {
		log.Fatal("Error: %s\n", err)
	}

	inputs := mat64.NewDense(
		4, 3, []float64{
			1.0, 1.0, 1.0,
			1.0, 0.0, 1.0,
			0.0, 1.0, 1.0,
			0.0, 0.0, 1.0,
		},
	)

	target := mat64.NewDense(
		4, 1, []float64{
			0.0,
			1.0,
			1.0,
			0.0,
		},
	)

	for i := 0; i < 1000; i++ {
		mse, err := n1.Backprop(inputs, target, 0.5)
		if err != nil {
			log.Fatal("Error: %s\n", err)
		}
		// average MSE
		if i%100 == 0 {
			fmt.Printf("Average MSE at epoch %d: %f\n", i, mse)
		}
	}

	outputs, err := n1.FeedForward(inputs)
	if err != nil {
		log.Fatal("Error: %s\n", err)
	}

	outputFmt := mat64.Formatted(outputs, mat64.Prefix("    "), mat64.Squeeze())
	fmt.Printf("a = %v\n\n", outputFmt)
}

func TestDPPN(t *testing.T) {
	rand.Seed(0)

	// Perform white box testing on DPPN
	DPPNUnitTest()

	// Perform black box testing on DPPN
	DPPNAcceptanceTest()
}
