package mga

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
)

// MGAUnitTest is a series of white box testing for MGA.
func MGAUnitTest() {
	fmt.Println("=== Unit Test ===")

	// Test MGAConfig
	fmt.Println("=== Test MGAConfig ===")
	var config MGAConfig

	f, err := os.Open("test_config.json")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Println("TEST \033[1;31mFAILED\033[0m")
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	if err := dec.Decode(&config); err != nil {
		if err != io.EOF {
			fmt.Printf("Error: %s\n", err)
			fmt.Println("TEST \033[1;31mFAILED\033[0m")
		}
	}

	if config.NumInputs == 4 && config.NumOutputs == 1 &&
		config.NumInitHidden == 2 && config.PopulationSize == 50 &&
		config.NumTournaments == 500 && config.MutAddNodeRate == 0.5 &&
		config.MutAddEdgeRate == 0.7 {
		fmt.Println("TEST \033[1;32mPASSED\033[0m")
	}

	// Test NewMGA()
	cfn := InverseComparison()
	efn := func(numEpochs int, lr float64) EvalFn {
		return func(g *Genome) float64 {
			n := NewDPPN(g)
			avgmse := 0.0
			for i := 0; i < numEpochs; i++ {
				mse1, _ := n.Backprop([]float64{1.0, 1.0, 1.0},
					[]float64{0.0}, lr)
				mse2, _ := n.Backprop([]float64{0.0, 1.0, 1.0},
					[]float64{1.0}, lr)
				mse3, _ := n.Backprop([]float64{1.0, 0.0, 1.0},
					[]float64{1.0}, lr)
				mse4, _ := n.Backprop([]float64{0.0, 0.0, 1.0},
					[]float64{0.0}, lr)
				avgmse += (mse1 + mse2 + mse3 + mse4) / 4.0
			}
			return avgmse / float64(numEpochs)
		}
	}(1000, 0.01)
	ga := NewMGA(&config, cfn, efn)

}

// MGAAcceptanceTest is a series of black box testing for MGA.
func MGAAcceptanceTest() {

}

func TestMGA(t *testing.T) {
	rand.Seed(0)

	// white box testing for MGA
	MGAUnitTest()

	// black box testing for MGA
	MGAAcceptanceTest()
}
