package main

import (
	"fmt"
	"testing"
)

func TestScore(t *testing.T) {
	var st ScrabbleTrie
	st.Insert("the")
	st.Insert("he")
	st.Insert("heat")
	st.Insert("heater")
	st.Insert("eat")
	st.Insert("eater")
	st.Insert("at")
	st.Insert("ate")
	st.Insert("er")
	st.Insert("theater")

	b := Board{Raw: []rune{'t', 'h', 'e', 'a', 't', 'e', 'r'}, Clean: "theater", scoreTrie: &st}
	words := b.ScoreWords()
	for word, score := range words {
		fmt.Printf("%s: %d\n", word, score)
	}
}

func BenchmarkScore(b *testing.B) {
	st, _ := buildDictionary(dictionaryURL)
	board := NewBoard(st)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.ScoreWords()
	}
}
