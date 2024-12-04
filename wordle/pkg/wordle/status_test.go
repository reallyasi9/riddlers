package wordle

import (
	"math/rand"
	"testing"

	"github.com/kelindar/bitmap"
)

func TestPlayStatus_UpdateWithGuess_Possible(t *testing.T) {
	type args struct {
		target   Word
		expected []bool
		words    []Word
		wss      []WordStatus
	}

	axxxx := NewWordFromString("axxxx") // "+----"
	xaxxx := NewWordFromString("xaxxx") // "-+---"
	xxaxx := NewWordFromString("xxaxx") // "--?--"
	xxxaa := NewWordFromString("xxxaa") // "---??"
	xxaaa := NewWordFromString("xxaaa") // "--??-" (tricky)
	xaxaa := NewWordFromString("xaxaa") // "-+-?-" (tricky)
	aabbb := NewWordFromString("aabbb") // "+++++" (target)
	qqqqq := NewWordFromString("qqqqq") // "?????" (not possible, no more possible solutions)

	tests := []struct {
		name   string
		ps     *PlayStatus
		args   args
		exppos [][WORD_SIZE]bitmap.Bitmap
		expmin [][N_LETTERS]int
		expmax [][N_LETTERS]int
	}{
		{
			name: "guess-chain",
			ps:   NewPlayStatus(),
			args: args{
				target:   aabbb,
				expected: []bool{true, true, true, true, true, true, true, false},
				words:    []Word{axxxx, xaxxx, xxaxx, xxxaa, xxaaa, xaxaa, aabbb, qqqqq},
				wss: []WordStatus{
					{CORRECT, ABSENT, ABSENT, ABSENT, ABSENT},
					{ABSENT, CORRECT, ABSENT, ABSENT, ABSENT},
					{ABSENT, ABSENT, PRESENT, ABSENT, ABSENT},
					{ABSENT, ABSENT, ABSENT, PRESENT, PRESENT},
					{ABSENT, ABSENT, PRESENT, PRESENT, ABSENT},
					{ABSENT, CORRECT, ABSENT, PRESENT, ABSENT},
					{CORRECT, CORRECT, CORRECT, CORRECT, CORRECT},
					{PRESENT, PRESENT, PRESENT, PRESENT, PRESENT},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, guess := range tt.args.words {
				tt.ps.UpdateWithGuess(guess, tt.args.wss[i])
				pos := tt.ps.Possible(tt.args.target)
				if pos != tt.args.expected[i] {
					t.Errorf("after guessing %v (%v), expected %s possible %t, got %t", guess, tt.args.wss[i], tt.args.target, tt.args.expected[i], pos)
				}
			}
		})
	}
}

func BenchmarkPlayStatus_UpdateWithGuess(b *testing.B) {
	ps := NewPlayStatus()
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

func TestWordStatus_TrinaryIndex(t *testing.T) {
	tests := []struct {
		name string
		ws   WordStatus
		want int
	}{
		{
			name: "all-absent",
			ws:   NewWordStatus("-----"),
			want: 0,
		},
		{
			name: "lsb-present",
			ws:   NewWordStatus("?----"),
			want: 1,
		},
		{
			name: "msb-present",
			ws:   NewWordStatus("----?"),
			want: 81,
		},
		{
			name: "all-present",
			ws:   NewWordStatus("?????"),
			want: 1 + 3 + 9 + 27 + 81,
		},
		{
			name: "lsb-correct",
			ws:   NewWordStatus("+----"),
			want: 2,
		},
		{
			name: "msb-correct",
			ws:   NewWordStatus("----+"),
			want: 162,
		},
		{
			name: "all-correct",
			ws:   NewWordStatus("+++++"),
			want: N_STATUS_OUTCOMES - 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ws.TrinaryIndex(); got != tt.want {
				t.Errorf("WordStatus (%v) TrinaryIndex() = %v, want %v", tt.ws, got, tt.want)
			}
		})
	}
}

func BenchmarkWordStatus_TrinaryIndex(b *testing.B) {
	s := rand.NewSource(0x42)
	status := randomWord(s).Compare(randomWord(s))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		status.TrinaryIndex()
	}
}
