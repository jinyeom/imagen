package main

import (
	"fmt"
	"github.com/gonum/matrix/mat64"
	"github.com/jinyeom/imagen/dppn"
	"github.com/jinyeom/imagen/mga"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
)

func help() {
	fmt.Println("Imagen (Image Generation via DPPN)")
	fmt.Println("Copyright (c) 2017 by Jin Yeom")
	fmt.Println("User Manual:")
	fmt.Println("\timagen [config].json [filename].png")
}

func draw(g *mga.Genome) {
	width := 50
	height := 41

	n, _ := dppn.NewDPPN(g, 1)
	imge := image.NewGray(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fx := float64(x)
			fy := float64(y)
			d := math.Sqrt((fx-float64(width)/2.0)*(fx-float64(width)/2.0) +
				(fy-float64(height)/2.0)*(fy-float64(height)/2.0))
			inputs := []float64{fx * 0.2, fy * 0.2, d * 0.2, 1.0}
			inputVec := mat64.NewDense(1, 4, inputs)

			outputVec, _ := n.FeedForward(inputVec)
			outputs := outputVec.RawMatrix().Data

			c := color.Gray{uint8(outputs[0] * 255.0)}
			imge.Set(x, y, c)
		}
	}

	f1, err := os.Create(fmt.Sprintf("estimated_%d.png", g.ID))
	if err != nil {
		fmt.Println(err)
	}
	defer f1.Close()

	png.Encode(f1, imge)
}

// genImage returns an evaluation function for fitting the argument image's
// pixel value distribution.
func genImage(filename string, numBatch, numEpochs int, lr float64) mga.EvalFn {
	// target image
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		fmt.Println(err)
	}

	width, height := img.Bounds().Max.X-img.Bounds().Min.X,
		img.Bounds().Max.Y-img.Bounds().Min.Y

	return func(g *mga.Genome) float64 {
		n, _ := dppn.NewDPPN(g, numBatch)
		avg := 0.0

		for i := 0; i < numEpochs; i++ {
			// process a random batch of inputs and target outputs
			inputs := make([]float64, 0, 4*numBatch)
			target := make([]float64, 0, numBatch)
			for j := 0; j < numBatch; j++ {
				x := rand.Intn(width)
				y := rand.Intn(height)

				// input
				fx := float64(x)
				fy := float64(y)
				d := math.Sqrt((fx-float64(width)/2.0)*(fx-float64(width)/2.0) +
					(fy-float64(height)/2.0)*(fy-float64(height)/2.0))
				inputs = append(inputs, fx*0.2, fy*0.2, d*0.2, 1.0)

				// target
				c := img.At(x, y).(color.Gray)
				target = append(target, float64(c.Y)/255.0)
			}

			inputBatch := mat64.NewDense(numBatch, 4, inputs)
			targetBatch := mat64.NewDense(numBatch, 1, target)

			mse, _ := n.Backprop(inputBatch, targetBatch, lr)

			avg += mse
		}

		// encode the DPPN's learned connection weights back to its genotype.
		n.Encode(g)

		return avg / float64(numEpochs)
	}
}

func main() {
	if len(os.Args) != 3 {
		help()
		return
	}

	config, err := mga.NewMGAConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	env, err := mga.NewMGA(config,
		mga.InverseComparison(),             // comparison function
		genImage(os.Args[2], 16, 2000, 0.1)) // evaluation function
	if err != nil {
		panic(err)
	}
	env.Run(true, true)

	// export all the images in the population
	for _, genome := range env.Population {
		draw(genome)
	}
}
