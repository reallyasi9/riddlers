package wordle

import (
	"testing"

	"github.com/kelindar/bitmap"
)

func TestPlayStatus_UpdateWithGuess(t *testing.T) {
	type args struct {
		word Word
		ws   WordStatus
	}

	aaaaa := NewWordFromString("aaaaa")
	abcde := NewWordFromString("abcde")

	expAbsent := bitmap.Bitmap{}
	expAbsent.Set('b' - ZERO_CHAR)
	expAbsent.Set('c' - ZERO_CHAR)
	expAbsent.Set('d' - ZERO_CHAR)
	expAbsent.Set('e' - ZERO_CHAR)

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
				if tt.ps.presentWrongPosition[i].Count() != 1 || tt.ps.presentWrongPosition[i].Contains('a'-ZERO_CHAR) {
					t.Errorf("expected letter %d presentWrongPosition 'a', got %v", i, tt.ps.presentWrongPosition[i])
				}
			}
			if !tt.ps.otherwisePresent.Contains('a' - ZERO_CHAR) {
				t.Errorf("expected 'a' otherwisePresent, got %v", tt.ps.otherwisePresent)
			}
			abs := tt.ps.absent.Clone(nil)
			abs.Xor(expAbsent)
			if abs.Count() != 0 {
				t.Errorf("expected absent to contain only 'bcde', saw %v", tt.ps.absent)
			}
		})
	}
}
