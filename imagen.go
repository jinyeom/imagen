/*


imagen.go implementation of image generation via DPPN and mGA.

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
	"github.com/gonum/matrix/mat64"
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
	fmt.Println("  imagen [filename].png [config].json")
}

func draw(g *Genome, width, height int) {
	n, _ := NewDPPN(g, 1)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fx := float64(x)
			fy := float64(y)
			d := math.Sqrt((fx-float64(width)/2.0)*(fx-float64(width)/2.0) +
				(fy-float64(height)/2.0)*(fy-float64(height)/2.0))
			inputs := []float64{fx * 0.1, fy * 0.1, d * 0.1, 1.0}
			inputVec := mat64.NewDense(1, 4, inputs)

			outputVec, _ := n.FeedForward(inputVec)
			outputs := outputVec.RawMatrix().Data

			c := color.RGBA{uint8(outputs[0] * 255.0), uint8(outputs[1] * 255.0),
				uint8(outputs[2] * 255.0), 255}
			img.Set(x, y, c)
		}
	}

	f1, err := os.Create(fmt.Sprintf("estimated_%d.png", g.ID))
	if err != nil {
		fmt.Println(err)
	}
	defer f1.Close()

	png.Encode(f1, img)
}

// genImage returns an evaluation function for fitting the argument image's
// pixel value distribution.
func genImage(img *image.RGBA, numBatch, numEpochs int,
	learningRate float64) EvaluationFunc {
	width, height := img.Bounds().Max.X-img.Bounds().Min.X,
		img.Bounds().Max.Y-img.Bounds().Min.Y

	return func(g *Genome) float64 {
		n, _ := NewDPPN(g, numBatch)
		score := 0.0

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
				inputs = append(inputs, fx*0.1, fy*0.1, d*0.1, 1.0)

				// target
				c := img.At(x, y).(color.RGBA)
				target = append(target, float64(c.R)/255.0,
					float64(c.G)/255.0, float64(c.B)/255.0)
			}

			inputBatch := mat64.NewDense(numBatch, 4, inputs)
			targetBatch := mat64.NewDense(numBatch, 3, target)

			mse, err := n.Backprop(inputBatch, targetBatch, learningRate)
			if err != nil {
				panic(err)
			}
			score += mse
		}

		// encode the DPPN's learned connection weights back to its genotype.
		n.Encode(g)

		return score
	}
}

func main() {
	if len(os.Args) != 3 {
		help()
		return
	}

	imgFile := os.Args[1]
	configFile := os.Args[2]

	// image file
	f, err := os.Open(imgFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	// config file
	config, err := NewConfiguration(configFile)
	if err != nil {
		panic(err)
	}

	rand.Seed(config.Seed)

	env, err := NewMGA(config,
		InverseComparison(),
		genImage(img.(*image.RGBA), config.BatchSize,
			config.NumEpochs, config.LearningRate))
	if err != nil {
		panic(err)
	}
	env.Run(true, true)

	// export all the images and genomes in the population
	width, height := img.Bounds().Max.X-img.Bounds().Min.X,
		img.Bounds().Max.Y-img.Bounds().Min.Y
	for _, genome := range env.Population {
		draw(genome, width, height)
		genome.Export()
	}
}
