package main

import (
	"sort"
	"sync"
)

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

// Iterate iterates the generation, mating the top survivors and mutating the rest randomly
func (g Generation) Iterate(survivors, spawn int, temperature float64) {
	// TODO: figure out some way of implementing offspring
	// offspring := survivors / 2
	// offspring := 0
	// mutations := g.Len() - survivors - offspring
	// if mutations < 0 {
	// 	log.Panicf("survivors (%d) plus offspring (%d) greater than generation size (%d)", survivors, offspring, g.Len())
	// }

	sort.Sort(sort.Reverse(g))

	// Pair off survivors
	// for i := 0; i < survivors; i += 2 {
	// 	g[survivors+i/2].ReplaceWithOffspring(&g[i], &g[i+1])
	// }

	// Clone survivors as mutants
	var wg sync.WaitGroup
	for i := survivors; i < g.Len()-spawn; i++ {
		wg.Add(1)
		go func(i int) {
			g[i].ReplaceWithMutation(&g[i%survivors], temperature)
			wg.Done()
		}(i)
	}

	// Spawn new
	for i := g.Len() - spawn; i < g.Len(); i++ {
		wg.Add(1)
		go func(i int) {
			g[i] = *NewBoard()
			wg.Done()
		}(i)
	}
	wg.Wait()
}
