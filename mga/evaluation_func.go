package mga

// EvalFn defines a type of function that evaluates a genome and returns the
// fitness score of the genome.
type EvalFn func(g *Genome) float64
