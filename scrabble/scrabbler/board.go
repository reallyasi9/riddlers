package scrabbler

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
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
var letterCounts map[rune]int

func init() {
	letterCounts = make(map[rune]int)
	for _, r := range letterCollection {
		letterCounts[r]++
	}
}

// Len returns the number of letters in the scrabble permutation
func (b *Board) Len() int {
	return len(b.Raw)
}

// NewBoard creates a new scrabble permutation
func NewBoard() *Board {
	b := &Board{Raw: make([]rune, len(letterCollection)), rng: rand.New(rand.NewSource(rand.Int63()))}
	copy(b.Raw, letterCollection)
	b.Mutate(math.Inf(1))
	return b
}

func validateBoard(r []rune) error {
	if len(r) != len(letterCollection) {
		return fmt.Errorf("number of letters in board (%d) does not match the number of tiles (%d)", len(r), len(letterCollection))
	}
	myCounts := make(map[rune]int)
	for _, x := range r {
		myCounts[x]++
	}
	for key, value := range myCounts {
		n, ok := letterCounts[key]
		if !ok {
			return fmt.Errorf("letter '%c' not a valid tile", key)
		}
		if value != n {
			return fmt.Errorf("count (%d) of letter '%c' not equal to number of tiles of that letter (%d)", value, key, n)
		}
	}
	return nil
}

// MakeBoard makes a board from a string.
// The format is either capital or lower-case letters with wild tiles surrounded by square brackets,
// like "AbcDE[f]GHijKLM[N]opq"
func MakeBoard(s string) *Board {
	s = strings.ToLower(s)

	raw := make([]rune, 0)
	clean := make([]byte, 0)
	for i := 0; i < len(s); i++ {
		if s[i] == '[' {
			if i+2 >= len(s) || s[i+2] != ']' {
				log.Panicf("missing closing bracket ']' after character %d", i)
			}
			clean = append(clean, s[i+1])
			raw = append(raw, '?')
			i++
			continue
		}

		if s[i] == ']' {
			if i-2 < 0 || s[i-2] != '[' {
				log.Panicf("missing opening bracket '[' before character %d", i)
			}
			continue
		}

		clean = append(clean, s[i])
		raw = append(raw, rune(s[i]))
	}

	err := validateBoard(raw)
	if err != nil {
		log.Panic(err)
	}

	b := &Board{Raw: raw, Clean: string(clean), rng: rand.New(rand.NewSource(rand.Int63()))}
	b.score()
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

// Score scores a scrabble permutation
func (b *Board) score() {
	b.Score = 0
	sw := b.ScoreWords()
	for _, v := range sw {
		b.Score += v
	}
}

// Mutate the board a bit in place.  The higher the temperature, the more random the shuffle.
func (b *Board) Mutate(temperature float64) {
	p := math.Exp(-float64(b.Score) / temperature)
	b.rng.Shuffle(len(b.Raw), func(i, j int) {
		if b.rng.Float64() < p {
			b.Raw[i], b.Raw[j] = b.Raw[j], b.Raw[i]
		}
	})
	b.replaceQMs()
	b.score()
}

// ReplaceWithMutation replaces this board with a mutated version of the parent
func (b *Board) ReplaceWithMutation(b2 *Board, temperature float64) {
	copy(b.Raw, b2.Raw)
	b.replaceQMs()
	b.Mutate(temperature)
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

// ScoreWords finds and scores all the words in the board
func (b *Board) ScoreWords() map[string]int {
	found := make(map[string]int)

Loop:
	for i := 0; i < b.Len(); i++ {

		var wb strings.Builder
		scoreModifier := 0
		r := b.Clean[i]

		branch := ScoreTrie.Step(r)
		if branch == nil {
			continue // That letter doesn't start a word?  Huh.
		}

		// Modify score
		if b.Raw[i] == '?' {
			scoreModifier -= runeScores[r-'a']
		}

		wb.WriteByte(r)
		for j := i + 1; j < b.Len(); j++ {

			r = b.Clean[j]
			branch = branch.Step(r)
			if branch == nil {
				// Not a prefix: break out of j loop
				continue Loop
			}

			// Modify score
			if b.Raw[j] == '?' {
				scoreModifier -= runeScores[r-'a']
			}

			wb.WriteByte(r)
			if branch.Score > 0 {
				// Is a word: add it to the list with a modified score
				word := wb.String()
				score := branch.Score + scoreModifier
				if score > found[word] {
					found[word] = score
				}
			}

		}
	}

	return found
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
