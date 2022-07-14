package wordle

import "github.com/kelindar/bitmap"

type LetterStatusCode int

const (
	ABSENT LetterStatusCode = iota
	PRESENT
	CORRECT
)

type WordStatus [WORD_SIZE]LetterStatusCode

type PlayStatus struct {
	Guesses      []Word
	WordStatuses []WordStatus

	absent               bitmap.Bitmap
	presentWrongPosition [WORD_SIZE]bitmap.Bitmap
	otherwisePresent     bitmap.Bitmap
	correct              Word
}

func NewStatus() *PlayStatus {
	pwp := [WORD_SIZE]bitmap.Bitmap{}
	for i := range pwp {
		pwp[i] = bitmap.Bitmap{}
	}
	return &PlayStatus{
		Guesses:              []Word{},
		WordStatuses:         []WordStatus{},
		absent:               bitmap.Bitmap{},
		presentWrongPosition: pwp,
		otherwisePresent:     bitmap.Bitmap{},
		correct:              Word{},
	}
}

func (ps *PlayStatus) Possible(soln Word) bool {
	for i, c := range soln {
		if ps.correct[i] != 0 && c != ps.correct[i] {
			return false
		}
		if ps.absent.Contains(uint32(c)) {
			return false
		}
		if ps.presentWrongPosition[i].Contains(uint32(c)) {
			return false
		}
		if !ps.otherwisePresent.Contains(uint32(c)) {
			return false
		}
	}
	return true
}

func (ps *PlayStatus) UpdateWithGuess(word Word, ws WordStatus) {
	ps.Guesses = append(ps.Guesses, word)
	ps.WordStatuses = append(ps.WordStatuses, ws)

	for i, st := range ws {
		switch st {
		case ABSENT:
			ps.absent.Set(uint32(word[i]))
		case PRESENT:
			ps.presentWrongPosition[i].Set(uint32(word[i]))
			ps.otherwisePresent.Set(uint32(word[i]))
		case CORRECT:
			ps.correct[i] = word[i]
		}
	}

}

func (ps *PlayStatus) Clone() *PlayStatus {
	guesses := make([]Word, len(ps.Guesses))
	copy(guesses, ps.Guesses)
	statuses := make([]WordStatus, len(ps.WordStatuses))
	copy(statuses, ps.WordStatuses)

	absent := ps.absent.Clone(nil)
	prp := [WORD_SIZE]bitmap.Bitmap{}
	for i := range prp {
		ps.presentWrongPosition[i].Clone(&prp[i])
	}
	op := ps.otherwisePresent.Clone(nil)
	correct := ps.correct

	return &PlayStatus{
		Guesses:              guesses,
		WordStatuses:         statuses,
		absent:               absent,
		presentWrongPosition: prp,
		otherwisePresent:     op,
		correct:              correct,
	}
}
