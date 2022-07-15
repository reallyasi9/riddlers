package wordle

import (
	"math/rand"
	"testing"

	"github.com/kelindar/bitmap"
)

func randomWord(s rand.Source) Word {
	var w Word
	for i := range w {
		w[i] = byte(s.Int63()%ALPHABET_SIZE) + 1
	}
	return w
}

func TestPlayStatus_UpdateWithGuess(t *testing.T) {
	type args struct {
		word Word
		ws   WordStatus
	}

	aaaaa := NewWordFromString("aaaaa")
	abcde := NewWordFromString("abcde")

	expAbsent := bitmap.Bitmap{}

	tests := []struct {
		name string
		ps   *PlayStatus
		args args
	}{
		// TODO: Add test cases.
		{
			name: "'aaaaa' vs 'abcde' from (empty)",
			ps:   NewStatus(),
			args: args{
				word: aaaaa,
				ws:   aaaaa.Compare(abcde),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ps.UpdateWithGuess(tt.args.word, tt.args.ws)
			if tt.ps.absent.Count() != 0 {
				t.Errorf("expected no absent letters, got %v", tt.ps.absent)
			}
			if tt.ps.presentWrongPosition[0].Count() != 0 {
				t.Errorf("expected first letter not presentWrongPosition, got %v", tt.ps.presentWrongPosition[0])
			}
			for i := 1; i < WORD_SIZE; i++ {
				if tt.ps.presentWrongPosition[i].Count() != 1 || !tt.ps.presentWrongPosition[i].Contains('a'-ZERO_CHAR) {
					t.Errorf("expected letter %d presentWrongPosition 'a', got %v", i, tt.ps.presentWrongPosition[i])
				}
			}
			if !tt.ps.otherwisePresent.Contains('a' - ZERO_CHAR) {
				t.Errorf("expected 'a' otherwisePresent, got %v", tt.ps.otherwisePresent)
			}
			abs := tt.ps.absent.Clone(nil)
			abs.Xor(expAbsent)
			if abs.Count() != 0 {
				t.Errorf("expected absent be empty, saw %v", tt.ps.absent)
			}
		})
	}
}

func BenchmarkPlayStatus_UpdateWithGuess(b *testing.B) {
	ps := NewStatus()
	s := rand.NewSource(0x42)
	soln := randomWord(s)

	guesses := make([]Word, b.N)
	statuses := make([]WordStatus, b.N)
	for i := 0; i < b.N; i++ {
		word := randomWord(s)
		guesses[i] = word
		statuses[i] = word.Compare(soln)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.UpdateWithGuess(guesses[i], statuses[i])
	}
}
