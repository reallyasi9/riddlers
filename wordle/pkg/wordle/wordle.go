package wordle

import (
	"fmt"
	"math"
)

type Wordle struct {
	solutions []Word
}

func (w Wordle) Len() int {
	return len(w.solutions)
}

func (w Wordle) Solutions() []Word {
	return w.solutions
}

func NewWordle(words []Word) *Wordle {
	ws := make([]Word, len(words))
	copy(ws, words)

	wordle := Wordle{solutions: ws}
	return &wordle
}

func (w *Wordle) Filter(guesses []Word, ws []WordStatus) int {
	ps := NewPlayStatus()
	for i := range guesses {
		ps.UpdateWithGuess(guesses[i], ws[i])
	}

	last := len(w.solutions) - 1
	for i := range w.solutions {
		j := last - i
		soln := w.solutions[j]
		if !ps.Possible(soln) {
			w.solutions = append(w.solutions[:j], w.solutions[j+1:]...)
		}
	}

	return len(w.solutions)
}

const MAX_GUESSES = 6

type guessIndices [MAX_GUESSES]uint8

func (w Wordle) tryGroups(guesses []Word) map[guessIndices]int {
	if len(guesses) > MAX_GUESSES {
		panic(fmt.Errorf("maximum of %d guesses allowed", MAX_GUESSES))
	}
	indices := make(map[Word]guessIndices)
	for _, soln := range w.solutions {
		for i, guess := range guesses {
			status := guess.Compare(soln)
			index := indices[soln]
			index[i] = uint8(status.TrinaryIndex())
			indices[soln] = index
		}
	}

	// Invert into count of indices
	groups := make(map[guessIndices]int)
	for _, index := range indices {
		groups[index]++
	}

	return groups
}

func (w Wordle) Try(guesses []Word) (float64, float64, int) {
	groups := w.tryGroups(guesses)

	n := float64(len(w.solutions))
	entropy := float64(0)
	prob := float64(len(groups)) / n
	deduced := int(0)
	for _, count := range groups {
		if count == 0 {
			continue
		}
		if count == 1 {
			deduced++
		}
		entropy += math.Log2(float64(count)) * float64(count) / n
	}
	return entropy, prob, deduced
}
