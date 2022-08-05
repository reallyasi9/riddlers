package wordle

import (
	"math/rand"
	"reflect"
	"testing"
)

func randomWord(s rand.Source) Word {
	var w Word
	for i := range w {
		w[i] = byte(s.Int63()%ALPHABET_SIZE) + 1
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
			name: "single-single",
			w:    NewWordFromString("abcde"),
			args: args{
				soln: NewWordFromString("acegi"),
			},
			want: WordStatus{CORRECT, ABSENT, PRESENT, ABSENT, PRESENT},
		},
		{
			name: "multiple-single",
			w:    NewWordFromString("aabbc"),
			args: args{
				soln: NewWordFromString("abcde"),
			},
			want: WordStatus{CORRECT, ABSENT, PRESENT, ABSENT, PRESENT},
		},
		{
			name: "single-multiple",
			w:    NewWordFromString("abcde"),
			args: args{
				soln: NewWordFromString("aabbc"),
			},
			want: WordStatus{CORRECT, PRESENT, PRESENT, ABSENT, ABSENT},
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
