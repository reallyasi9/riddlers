package main

import (
	"fmt"
	"math/rand"
	"sort"
)

const nPlayers = 2000
const nSoldiers = 100
const nCastles = 10
const nSims = 10000
const nSurvive = 250 // Should be < nPlayers
const nNew = 1500    // Should be < nPlayers - nSurvive

type slice struct {
	sort.Float64Slice
	idx []int
}

func (s slice) Swap(i, j int) {
	s.Float64Slice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func newSlice(f []float64) *slice {
	s := &slice{Float64Slice: sort.Float64Slice(f), idx: make([]int, len(f))}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

func main() {

	picks := generateDistributions(nPlayers, nCastles, nSoldiers)
	best := 0.

	for iSim := 0; iSim < nSims; iSim++ {
		scores := battle(picks)
		scoreSlice := newSlice(scores)
		sort.Sort(sort.Reverse(scoreSlice))

		if scoreSlice.Float64Slice[0] > best {
			best = scoreSlice.Float64Slice[0]
			fmt.Printf("(%d) Score: %f\nTop dist:  %v\n", iSim, scoreSlice.Float64Slice[0], picks[scoreSlice.idx[0]])
		}

		newPicks := make([][]float64, nPlayers)
		for iPlayer := 0; iPlayer < nSurvive; iPlayer++ {
			newPicks[iPlayer] = picks[scoreSlice.idx[iPlayer]]
		}
		for iPlayer := nSurvive; iPlayer < nPlayers-nNew; iPlayer++ {
			j := rand.Intn(nSurvive)
			k := rand.Intn(nSurvive)
			newPicks[iPlayer] = mean(picks[scoreSlice.idx[j]], picks[scoreSlice.idx[k]])
		}
		newDistributions := generateDistributions(nNew, nCastles, nSoldiers)
		newPicks = append(newPicks[:nPlayers-nNew], newDistributions...)
		picks = newPicks

	}

}

func generateDistributions(p, c, s int) [][]float64 {
	picks := make([][]float64, p)
	for i := 0; i < p; i++ {
		divs := generateDivisions(c, s)
		picks[i] = distributeSoldiers(divs, s)
	}
	return picks
}

func generateDivisions(c, s int) []int {
	divs := make([]int, c-1)
	for i := 0; i < c-1; i++ {
		divs[i] = rand.Intn(s + 1)
	}
	return divs
}

func distributeSoldiers(d []int, ns int) []float64 {
	soldiers := make([]float64, len(d)+1)
	sort.Ints(d)
	soldiers[0] = float64(d[0])
	for i := 1; i < len(d); i++ {
		soldiers[i] = float64(d[i] - d[i-1])
	}
	soldiers[len(d)] = float64(nSoldiers - d[len(d)-1])
	return soldiers
}

func battle(picks [][]float64) []float64 {
	scores := make([]float64, len(picks))
	for i := 0; i < len(picks); i++ {
		for j := i + 1; j < len(picks); j++ {
			for k := 0; k < len(picks[i]); k++ {
				if picks[i][k] > picks[j][k] {
					scores[i] += float64(k) + 1.
				} else if picks[i][k] < picks[j][k] {
					scores[j] += float64(k) + 1.
				} else {
					scores[i] += (float64(k) + 1.) / 2.
					scores[j] += (float64(k) + 1.) / 2.
				}
			}
		}
	}
	return scores
}

func mean(p1, p2 []float64) []float64 {
	r := make([]float64, len(p1))
	for i := 0; i < len(p1); i++ {
		r[i] = (p1[i] + p2[i]) / 2.
	}
	return r
}
