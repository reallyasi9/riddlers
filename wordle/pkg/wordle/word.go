package wordle

const WORD_SIZE = 5
const ZERO_CHAR = 'a' - 1

type Word [WORD_SIZE]byte

func NewWord(bs []byte) Word {
	var w Word
	var i int
	var b byte
	for i, b = range bs {
		if i >= WORD_SIZE {
			break
		}
		w[i] = b - ZERO_CHAR
	}
	return w
}

func NewWordFromString(s string) Word {
	return NewWord([]byte(s))
}

func (w Word) String() string {
	out := [WORD_SIZE]byte{}
	for i, c := range w {
		out[i] = c + ZERO_CHAR
	}
	return string(out[:])
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
