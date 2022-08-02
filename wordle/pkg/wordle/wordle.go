package wordle

import (
	"github.com/kelindar/bitmap"
)

const ALPHABET_SIZE = 26

type Wordle struct {
	words []Word

	wordsContainingLetter           [ALPHABET_SIZE]bitmap.Bitmap
	wordsContainingLetterByPosition [WORD_SIZE][ALPHABET_SIZE]bitmap.Bitmap
	wordsNotContainingLetter        [ALPHABET_SIZE]bitmap.Bitmap
}

func NewWordle(words []Word) *Wordle {
	ws := make([]Word, len(words))
	copy(ws, words)

	wordle := Wordle{words: ws}
	for i, word := range words {
		for letter := byte(0); letter < ALPHABET_SIZE; letter++ {
			found := false
			for pos, character := range word {
				if character == letter {
					wordle.wordsContainingLetterByPosition[pos][character].Set(uint32(i))
					wordle.wordsContainingLetter[character].Set(uint32(i))
					found = true
				}
			}
			if !found {
				wordle.wordsNotContainingLetter[letter].Set(uint32(i))
			}
		}
	}

	return &wordle
}

func (w *Wordle) GetWord(i int) Word {
	return w.words[i]
}

func (w *Wordle) NWords() int {
	return len(w.words)
}

func (w *Wordle) Ambiguities(guesses []Word, soln Word) bitmap.Bitmap {
	var possible bitmap.Bitmap
	possible.Grow(uint32(len(w.words)))
	possible.Ones()

	for _, guess := range guesses {
		status := guess.Compare(soln)
		for letter, stat := range status {
			switch stat {
			case CORRECT:
				possible.And(w.wordsContainingLetterByPosition[letter][guess[letter]-1])
			case PRESENT:
				possible.And(w.wordsContainingLetter[guess[letter]-1])
				possible.AndNot(w.wordsContainingLetterByPosition[letter][guess[letter]-1])
			case ABSENT:
				possible.And(w.wordsNotContainingLetter[guess[letter]-1])
			default:
			}
		}
	}
	return possible
}
