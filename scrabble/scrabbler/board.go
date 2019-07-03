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
	Raw []byte
	// Clean is the cleaned permutation, wilds replaced with random values
	Clean string
	// Score is the score of the board
	Score int

	rng *rand.Rand
}

var letterCollection = []byte("??aaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyyz")
var letterCounts map[byte]int

func init() {
	letterCounts = make(map[byte]int)
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
	b := &Board{Raw: make([]byte, len(letterCollection)), rng: rand.New(rand.NewSource(rand.Int63()))}
	copy(b.Raw, letterCollection)
	b.Shuffle()
	return b
}

func validateBoard(r []byte) error {
	if len(r) != len(letterCollection) {
		return fmt.Errorf("number of letters in board (%d) does not match the number of tiles (%d)", len(r), len(letterCollection))
	}
	myCounts := make(map[byte]int)
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

	raw := make([]byte, 0)
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
		raw = append(raw, s[i])
	}

	err := validateBoard(raw)
	if err != nil {
		log.Panic(err)
	}

	b := &Board{Raw: raw, Clean: string(clean), rng: rand.New(rand.NewSource(rand.Int63()))}
	b.score()
	return b
}

func (b *Board) replaceQMs() {
	var builder strings.Builder
	for _, r := range b.Raw {
		if r == '?' {
			builder.WriteByte(byte(b.rng.Intn('z'-'a') + 'a'))
		} else {
			builder.WriteByte(r)
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

// Nudge the board a bit in place.  The number of characters randomly shuffled is determined by the temperature.
func (b *Board) Nudge(temperature float64) {
	n := int(math.Ceil(math.Exp(-float64(b.Score)/temperature) * float64(b.Len())))

	for i := 0; i < n; i++ {
		j1 := b.rng.Intn(b.Len())
		j2 := b.rng.Intn(b.Len())
		b.Raw[j1], b.Raw[j2] = b.Raw[j2], b.Raw[j1]
	}

	b.replaceQMs()
	b.score()
}

// Shuffle the board in place
func (b *Board) Shuffle() {
	b.rng.Shuffle(b.Len(), func(i, j int) {
		b.Raw[i], b.Raw[j] = b.Raw[j], b.Raw[i]
	})

	b.replaceQMs()
	b.score()
}

// ReplaceWithNudge replaces this board with a mutated version of the parent
func (b *Board) ReplaceWithNudge(b2 *Board, temperature float64) {
	copy(b.Raw, b2.Raw)
	b.Score = b2.Score
	b.Nudge(temperature)
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
			builder.WriteByte('[')
			builder.WriteByte(b.Clean[i])
			builder.WriteByte(']')
		} else {
			builder.WriteByte(r)
		}
	}
	return fmt.Sprintf("%s %d", strings.ToUpper(builder.String()), b.Score)
}
