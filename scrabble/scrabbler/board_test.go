package scrabbler

import (
	"math"
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

	b := Board{Raw: []byte{'t', 'h', 'e', 'a', 't', 'e', 'r'}, Clean: "theater"}
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

type rawDiff struct {
	Index    int
	Expected byte
	Got      byte
}

func rawDifference(expected, got []byte) []rawDiff {
	diffs := make([]rawDiff, 0)
	for i, r := range expected {
		if len(got) <= i {
			diffs = append(diffs, rawDiff{Index: i, Expected: r, Got: 0x00})
			continue
		}
		if got[i] != r {
			diffs = append(diffs, rawDiff{Index: i, Expected: r, Got: got[i]})
		}
	}
	return diffs
}

func TestMutate(t *testing.T) {
	b1 := NewBoard()
	r1 := make([]byte, len(b1.Raw))
	copy(r1, b1.Raw)
	b1.Nudge(0)

	if diff := rawDifference(r1, b1.Raw); len(diff) != 0 {
		t.Log("zero-temperature mutation resulted in new board:")
		for _, d := range diff {
			t.Logf("Index (%d) Expected '%c' Got '%c'", d.Index, d.Expected, d.Got)
		}
		t.Fail()
	}

	b2 := NewBoard()
	b2.ReplaceWithNudge(b1, math.Inf(1))
	if diff := rawDifference(b1.Raw, b2.Raw); len(diff) == 0 {
		t.Log("infinite-temperature mutation resulted in same board!")
		t.Fail()
	}

	// Calculate minimum temperature for board to be identical
	temperature := 1024.

	for i := 0; i < 100000; i++ {
		b2.ReplaceWithNudge(b1, temperature)
		diff := rawDifference(b1.Raw, b2.Raw)
		if len(diff) == 0 {
			break
		}
		temperature /= 2.
	}

	t.Logf("minimum temperature for score %d: %f", b1.Score, temperature)
}

func BenchmarkMutate(b *testing.B) {
	board := NewBoard()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.Nudge(200.)
	}
}

func BenchmarkScore(b *testing.B) {
	board := NewBoard()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		board.ScoreWords()
	}
}
