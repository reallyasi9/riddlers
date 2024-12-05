package wordle

import "math"

const ALPHABET_SIZE = 26

type Wordle struct {
	solutions []Word
	status    *PlayStatus
}

func NewWordle(words []Word) *Wordle {
	ws := make([]Word, len(words))
	ss := make([]Word, len(words))
	copy(ws, words)
	copy(ss, words)

	wordle := Wordle{solutions: ws, status: NewPlayStatus()}
	return &wordle
}

func (w *Wordle) GetSolution(i int) Word {
	return w.solutions[i]
}

func (w *Wordle) NSolutions() int {
	return len(w.solutions)
}

func (w *Wordle) AddGuess(guess Word, ws WordStatus) int {
	w.status.UpdateWithGuess(guess, ws)
	cut := []int{}
	for i, soln := range w.solutions {
		if !w.status.Possible(soln) {
			cut = append(cut, i)
		}
	}
	last := len(cut) - 1
	for i := range cut {
		c := cut[last-i]
		w.solutions = append(w.solutions[:c], w.solutions[c+1:]...)
	}
	return len(w.solutions)
}

func (w *Wordle) Try(guesses []Word) float64 {
	groups := [][]Word{w.solutions}
	for _, guess := range guesses {
		newGroups := [][]Word{}
		for _, group := range groups {
			if len(group) == 0 {
				continue
			}
			if len(group) == 1 {
				newGroups = append(newGroups, group)
				continue
			}
			outcomeTabulation := [N_STATUS_OUTCOMES][]Word{}
			for _, soln := range group {
				status := guess.Compare(soln)
				outcomeTabulation[status.TrinaryIndex()] = append(outcomeTabulation[status.TrinaryIndex()], soln)
			}
			for _, solns := range outcomeTabulation {
				if len(solns) > 0 {
					newGroups = append(newGroups, solns)
				}
			}
		}
		groups = newGroups
	}
	entropy := float64(0)
	for _, group := range groups {
		entropy += math.Log2(float64(len(group)))
	}
	return entropy
}
