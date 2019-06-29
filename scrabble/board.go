package main

import (
	"fmt"
	"log"
	"math"
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

	rng *rand.Rand
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
	b := &Board{Raw: make([]rune, len(letterCollection)), rng: rand.New(rand.NewSource(rand.Int63()))}
	copy(b.Raw, letterCollection)
	b.Shuffle()
	return b
}

// Shuffle shuffles a permutation
func (b *Board) Shuffle() {
	b.rng.Shuffle(len(b.Raw), func(i, j int) {
		b.Raw[i], b.Raw[j] = b.Raw[j], b.Raw[i]
	})
	b.replaceQMs()
	b.score()
}

func (b *Board) replaceQMs() {
	var builder strings.Builder
	for _, r := range b.Raw {
		if r == '?' {
			builder.WriteRune(rune(b.rng.Int31n('z'-'a') + 'a'))
		} else {
			builder.WriteRune(r)
		}
	}
	b.Clean = builder.String()
}

type prefixWalker struct {
	isPrefix bool
}

func (p *prefixWalker) walk(s string, v interface{}) bool {
	p.isPrefix = true
	return false // stop iterating immediately
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

// Score scores a scrabble permutation
func (b *Board) score() {
	b.Score = 0
	sw := b.ScoreWords()
	for _, v := range sw {
		b.Score += v.Score
	}
}

// Mutate the board a bit in place
func (b *Board) Mutate(temp float64) {
	exponent := -float64(b.Score) / temp
	p := math.Exp(exponent)

	b.rng.Shuffle(len(b.Raw), func(i, j int) {
		if b.rng.Float64() < p {
			b.Raw[i], b.Raw[j] = b.Raw[j], b.Raw[i]
		}
	})
	b.replaceQMs()
	b.score()
}

// ReplaceWithMutation replaces this board with a mutated version of the parent
// FIXME: This makes a potentially illegal board.
func (b *Board) ReplaceWithMutation(b2 *Board, temperature float64) {
	copy(b.Raw, b2.Raw)
	b.Mutate(temperature)
	b.replaceQMs()
	b.score()
}

// ReplaceWithOffspring replaces this board with an offspring based on the parents
// FIXME: This makes a potentially illegal board.
func (b *Board) ReplaceWithOffspring(b1, b2 *Board) {
	for i := 0; i < b.Len(); i++ {
		if b.rng.Float32() < .5 {
			b.Raw[i] = b1.Raw[i]
		} else {
			b.Raw[i] = b2.Raw[i]
		}
	}
	b.replaceQMs()
	b.score()
}

// ScoredWord is a word with an associated score
type ScoredWord struct {
	Word  string
	Score int
}

// ScoreWords finds and scores all the words in the board
func (b *Board) ScoreWords() []ScoredWord {
	words := make([]ScoredWord, 0)
	foundTrie := radix.New()
	var walker prefixWalker

	i := 0
	j := 1
	for i < b.Len() {
		sub := b.Clean[i:j]

		if val, ok := scoreTrie.Get(sub); ok {
			if _, alreadyFound := foundTrie.Get(sub); !alreadyFound {
				// Is a word: score it
				score := val.(int)
				// Replace ? values
				raw := b.Raw[i:j]
				for k, r := range raw {
					if r == '?' {
						score -= runeScores[rune(sub[k])]
						score += runeScores[r]
					}
				}
				words = append(words, ScoredWord{Word: sub, Score: val.(int)})
				foundTrie.Insert(sub, true)
			}
		}

		walker.isPrefix = false
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
	return words
}

func (b Board) String() string {
	var builder strings.Builder
	for i, r := range b.Raw {
		if r == '?' {
			builder.WriteRune('[')
			builder.WriteByte(b.Clean[i])
			builder.WriteRune(']')
		} else {
			builder.WriteRune(r)
		}
	}
	return fmt.Sprintf("%s %d", strings.ToUpper(builder.String()), b.Score)
}
