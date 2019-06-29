package scrabbler

import (
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

	b := Board{Raw: []rune{'t', 'h', 'e', 'a', 't', 'e', 'r'}, Clean: "theater"}
	words := b.ScoreWords()
	for word, score := range words {
		t.Logf("%s: %d\n", word, score)
	}
}

func BenchmarkScore(b *testing.B) {
	board := NewBoard()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.ScoreWords()
	}
}
