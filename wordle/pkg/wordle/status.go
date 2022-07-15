package wordle

import (
	"github.com/kelindar/bitmap"
	"github.com/segmentio/fasthash/fnv1a"
)

type LetterStatusCode int

const (
	ABSENT LetterStatusCode = iota
	PRESENT
	CORRECT
)

type WordStatus [WORD_SIZE]LetterStatusCode

type PlayStatus struct {
	absent               bitmap.Bitmap
	presentWrongPosition [wordle.WORD_SIZE]bitmap.Bitmap
	otherwisePresent     bitmap.Bitmap
	correct              Word
}

func NewStatus() *PlayStatus {
	pwp := [WORD_SIZE]bitmap.Bitmap{}
	for i := range pwp {
		pwp[i] = bitmap.Bitmap{}
	}
	return &PlayStatus{
		absent:               bitmap.Bitmap{},
		presentWrongPosition: pwp,
		otherwisePresent:     bitmap.Bitmap{},
		correct:              Word{},
	}
}

func (ps *PlayStatus) Possible(soln Word) bool {
	for i, c := range soln {
		if ps.correct[i] != ZERO_CHAR && c != ps.correct[i] {
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
	absent := ps.absent.Clone(nil)
	prp := [WORD_SIZE]bitmap.Bitmap{}
	for i := range prp {
		ps.presentWrongPosition[i].Clone(&prp[i])
	}
	op := ps.otherwisePresent.Clone(nil)
	correct := ps.correct

	return &PlayStatus{
		absent:               absent,
		presentWrongPosition: prp,
		otherwisePresent:     op,
		correct:              correct,
	}
}

func (ps *PlayStatus) Hash() uint64 {
	h := fnv1a.Init64
	for _, val := range ps.absent {
		h = fnv1a.AddUint64(h, val)
	}
	for _, bs := range ps.presentWrongPosition {
		for _, val := range bs {
			h = fnv1a.AddUint64(h, val)
		}
	}
	for _, val := range ps.otherwisePresent {
		h = fnv1a.AddUint64(h, val)
	}
	h = fnv1a.AddBytes64(h, ps.correct[:])

	return h
}
