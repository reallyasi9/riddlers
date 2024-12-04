package wordle

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
	for _, i := range cut {
		w.solutions[i] = w.solutions[len(w.solutions)-1]
		w.solutions = w.solutions[:len(w.solutions)-1]
	}
	return len(w.solutions)
}

func (w *Wordle) Try(guesses []Word) float32 {
	groups := [][]Word{w.solutions}
	for _, guess := range guesses {
		newGroups := [][]Word{}
		for _, group := range groups {
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
	prob := float32(0)
	for _, group := range groups {
		prob += 1 / float32(len(group))
	}
	return prob / float32(len(w.solutions))
}
