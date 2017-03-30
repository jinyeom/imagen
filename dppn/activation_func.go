package dppn

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
