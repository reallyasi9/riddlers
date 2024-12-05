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
		panic(fmt.Sprintf("word status length must be %d", WORD_SIZE))
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

func (ws WordStatus) TrinaryIndex() int {
	idx := 0
	pow := 1
	for _, code := range ws {
		idx += pow * int(code)
		pow *= 3
	}
	return idx
}

const N_STATUS_OUTCOMES = 243 // 3^5

const N_LETTERS = 26

const ZERO_LETTER = 'a'

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
		possible[i] = bitmap.Bitmap{1<<(N_LETTERS+1) - 1}
	}
	maximumPresent := [N_LETTERS]int{}
	for c := range maximumPresent {
		maximumPresent[c] = WORD_SIZE
	}
	return &PlayStatus{
		possible:       possible,
		minimumPresent: [N_LETTERS]int{},
		maximumPresent: maximumPresent,
	}
}

func (ps *PlayStatus) Possible(soln Word) bool {
	letterCounts := [N_LETTERS]int{}
	for i, c := range soln {
		cint := uint32(c) - ZERO_LETTER
		if !ps.possible[i].Contains(cint) {
			return false
		}
		letterCounts[cint]++
		if letterCounts[cint] > ps.maximumPresent[cint] {
			return false
		}
	}
	for cint, n := range letterCounts {
		if n < ps.minimumPresent[cint] {
			return false
		}
	}
	return true
}

func (ps *PlayStatus) UpdateWithGuess(word Word, ws WordStatus) {
	localMins := [N_LETTERS]int{}
	localMaxs := [N_LETTERS]int{}
	for i := range localMaxs {
		localMaxs[i] = WORD_SIZE
	}

	for i, st := range ws {
		cint := uint32(word[i]) - ZERO_LETTER
		switch st {
		case ABSENT:
			// Only eliminate from this position, but set the known maximum letter count
			ps.possible[i].Remove(cint)
			localMaxs[cint] = localMins[cint]
		case PRESENT:
			// Eliminate from this position only
			ps.possible[i].Remove(cint)
			// Increase the known number of minimum counts
			localMins[cint]++
			// Trick: increase the known number of maximum counts to deal with a CORRECT or PRESENT after an ABSENT
			localMaxs[cint]++
		case CORRECT:
			// Eliminate all other options from this position
			ps.possible[i].Clear()
			ps.possible[i].Set(cint)
			// Increase the known number of minimum counts
			localMins[cint]++
			// Trick: increase the known number of maximum counts to deal with a CORRECT or PRESENT after an ABSENT
			localMaxs[cint]++
		}
	}
	// Update counts
	for i := range localMaxs {
		if ps.minimumPresent[i] < localMins[i] {
			ps.minimumPresent[i] = localMins[i]
		}
		if ps.maximumPresent[i] > localMaxs[i] {
			if localMaxs[i] > 5 {
				ps.maximumPresent[i] = 5
			} else {
				ps.maximumPresent[i] = localMaxs[i]
			}
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
