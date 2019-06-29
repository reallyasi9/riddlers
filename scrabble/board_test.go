package main

import (
	"fmt"
	"testing"
)

func TestScore(t *testing.T) {
	buildDictionary(dictionaryURL)

	b := Board{Raw: []rune{'t', 'h', 'e', 'a', 't', 'e', 'r'}, Clean: "theater"}
	words := b.ScoreWords()
	for _, w := range words {
		fmt.Printf("%s: %d\n", w.Word, w.Score)
	}
}
