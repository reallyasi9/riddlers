package main

import (
	"math/rand"
	"os"
	"testing"
)

func BenchmarkDisjointLetters(b *testing.B) {
	guesses := []Word{
		makeWord([]byte("hello")),
		makeWord([]byte("smato")),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		disjointLetters(guesses)
	}
}

func BenchmarkCompare(b *testing.B) {
	word := makeWord([]byte("hello"))
	other := makeWord([]byte("shlep"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word.Compare(other)
	}
}

func randomWord() Word {
	var w Word
	for i := 0; i < 5; i++ {
		w[i] = byte(rand.Intn(26))
	}
	return w
}

func BenchmarkAmbiguities(b *testing.B) {
	solnFile, err := os.Open("solutions.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer solnFile.Close()
	solns := readWords(solnFile)
	wordle := NewWordle(solns)

	guesses := []Word{
		makeWord([]byte("hello")),
		makeWord([]byte("shaps")),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, target := range solns {
			wordle.Ambiguities(guesses, target)
		}
	}
}
