package wordle

import (
	"fmt"

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

func NewWordStatus(s string) WordStatus {
	if len(s) != WORD_SIZE {
		panic("word status length must be 5")
	}
	var ws WordStatus
	for i, r := range s {
		switch r {
		case '+':
			ws[i] = CORRECT
		case '-':
			ws[i] = ABSENT
		case '?':
			ws[i] = PRESENT
		default:
			panic(fmt.Sprintf("words status character '%c' not recognized", r))
		}
	}
	return ws
}

type PlayStatus struct {
	absent               bitmap.Bitmap
	presentWrongPosition [WORD_SIZE]bitmap.Bitmap
	otherwisePresent     bitmap.Bitmap
	correct              Word
}

func NewPlayStatus() *PlayStatus {
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
	var otherwiseFound bitmap.Bitmap
	for i, c := range soln {
		if ps.correct[i] != 0 && c != ps.correct[i] {
			return false
		}
		cint := uint32(c)
		if ps.absent.Contains(cint) {
			return false
		}
		if ps.presentWrongPosition[i].Contains(cint) {
			return false
		}
		if ps.otherwisePresent.Contains(cint) {
			otherwiseFound.Set(cint)
		}
	}
	otherwiseFound.Xor(ps.otherwisePresent)
	return otherwiseFound.Count() == 0
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
