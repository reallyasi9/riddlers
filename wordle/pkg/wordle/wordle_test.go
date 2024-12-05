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

	guesses := make([]Word, 6)
	for i := range guesses {
		guesses[i] = randomWord(s)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wordle.Try(guesses)
	}
}

func TestWordle_Filter(t *testing.T) {
	type fields struct {
		solutions []Word
		status    *PlayStatus
	}
	type args struct {
		guesses []Word
		ws      []WordStatus
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
				guesses: []Word{NewWordFromString("axxxx")},
				ws:      []WordStatus{NewWordStatus("+----")},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wordle{
				solutions: tt.fields.solutions,
			}
			if got := w.Filter(tt.args.guesses, tt.args.ws); got != tt.want {
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
		want1  float64
		want2  float64
		want3  int
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
			want1: 0.0,
			want2: 1.0,
			want3: 2,
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
			want1: 1.0,
			want2: 0.5,
			want3: 0,
		},
		{
			name: "multiple-guesses-perfect-determination",
			fields: fields{
				solutions: []Word{NewWordFromString("aaaaa"), NewWordFromString("bbbbb"), NewWordFromString("aabbb")},
				status:    NewPlayStatus(),
			},
			args: args{
				guesses: []Word{NewWordFromString("xaxxx"), NewWordFromString("xxbxx")},
			},
			want1: 0.0,
			want2: 1.0,
			want3: 3,
		},
		{
			name: "multiple-guesses-partial-determination",
			fields: fields{
				solutions: []Word{NewWordFromString("aaaaa"), NewWordFromString("bbbbb"), NewWordFromString("aabbb")},
				status:    NewPlayStatus(),
			},
			args: args{
				guesses: []Word{NewWordFromString("xaxxx"), NewWordFromString("axxxx")},
			},
			want1: 1.0,
			want2: 2. / 3.,
			want3: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wordle{
				solutions: tt.fields.solutions,
			}
			if got1, got2, got3 := w.Try(tt.args.guesses); got1 != tt.want1 || got2 != tt.want2 || got3 != tt.want3 {
				t.Errorf("Wordle.Try() = %v, %v, %v want %v, %v, %v", got1, got2, got3, tt.want1, tt.want2, tt.want3)
			}
		})
	}
}
