package wordle

import (
	"math/rand"
	"testing"

	"github.com/kelindar/bitmap"
)

func TestPlayStatus_UpdateWithGuess(t *testing.T) {
	type args struct {
		target Word
		words  []Word
		wss    []WordStatus
	}

	ceorl := NewWordFromString("ceorl")
	saint := NewWordFromString("saint")
	crazy := NewWordFromString("crazy")
	cramp := NewWordFromString("cramp")

	tests := []struct {
		name   string
		ps     *PlayStatus
		args   args
		exppos [][WORD_SIZE]bitmap.Bitmap
		expmin [][N_LETTERS]int
		expmax [][N_LETTERS]int
	}{
		// TODO: Add test cases.
		{
			name: "backtrack-cramp",
			ps:   NewPlayStatus(),
			args: args{
				target: cramp,
				words:  []Word{ceorl, saint, crazy},
				wss: []WordStatus{
					{CORRECT, ABSENT, ABSENT, PRESENT, ABSENT},
					{ABSENT, PRESENT, ABSENT, ABSENT, ABSENT},
					{CORRECT, CORRECT, CORRECT, ABSENT, ABSENT},
				},
			},
			exppos: [][WORD_SIZE]bitmap.Bitmap{
				{
					//              zyxwvutsrqponmlkjihgfedcba.
					bitmap.Bitmap{1 << 3},                        // Only c
					bitmap.Bitmap{0b111111111110110111111011111}, // Not: elo
					bitmap.Bitmap{0b111111111110110111111011111}, // Not: elo
					bitmap.Bitmap{0b111111110110110111111011111}, // Not: elor
					bitmap.Bitmap{0b111111111110110111111011111}, // Not: elo
				},
				{
					//              zyxwvutsrqponmlkjihgfedcba.
					bitmap.Bitmap{1 << 3},                        // Only c
					bitmap.Bitmap{0b111111001110010110111011101}, // Not: aeilnost
					bitmap.Bitmap{0b111111001110010110111011111}, // Not: eilnost
					bitmap.Bitmap{0b111111000110010110111011111}, // Not: eilnorst
					bitmap.Bitmap{0b111111001110010110111011111}, // Not: eilnost
				},
				{
					//              zyxwvutsrqponmlkjihgfedcba.
					bitmap.Bitmap{1 << 3},                        // Only c
					bitmap.Bitmap{1 << 18},                       // Only r
					bitmap.Bitmap{1 << 1},                        // Only a
					bitmap.Bitmap{0b001111000110010110111011111}, // Not: eilnorstyz
					bitmap.Bitmap{0b001111001110010110111011111}, // Not: eilnostyz
				},
			},
			expmin: [][N_LETTERS]int{
				//  a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z
				{
					0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0,
				},
				{
					1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0,
				},
				{
					1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0,
				},
			},
			expmax: [][N_LETTERS]int{
				//   a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s, t, u, v, w, x, y, z
				{
					-1, -1, -1, -1, 0, -1, -1, -1, -1, -1, -1, 0, -1, -1, 0, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
				},
				{
					-1, -1, -1, -1, 0, -1, -1, -1, 0, -1, -1, 0, -1, 0, 0, -1, -1, -1, 0, 0, -1, -1, -1, -1, -1, -1,
				},
				{
					-1, -1, -1, -1, 0, -1, -1, -1, 0, -1, -1, 0, -1, 0, 0, -1, -1, -1, 0, 0, -1, -1, -1, -1, 0, 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, guess := range tt.args.words {
				tt.ps.UpdateWithGuess(guess, tt.args.wss[i])
				for letter := 0; letter < WORD_SIZE; letter++ {
					var testBitmap bitmap.Bitmap
					testBitmap.Or(tt.ps.possible[letter])
					testBitmap.Xor(tt.exppos[i][letter])
					if testBitmap.Count() != 0 {
						t.Errorf("after guessing %s against %s (%v), expected letter %d present %v, got %v", guess, tt.args.target, tt.args.wss[i], letter, tt.exppos[i][letter], tt.ps.possible[letter])
					}
				}
				if tt.ps.minimumPresent != tt.expmin[i] {
					t.Errorf("after guessing %s against %s (%v), expected minimumPresent %v, got %v", guess, tt.args.target, tt.args.wss[i], tt.expmin[i], tt.ps.minimumPresent)
				}
				if tt.ps.maximumPresent != tt.expmax[i] {
					t.Errorf("after guessing %s against %s (%v), expected maximumPresent %v, got %v", guess, tt.args.target, tt.args.wss[i], tt.expmax[i], tt.ps.maximumPresent)
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
