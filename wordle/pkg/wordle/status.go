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

const N_LETTERS = 26

type PlayStatus struct {
	// Possible solutions for each position
	possible [WORD_SIZE]bitmap.Bitmap
	// Minimum number of each letter present
	minimumPresent [N_LETTERS]int
	// Maximum number of each letter present
	maximumPresent [N_LETTERS]int
}

func NewPlayStatus() *PlayStatus {
	possible := [WORD_SIZE]bitmap.Bitmap{}
	for i := range possible {
		possible[i] = bitmap.Bitmap{(1 << (N_LETTERS + 1)) - 1}
	}
	maximumPresent := [N_LETTERS]int{}
	for c := range maximumPresent {
		maximumPresent[c] = -1
	}
	return &PlayStatus{
		possible:       possible,
		minimumPresent: [N_LETTERS]int{},
		maximumPresent: maximumPresent,
	}
}

func (ps *PlayStatus) Possible(soln Word) bool {
	letterCounts := make(map[uint32]int)
	for i, c := range soln {
		cint := uint32(c)
		if !ps.possible[i].Contains(cint) {
			return false
		}
		letterCounts[cint]++
	}
	for cint, n := range letterCounts {
		if n < ps.minimumPresent[cint-1] {
			return false
		}
		if ps.maximumPresent[cint-1] >= 0 && n > ps.maximumPresent[cint-1] {
			return false
		}
	}
	return true
}

func (ps *PlayStatus) UpdateWithGuess(word Word, ws WordStatus) {
	letterCounts := make(map[uint32]int)
	maxFound := make(map[uint32]struct{})
	for i, st := range ws {
		cint := uint32(word[i])
		switch st {
		case ABSENT:
			// Only eliminate from positions that are not solved
			for j := 0; j < WORD_SIZE; j++ {
				if ps.possible[j].Count() > 1 {
					ps.possible[j].Remove(cint)
				}
			}
			letterCounts[cint] += 0 // set to 0 if not in map, otherwise make no change
			maxFound[cint] = struct{}{}
		case PRESENT:
			// Eliminate from this position only
			ps.possible[i].Remove(cint)
			letterCounts[cint]++
		case CORRECT:
			// Eliminate all other options
			ps.possible[i].Clear()
			ps.possible[i].Set(cint)
			letterCounts[cint]++
		}
	}
	// Update counts
	for cint, n := range letterCounts {
		if n > ps.minimumPresent[cint-1] {
			ps.minimumPresent[cint-1] = n
		}
		if _, ok := maxFound[cint]; ok {
			ps.maximumPresent[cint-1] = n
		}
	}
}

func (ps *PlayStatus) Clone() *PlayStatus {
	possible := [WORD_SIZE]bitmap.Bitmap{}
	for i := range possible {
		ps.possible[i].Clone(&possible[i])
	}
	return &PlayStatus{
		possible:       possible,
		minimumPresent: ps.minimumPresent,
		maximumPresent: ps.maximumPresent,
	}
}

func (ps *PlayStatus) Hash() uint64 {
	h := fnv1a.Init64
	for _, possible := range ps.possible {
		for _, val := range possible {
			h = fnv1a.AddUint64(h, val)
		}
	}
	for _, n := range ps.minimumPresent {
		h = fnv1a.AddUint64(h, uint64(n))
	}
	for _, n := range ps.maximumPresent {
		h = fnv1a.AddUint64(h, uint64(n))
	}

	return h
}
