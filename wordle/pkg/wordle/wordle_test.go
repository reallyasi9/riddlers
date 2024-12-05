package wordle

import (
	"math/rand"
	"testing"
)

func BenchmarkTry(b *testing.B) {
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
		wordle.Try(guesses)
	}
}

func TestWordle_AddGuess(t *testing.T) {
	type fields struct {
		solutions []Word
		status    *PlayStatus
	}
	type args struct {
		guess Word
		ws    WordStatus
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "fake-words",
			fields: fields{
				solutions: []Word{NewWordFromString("aabbb"), NewWordFromString("bbbaa")},
				status:    NewPlayStatus(),
			},
			args: args{
				guess: NewWordFromString("axxxx"),
				ws:    NewWordStatus("+----"),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wordle{
				solutions: tt.fields.solutions,
				status:    tt.fields.status,
			}
			if got := w.AddGuess(tt.args.guess, tt.args.ws); got != tt.want {
				t.Errorf("Wordle.AddGuess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWordle_Try(t *testing.T) {
	type fields struct {
		solutions []Word
		status    *PlayStatus
	}
	type args struct {
		guesses []Word
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "perfect-determination",
			fields: fields{
				solutions: []Word{NewWordFromString("aabbb"), NewWordFromString("bbbaa")},
				status:    NewPlayStatus(),
			},
			args: args{
				guesses: []Word{NewWordFromString("axxxx")},
			},
			want: 0,
		},
		{
			name: "no-information",
			fields: fields{
				solutions: []Word{NewWordFromString("aabbb"), NewWordFromString("bbbaa")},
				status:    NewPlayStatus(),
			},
			args: args{
				guesses: []Word{NewWordFromString("xxaxx")},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wordle{
				solutions: tt.fields.solutions,
				status:    tt.fields.status,
			}
			if got := w.Try(tt.args.guesses); got != tt.want {
				t.Errorf("Wordle.Try() = %v, want %v", got, tt.want)
			}
		})
	}
}
