package mga

// CompFn defines a type of function that compares two fitness scores (float64)
// and returns true if the first argument fitness score is more fit than the
// second, false if otherwise.
type CompFn func(float64, float64) bool

// DirectComparison returns a comparison function that returns true if the
// first argument fitness score is higher than the second.
func DirectComparison() CompFn {
	return func(score0, score1 float64) bool {
		return score0 > score1
	}
}

// InverseComparison returns a comparison function that returns true if the
// first argument fitness score is lower than the second.
func InverseComparison() CompFn {
	return func(score0, score1 float64) bool {
		return score0 < score1
	}
}
