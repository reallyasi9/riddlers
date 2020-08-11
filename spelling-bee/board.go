package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gonum/stat/combin"
)

// A Board represents a spelling bee hexagramatical board.
type Board struct {
	center  rune
	letters map[rune]struct{}
}

// NewBoard creates a board object from a string of letters.
func NewBoard(letters string) (*Board, error) {
	if len(letters) < 7 {
		return nil, fmt.Errorf("board needs at least seven letters")
	}
	l := make(map[rune]struct{})
	var center rune
	for _, letter := range letters {
		lower := unicode.ToLower(letter)
		if lower == 's' {
			return nil, fmt.Errorf("board cannot contain letter 's'")
		}
		if lower != letter {
			center = lower
		}
		l[lower] = struct{}{}
	}
	if len(l) < 7 {
		return nil, fmt.Errorf("board needs seven unique letters")
	}
	return &Board{center: center, letters: l}, nil
}

// NewBoards creates as many unique boards from a word as there are pangrams.
func NewBoards(word string) ([]*Board, error) {
	word = strings.ToLower(word)
	if len(word) < 7 {
		return nil, fmt.Errorf("board needs at least seven letters")
	}
	letters := make(map[rune]struct{})
	for _, letter := range word {
		if letter == 's' {
			return nil, fmt.Errorf("board cannot contain letter 's'")
		}
		letters[letter] = struct{}{}
	}
	if len(letters) < 7 {
		return nil, fmt.Errorf("board needs at least one pangram")
	}
	unique := make([]rune, len(letters))
	i := 0
	for letter := range letters {
		unique[i] = letter
		i++
	}
	gen := combin.NewCombinationGenerator(len(letters), 7)
	boards := make([]*Board, 0)
	comb := make([]int, 7)
	for gen.Next() {
		gen.Combination(comb)
		l := make(map[rune]struct{})
		for _, i := range comb {
			l[unique[i]] = struct{}{}
		}
		for i := 0; i < 7; i++ {
			boards = append(boards, &Board{center: unique[comb[i]], letters: l})
		}
	}
	return boards, nil
}

// Check if a word is valid for a board.  Returns two values: a bool describing wheather or not the word is valid, and an int counting the number of letters of the board used.
func (b *Board) Check(word string) (bool, int) {
	// the word must be at least 4 letters long
	if len(word) < 4 {
		return false, 0
	}
	center := false
	used := make(map[rune]bool)
	for _, letter := range word {
		// words cannot use letters not on the board
		if _, ok := b.letters[letter]; !ok {
			return false, 0
		}
		// words must use the central element
		if letter == b.center {
			center = true
		}
		// count the letters used
		used[letter] = true
	}
	return center, len(used)
}

// Score a word and return whether or not it is a pangram.
func (b *Board) Score(word string) (int, bool) {
	valid, letters := b.Check(word)
	if !valid {
		return 0, false
	}
	if len(word) == 4 {
		return 1, false
	}
	bonus := 0
	if letters == 7 {
		bonus = 7
	}
	return len(word) + bonus, bonus == 7
}

func (b Board) String() string {
	var sb strings.Builder
	for letter := range b.letters {
		if letter == b.center {
			sb.WriteString(strings.ToUpper(string(letter)))
		} else {
			sb.WriteString(strings.ToLower(string(letter)))
		}
	}
	return sb.String()
}
