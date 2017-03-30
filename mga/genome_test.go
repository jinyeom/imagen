package mga

import (
	"fmt"
	"math/rand"
	"testing"
)

// GenomeUnitTest performs a series of white box testing.
func GenomeUnitTest() {
	fmt.Println("=== UNIT TEST ===")

	// Test Genome.NewGenome() (constructor)
	fmt.Println("=== Test Genome.NewGenome() ===")
	g0 := NewGenome(0, 3, 4, 2)
	fmt.Println(g0.ToString())

	expectedNodes := 9
	expectedEdges := 20

	fmt.Printf("Expected number of node genes: %d\n", expectedNodes)
	fmt.Printf("Actual number of node genes: %d\n", len(g0.NodeGenes))
	if len(g0.NodeGenes) != expectedNodes {
		fmt.Println("TEST \033[1;31mFAILED\033[0m")
	} else {
		fmt.Println("TEST \033[1;32mPASSED\033[0m")
	}

	fmt.Printf("Expected number of edge genes: %d\n", expectedEdges)
	fmt.Printf("Actual number of edge genes: %d\n", len(g0.EdgeGenes))
	if len(g0.EdgeGenes) != expectedEdges {
		fmt.Println("TEST \033[1;31mFAILED\033[0m")
	} else {
		fmt.Println("TEST \033[1;32mPASSED\033[0m")
	}

	// Test Genome.pathSearch()
	fmt.Println("=== Test Genome.pathSearch() ===")
	from := g0.NodeGenes[0] // input node
	to := g0.NodeGenes[3]   // hidden node

	fmt.Println("Testing pathSearch from an input node to a hidden node")
	fmt.Printf("Expected result of path search: %t\n", true)
	fmt.Printf("Actual result of path search: %t\n", g0.pathSearch(from, to))
	if !g0.pathSearch(from, to) {
		fmt.Println("TEST \033[1;31mFAILED\033[0m")
	} else {
		fmt.Println("TEST \033[1;32mPASSED\033[0m")
	}

	to = g0.NodeGenes[7] // output node

	fmt.Println("Testing pathSearch from an input node to an output node")
	fmt.Printf("Expected result of path search: %t\n", true)
	fmt.Printf("Actual result of path search: %t\n", g0.pathSearch(from, to))
	if !g0.pathSearch(from, to) {
		fmt.Println("TEST \033[1;31mFAILED\033[0m")
	} else {
		fmt.Println("TEST \033[1;32mPASSED\033[0m")
	}

	fmt.Println("Testing pathSearch from an output node to an input node")
	fmt.Printf("Expected result of path search: %t\n", false)
	fmt.Printf("Actual result of path search: %t\n", g0.pathSearch(to, from))
	if g0.pathSearch(to, from) {
		fmt.Println("TEST \033[1;31mFAILED\033[0m")
	} else {
		fmt.Println("TEST \033[1;32mPASSED\033[0m")
	}

	// Test Genome.AddNode()
	fmt.Println("=== Test Genome.AddNode() ===")
	fmt.Println(g0.ToString())
	g0.AddNode()
	expectedNodes++
	expectedEdges += 2
	fmt.Println(g0.ToString())

	fmt.Printf("Expected number of node genes: %d\n", expectedNodes)
	fmt.Printf("Actual number of node genes: %d\n", len(g0.NodeGenes))
	if len(g0.NodeGenes) != expectedNodes {
		fmt.Println("TEST \033[1;31mFAILED\033[0m")
	} else {
		fmt.Println("TEST \033[1;32mPASSED\033[0m")
	}

	// Test Genome.AddEdge()
	fmt.Println("=== Test Genome.AddEdge() ===")
	fmt.Println(g0.ToString())
	newFrom, newTo := g0.AddEdge()
	fmt.Println(g0.ToString())

	if newFrom != -1 && newTo != -1 {
		expectedEdges++
		fmt.Printf("Expected number of edge genes: %d\n", expectedEdges)
		fmt.Printf("Actual number of edge genes: %d\n", len(g0.EdgeGenes))
		if len(g0.EdgeGenes) != expectedEdges {
			fmt.Println("TEST \033[1;31mFAILED\033[0m")
		} else {
			fmt.Println("TEST \033[1;32mPASSED\033[0m")
		}
	} else {
		fmt.Printf("Expected number of edge genes: %d\n", expectedEdges)
		fmt.Printf("Actual number of edge genes: %d\n", len(g0.EdgeGenes))
		if len(g0.EdgeGenes) != expectedEdges {
			fmt.Println("TEST \033[1;31mFAILED\033[0m")
		} else {
			fmt.Println("TEST \033[1;32mPASSED\033[0m")
		}
	}

	// Test Genome.Mutate()
	fmt.Println("=== Test Genome.Mutate() ===")
	fmt.Println(g0.ToString())
	newNID, newFrom, newTo := g0.Mutate(0.5, 0.5)
	fmt.Println(g0.ToString())

	if newNID != -1 {
		expectedNodes++
		expectedEdges += 2
		fmt.Printf("Expected number of node genes: %d\n", expectedNodes)
		fmt.Printf("Actual number of node genes: %d\n", len(g0.NodeGenes))
		if len(g0.NodeGenes) != expectedNodes {
			fmt.Println("TEST \033[1;31mFAILED\033[0m")
		} else {
			fmt.Println("TEST \033[1;32mPASSED\033[0m")
		}
	} else {
		fmt.Printf("Expected number of node genes: %d\n", expectedNodes)
		fmt.Printf("Actual number of node genes: %d\n", len(g0.NodeGenes))
		if len(g0.NodeGenes) != expectedNodes {
			fmt.Println("TEST \033[1;31mFAILED\033[0m")
		} else {
			fmt.Println("TEST \033[1;32mPASSED\033[0m")
		}
	}

	if newFrom != -1 && newTo != -1 {
		expectedEdges++
		fmt.Printf("Expected number of edge genes: %d\n", expectedEdges)
		fmt.Printf("Actual number of edge genes: %d\n", len(g0.EdgeGenes))
		if len(g0.EdgeGenes) != expectedEdges {
			fmt.Println("TEST \033[1;31mFAILED\033[0m")
		} else {
			fmt.Println("TEST \033[1;32mPASSED\033[0m")
		}
	} else {
		fmt.Printf("Expected number of edge genes: %d\n", expectedEdges)
		fmt.Printf("Actual number of edge genes: %d\n", len(g0.EdgeGenes))
		if len(g0.EdgeGenes) != expectedEdges {
			fmt.Println("TEST \033[1;31mFAILED\033[0m")
		} else {
			fmt.Println("TEST \033[1;32mPASSED\033[0m")
		}
	}

	// Test Genome.Crossover()
	fmt.Println("=== Test Genome.Crossover() ===")
	g1 := NewGenome(1, 3, 4, 2)
	for i := 0; i < 10; i++ {
		g0.Mutate(0.5, 0.5)
		g1.Mutate(0.5, 0.5)
	}

	fmt.Println(g0.ToString())
	fmt.Println(g1.ToString())

	expectedNodes = len(g0.NodeGenes) + len(g1.NodeGenes) - 5
	expectedEdges = len(g0.EdgeGenes) + len(g1.EdgeGenes)

	if err := g0.Crossover(g1); err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Println("TEST \033[1;31mFAILED\033[0m")
	} else {
		fmt.Println(g0.ToString())

		fmt.Printf("Expected number of node genes: %d\n", expectedNodes)
		fmt.Printf("Actual number of node genes: %d\n", len(g0.NodeGenes))
		if len(g0.NodeGenes) != expectedNodes {
			fmt.Println("TEST \033[1;31mFAILED\033[0m")
		} else {
			fmt.Println("TEST \033[1;32mPASSED\033[0m")
		}

		fmt.Printf("Expected number of edge genes: %d\n", expectedEdges)
		fmt.Printf("Actual number of edge genes: %d\n", len(g0.EdgeGenes))
		if len(g0.EdgeGenes) != expectedEdges {
			fmt.Println("TEST \033[1;31mFAILED\033[0m")
		} else {
			fmt.Println("TEST \033[1;32mPASSED\033[0m")
		}
	}
}

// GenomeAcceptanceTest performs a series of black box testing.
func GenomeAcceptanceTest() {
	fmt.Println("=== ACCEPTANCE TEST ===")

	// Genome black box test
	fmt.Println("=== Test Genome ===")
	g2 := NewGenome(2, 5, 8, 4)
	g3 := NewGenome(3, 5, 8, 4)

	g2ExpectedNodes := 17
	g2ExpectedEdges := 72

	g3ExpectedNodes := 17
	g3ExpectedEdges := 72

	for i := 0; i < 10; i++ {
		nid, from, to := g2.Mutate(0.5, 0.5)
		if nid != -1 {
			g2ExpectedNodes++
			g2ExpectedEdges += 2
		}
		if from != -1 && to != -1 {
			g2ExpectedEdges++
		}
		nid, from, to = g3.Mutate(0.5, 0.5)
		if nid != -1 {
			g3ExpectedNodes++
			g3ExpectedEdges += 2
		}
		if from != -1 && to != -1 {
			g3ExpectedEdges++
		}
	}

	if err := g2.Crossover(g3); err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Println("TEST \033[1;31mFAILED\033[0m")
	} else {
		g2ExpectedNodes += (g3ExpectedNodes - 9)
	}

	if len(g2.NodeGenes) == g2ExpectedNodes &&
		len(g3.EdgeGenes) == g3ExpectedEdges {
		fmt.Println("TEST \033[1;32mPASSED\033[0m")
	}
}

func TestGenome(t *testing.T) {
	rand.Seed(0)

	// white box test for each function of Genome
	GenomeUnitTest()

	// black box test for Genome
	GenomeAcceptanceTest()
}
