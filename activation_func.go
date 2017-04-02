/*


activation_func.go implementation of activation functions for DPPN.

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
	"math"
)

var (
	// aFuncSet is a list of activation functions that can be used by the DPPN.
	// Each function can be called via GetAFunc function.
	aFuncSet = map[string]*ActivationFunc{
		"identity": Identity(),
		"sigmoid":  Sigmoid(),
		"tanh":     Tanh(),
		"relu":     ReLU(),
		"elu":      ELU(1.0),
		"abs":      Abs(),
		"sine":     Sin(),
		"gaussian": Gaussian(0.0, 1.0),
	}
)

type ActivationFunc struct {
	Name string                // activation function name
	Fn   func(float64) float64 // transfer function
	DFn  func(float64) float64 // differential for backprop
}

func Identity() *ActivationFunc {
	return &ActivationFunc{
		Name: "identity",
		Fn: func(x float64) float64 {
			return x
		},
		DFn: func(x float64) float64 {
			return 1.0
		},
	}
}

func Sigmoid() *ActivationFunc {
	return &ActivationFunc{
		Name: "sigmoid",
		Fn: func(x float64) float64 {
			x = math.Max(-60.0, math.Min(60.0, 5.0*x))
			return 1.0 / (1.0 + math.Exp(-x))
		},
		DFn: func(x float64) float64 {
			x = math.Max(-60.0, math.Min(60.0, 5.0*x))
			sig := 1.0 / (1.0 + math.Exp(-x))
			return sig * (1.0 - sig)
		},
	}
}

func Tanh() *ActivationFunc {
	return &ActivationFunc{
		Name: "tanh",
		Fn: func(x float64) float64 {
			x = math.Max(-60.0, math.Min(60.0, 2.5*x))
			return math.Tanh(x)
		},
		DFn: func(x float64) float64 {
			x = math.Max(-60.0, math.Min(60.0, 2.5*x))
			return 1.0 - math.Pow(math.Tanh(x), 2.0)
		},
	}
}

func ReLU() *ActivationFunc {
	return &ActivationFunc{
		Name: "relu",
		Fn: func(x float64) float64 {
			if x >= 0.0 {
				return x
			}
			return 1.0
		},
		DFn: func(x float64) float64 {
			if x >= 0.0 {
				return 1.0
			}
			return 0.0
		},
	}
}

func ELU(a float64) *ActivationFunc {
	return &ActivationFunc{
		Name: "elu",
		Fn: func(x float64) float64 {
			if x >= 0.0 {
				return x
			}
			x = math.Max(-60.0, math.Min(60.0, 2.5*x))
			return a * (math.Exp(x) - 1.0)
		},
		DFn: func(x float64) float64 {
			if x >= 0.0 {
				return 1.0
			}
			x = math.Max(-60.0, math.Min(60.0, 2.5*x))
			return a*(math.Exp(x)-1.0) + a
		},
	}
}

func Abs() *ActivationFunc {
	return &ActivationFunc{
		Name: "abs",
		Fn: func(x float64) float64 {
			return math.Abs(x)
		},
		DFn: func(x float64) float64 {
			if x >= 0.0 {
				return 1.0
			}
			return -1.0
		},
	}
}

func Sin() *ActivationFunc {
	return &ActivationFunc{
		Name: "sine",
		Fn: func(x float64) float64 {
			return math.Sin(x)
		},
		DFn: func(x float64) float64 {
			return math.Cos(x)
		},
	}
}

func Gaussian(mu, sigma float64) *ActivationFunc {
	return &ActivationFunc{
		Name: "gaussian",
		Fn: func(x float64) float64 {
			x = math.Max(-3.4, math.Min(3.4, x))
			return (1.0 / math.Sqrt(2.0*sigma*math.Pi)) *
				math.Exp(-math.Pow(x-mu, 2.0)/(2.0*sigma*sigma))
		},
		DFn: func(x float64) float64 {
			x = math.Max(-3.4, math.Min(3.4, x))
			return (1.0 / math.Sqrt(2.0*sigma*math.Pi)) *
				math.Exp(-math.Pow(x-mu, 2.0)/(2.0*sigma*sigma)) *
				(mu - x) / (sigma * sigma)
		},
	}
}
