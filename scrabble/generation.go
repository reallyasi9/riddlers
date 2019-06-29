package main

// Generation represents a single generation of scrabble permutations
type Generation []Board

func (g Generation) Len() int           { return len(g) }
func (g Generation) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g Generation) Less(i, j int) bool { return g[i].Score < g[j].Score }

// MakeGeneration makes a new generation of length n
func MakeGeneration(n int) Generation {
	gen := make(Generation, n)
	for i := 0; i < n; i++ {
		gen[i] = *NewBoard()
	}
	return gen
}
