package wordle

import (
	"math/rand"
	"testing"
)

func BenchmarkAmbiguities(b *testing.B) {
	s := rand.NewSource(0x42)
	solns := make([]Word, 1024)
	for i := range solns {
		solns[i] = randomWord(s)
	}
	wordle := NewWordle(solns)

	guesses := make([]Word, b.N)
	for i := range guesses {
		guesses[i] = randomWord(s)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, target := range solns {
			wordle.Ambiguities(guesses, target)
		}
	}
}
