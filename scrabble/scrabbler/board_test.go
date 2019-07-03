package scrabbler

import (
	"testing"
)

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

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

func TestMake(t *testing.T) {
	b := MakeBoard("[a][a]aaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyyz")
	expected := "aaaaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyyz"
	if b.Clean != expected {
		t.Errorf("board does not match: expected '%s', got '%s'", expected, b.Clean)
	}

	b = MakeBoard("yyZ[X]AaaaaAaaabBccddddeeEEeeeeeeeeffgGgHHiiiiiiiiIJkllllmmnnnnnnoOOOOoooppqrrrrrrsssstTTTttuuu[y]uvvwwx")
	expected = "yyzxaaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuyuvvwwx"
	if b.Clean != expected {
		t.Errorf("board does not match: expected '%s', got '%s'", expected, b.Clean)
	}

	assertPanic(t, func() {
		MakeBoard("NotEnoughLetters")
	})
	assertPanic(t, func() {
		MakeBoard("[a][a]aaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyyzTooManyLetters")
	})
	assertPanic(t, func() {
		MakeBoard("[a][a]aaaaaaaaabccddeeeeeeeeeehiiijkllmmnooppqruuvvxyyzWrongDistributionOfLettersInThisLongStringOfWords")
	})
	assertPanic(t, func() {
		MakeBoard("[a]aaaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyyz")
	})
	assertPanic(t, func() {
		MakeBoard("[a][a][z]aaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyy")
	})
	assertPanic(t, func() {
		MakeBoard("[a][a]aaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyy0")
	})
}

func BenchmarkScore(b *testing.B) {
	board := NewBoard()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.ScoreWords()
	}
}
