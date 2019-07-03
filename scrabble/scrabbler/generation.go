package scrabbler

import (
	"log"
	"math"
	"math/rand"
	"time"
)

// Generation represents a single generation of scrabble permutations
type Generation struct {
	Board            *Board
	Steps            int
	StartTemperature float64
	TemperatureStep  float64
	ReportEvery      int
	Patience         int

	rng *rand.Rand
}

// NewGeneration creates a new Generation object
func NewGeneration(steps int, start, step float64, report int, patience int, rng *rand.Rand) *Generation {
	return &Generation{Board: NewBoard(rng), Steps: steps, StartTemperature: start, TemperatureStep: step, ReportEvery: report, Patience: patience, rng: rng}
}

// Energy is the probability of transition from state a to state b
func energy(a, b *Board, t float64) float64 {
	if b.Score > a.Score {
		return 1.
	}
	return math.Exp(-float64(a.Score) / (float64(b.Score) * t))
}

func copyBoard(to, from *Board) {
	copy(to.Raw, from.Raw)
	to.Clean = from.Clean
	to.Score = from.Score
}

// Anneal performs the simulated annealing
func (g *Generation) Anneal() {

	swapBoard := &Board{Raw: make([]byte, g.Board.Len()), Clean: g.Board.Clean, Score: g.Board.Score}
	copy(swapBoard.Raw, g.Board.Raw)

	lt := math.Log(g.StartTemperature)
	startTime := time.Now()
	lastStep := 0
	for itr := 0; itr < g.Steps; itr++ {
		swapBoard.Nudge(g.rng)
		energy := energy(g.Board, swapBoard, math.Exp(lt))
		if g.rng.Float64() < energy {
			copyBoard(g.Board, swapBoard)
			lastStep = itr
			lt -= g.TemperatureStep
		} else if itr-lastStep > g.Patience {
			return
		} else {
			// reset
			copyBoard(swapBoard, g.Board)
		}
		if g.ReportEvery > 0 && itr%g.ReportEvery == 0 {
			log.Printf("%4.0f Sec:  Step %8d  LTmp %8F  Enrg %8F  %s", time.Now().Sub(startTime).Seconds(), itr, lt, energy, g.Board)
		}
	}
}
