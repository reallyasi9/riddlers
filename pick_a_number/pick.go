package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

const nPlayers = 3000
const nSurvive = 50 // the top n survive to the next round--should be an even number <= nPlayers/2
const minN = 1
const maxN = 1000000000
const target = 2. / 3. // fraction of mean of submitted answers
const nSim = 100000
const nWarmup = 100 // how long to warm up before starting real simulations

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
	// Generate the inital picks randomly
	picks := randomizePicks(nPlayers)
	results := make([]int64, nSim)
	distances := make([]float64, nSim)

	for iSim := nWarmup; iSim < nSim; iSim++ {
		t := meanInt64(picks) * target
		dists := newSlice(distance(picks, t))
		// Sort to find who is closest
		sort.Sort(dists)
		if iSim%10000 == 0 {
			fmt.Printf("%d: Best so far: %d (score %f)\n", iSim, picks[dists.idx[0]], dists.Float64Slice[0])
		}
		if iSim >= 0 {
			results[iSim] = picks[dists.idx[0]]
			distances[iSim] = dists.Float64Slice[0]
		}
		// Prune those who are too far away from the target
		bestIdxs := dists.idx[:nSurvive]
		newPicks := make([]int64, nPlayers)
		for i, idx := range bestIdxs {
			newPicks[i] = picks[idx]
		}
		// Pair these up and add their "children"
		j := nSurvive
		for i := 0; i < nSurvive; i += 2 {
			newPicks[j] = (picks[i] + picks[i+1]) / 2
			j++
		}
		// Randomize the rest
		picks = append(newPicks[:j], randomizePicks(nPlayers-j)...)
	}

	// Mean of the results
	fmt.Printf("Mean pick after %d generations: %f (score %f)\n", nSim+nWarmup, meanInt64(results), meanFloat64(distances))
}

func randomizePicks(n int) []int64 {
	picks := make([]int64, n)
	for i := 0; i < n; i++ {
		picks[i] = rand.Int63n(maxN-minN) + minN
	}
	return picks
}

func meanInt64(picks []int64) float64 {
	m := float64(0.)
	for _, val := range picks {
		m += float64(val)
	}
	return m / float64(len(picks))
}

func meanFloat64(picks []float64) float64 {
	m := 0.
	for _, val := range picks {
		m += val
	}
	return m / float64(len(picks))
}

func distance(picks []int64, t float64) []float64 {
	dists := make([]float64, len(picks))
	for i, val := range picks {
		dists[i] = math.Abs(t - float64(val))
	}
	return dists
}
