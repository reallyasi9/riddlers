package wordle

const WORD_SIZE = 5

type Word [WORD_SIZE]byte

func NewWord(bs []byte) Word {
	var w Word
	var i int
	var b byte
	for i, b = range bs {
		if i >= WORD_SIZE {
			break
		}
		w[i] = b
	}
	return w
}

func NewWordFromString(s string) Word {
	return NewWord([]byte(s))
}

func (w Word) String() string {
	return string(w[:])
}

func (w Word) Compare(soln Word) WordStatus {
	var status WordStatus

OUTER:
	for i, c := range w {
		for j, x := range soln {
			if c == x {
				if i == j {
					status[i] = CORRECT
					continue OUTER
				} else {
					status[i] = PRESENT
				}
			}
		}
	}

	return status
}
