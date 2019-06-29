package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	radix "github.com/armon/go-radix"
)

// Board represents a permutation of scrabble letters
type Board struct {
	// Raw is raw permutation data, wilds included
	Raw []rune
	// Clean is the cleaned permutation, wilds replaced with random values
	Clean string
	// Score is the score of the board
	Score int
}

var letterCollection = []rune("??aaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyyz")

var runeScores = map[rune]int{
	'?': 0,
	'e': 1,
	'a': 1,
	'i': 1,
	'o': 1,
	'n': 1,
	'r': 1,
	't': 1,
	'l': 1,
	's': 1,
	'u': 1,
	'd': 2,
	'g': 2,
	'b': 3,
	'c': 3,
	'm': 3,
	'p': 3,
	'f': 4,
	'h': 4,
	'v': 4,
	'w': 4,
	'y': 4,
	'k': 5,
	'j': 8,
	'x': 8,
	'q': 10,
	'z': 10,
}

var scoreTrie = radix.New()

// Len returns the number of letters in the scrabble permutation
func (b *Board) Len() int {
	return len(b.Raw)
}

// NewBoard creates a new scrabble permutation
func NewBoard() *Board {
	b := &Board{Raw: make([]rune, len(letterCollection))}
	copy(b.Raw, letterCollection)
	rand.Shuffle(len(b.Raw), func(i, j int) {
		b.Raw[i], b.Raw[j] = b.Raw[j], b.Raw[i]
	})
	b.Clean = replaceQMs(b.Raw)
	b.score()
	return b
}

func replaceQMs(board []rune) string {
	var builder strings.Builder
	for _, r := range board {
		if r == '?' {
			builder.WriteRune(rune(rand.Int31n('z'-'a') + 'a'))
		} else {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

// Score scores a scrabble permutation
func (b *Board) score() {
	b.Score = 0
	foundTrie := radix.New()
	i := 0
	j := 1
	for i < b.Len() {
		sub := b.Clean[i:j]

		if val, ok := scoreTrie.Get(sub); ok {
			if _, alreadyFound := foundTrie.Get(sub); !alreadyFound {
				// Is a word: score it
				b.Score += val.(int)
				foundTrie.Insert(sub, true)
			}
		}

		var walker prefixWalker
		scoreTrie.WalkPrefix(sub, walker.walk)
		if walker.isPrefix {
			// Is a prefix: increase j
			j++
			if j > b.Len() {
				// No more letters
				break
			}
		} else {
			// Scoot forward
			i++
			if i == j {
				j++
			}
		}
	}
}

func scoreWord(word string) (int, error) {
	score := 0
	for _, r := range word {
		value, ok := runeScores[r]
		if !ok {
			err := fmt.Errorf("Rune %c not a recognized scrabble letter", r)
			log.Printf("Error looking up word: %s", err)
			return 0, err
		}
		score += value
	}
	return score, nil
}

func (b Board) String() string {
	wilds := make([]rune, 0)
	for i, r := range b.Raw {
		if r == '?' {
			wilds = append(wilds, rune(b.Clean[i]))
		}
	}
	return fmt.Sprintf("%s (%s) %d", string(b.Raw), string(wilds), b.Score)
}
