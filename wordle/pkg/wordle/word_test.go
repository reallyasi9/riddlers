package wordle

import (
	"math/rand"
	"reflect"
	"testing"
)

func randomWord(s rand.Source) Word {
	var w Word
	for i := range w {
		w[i] = byte(s.Int63() % N_LETTERS)
	}
	return w
}

func TestWord_Compare(t *testing.T) {
	type args struct {
		soln Word
	}
	tests := []struct {
		name string
		w    Word
		args args
		want WordStatus
	}{
		{
			name: "first-correct-of-many",
			w:    NewWordFromString("axxxx"),
			args: args{
				soln: NewWordFromString("aabbb"),
			},
			want: WordStatus{CORRECT, ABSENT, ABSENT, ABSENT, ABSENT},
		},
		{
			name: "second-correct-of-many",
			w:    NewWordFromString("xaxxx"),
			args: args{
				soln: NewWordFromString("aabbb"),
			},
			want: WordStatus{ABSENT, CORRECT, ABSENT, ABSENT, ABSENT},
		},
		{
			name: "third-present-of-many",
			w:    NewWordFromString("xxaxx"),
			args: args{
				soln: NewWordFromString("aabbb"),
			},
			want: WordStatus{ABSENT, ABSENT, PRESENT, ABSENT, ABSENT},
		},
		{
			name: "two-present-of-many",
			w:    NewWordFromString("xxxaa"),
			args: args{
				soln: NewWordFromString("aabbb"),
			},
			want: WordStatus{ABSENT, ABSENT, ABSENT, PRESENT, PRESENT},
		},
		{
			name: "two-present-of-three",
			w:    NewWordFromString("xxaaa"),
			args: args{
				soln: NewWordFromString("aabbb"),
			},
			want: WordStatus{ABSENT, ABSENT, PRESENT, PRESENT, ABSENT},
		},
		{
			name: "one-correct-one-present-of-three",
			w:    NewWordFromString("xaxaa"),
			args: args{
				soln: NewWordFromString("aabbb"),
			},
			want: WordStatus{ABSENT, CORRECT, ABSENT, PRESENT, ABSENT},
		},
		{
			name: "perfect-guess",
			w:    NewWordFromString("aabbb"),
			args: args{
				soln: NewWordFromString("aabbb"),
			},
			want: WordStatus{CORRECT, CORRECT, CORRECT, CORRECT, CORRECT},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.Compare(tt.args.soln); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Word.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkWord_Compare(b *testing.B) {
	s := rand.NewSource(0x42)
	w := randomWord(s)

	soln := make([]Word, b.N)
	for i := 0; i < b.N; i++ {
		soln[i] = randomWord(s)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Compare(soln[i])
	}
}
