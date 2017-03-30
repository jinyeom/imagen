package main

import (
	"encoding/json"
	"fmt"
	. "github.com/gonum/matrix/mat64"
	. "github.com/jinyeom/imagen/dppn"
	. "github.com/jinyeom/imagen/mga"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/rand"
	"os"
)

func help() {
	fmt.Println("EAN (Evolutionary Adversarial Nets)")
	fmt.Println("Copyright (c) 2017 by Jin Yeom")
	fmt.Println("User Manual:")
	fmt.Println("\tean [config].json")
}

func draw(g *Genome) {
	width := 50
	height := 41

	n, _ := NewDPPN(g, 1)
	imge := image.NewGray(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fx := float64(x)
			fy := float64(y)
			d := math.Sqrt((fx-float64(width)/2.0)*(fx-float64(width)/2.0) +
				(fy-float64(height)/2.0)*(fy-float64(height)/2.0))
			inputs := []float64{fx * 0.2, fy * 0.2, d * 0.2, 1.0}
			inputVec := NewDense(1, 4, inputs)

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

func ImageFit(numBatch, numEpochs int, lr float64) EvalFn {
	// target image
	f, err := os.Open("butterfly.png")
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

	return func(g *Genome) float64 {
		n, _ := NewDPPN(g, numBatch)
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

			inputBatch := NewDense(numBatch, 4, inputs)
			targetBatch := NewDense(numBatch, 1, target)

			mse, _ := n.Backprop(inputBatch, targetBatch, lr)

			avg += mse
		}

		// encode the DPPN's learned connection weights back to its genotype.
		n.Encode(g)

		return avg / float64(numEpochs)
	}
}

func main() {
	if len(os.Args) != 2 {
		help()
		return
	}

	// import configuration
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	var config MGAConfig
	if err := dec.Decode(&config); err != nil {
		if err != io.EOF {
			fmt.Println(err)
			return
		}
	}

	env, err := NewMGA(&config, InverseComparison(), ImageFit(16, 2000, 0.1))
	if err != nil {
		fmt.Println(err)
		return
	}
	env.Run(true, true)

	for _, genome := range env.Population {
		draw(genome)
	}
}
