package wordle

import "strings"

const WORD_SIZE = 5

type Word [WORD_SIZE]byte

func NewWord(bs []byte) Word {
	var w Word
	var i int
	var b byte
	for i, b = range bs[:WORD_SIZE] { // truncate
		w[i] = b
	}
	return w
}

func NewWordFromString(s string) Word {
	return NewWord([]byte(strings.ToLower(s)))
}

func (w Word) String() string {
	return string(w[:])
}

func (w Word) Compare(soln Word) WordStatus {
	var status WordStatus

	// correct first
	for i, c := range w {
		if soln[i] == c {
			status[i] = CORRECT
			soln[i] = 0 // prevent further matches
		}
	}
	// present second
OUTER:
	for i, c := range w {
		if status[i] == CORRECT {
			continue
		}
		for j, x := range soln {
			if c == x {
				status[i] = PRESENT
				soln[j] = 0 // prevent further matches
				continue OUTER
			}
		}
	}
	// default is absent

	return status
}
