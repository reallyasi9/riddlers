package castler

import (
	"strings"
	"testing"
)

func TestReadDistributions(t *testing.T) {
	in := `Castle 1,Castle 2,Castle 3,Castle 4,Castle 5,Castle 6,Castle 7,Castle 8,Castle 9,Castle 10,Why did you choose your troop deployment?
100,0,0,0,0,0,0,0,0,0,"test"
99.1,0,0,0,0,0,0,0,0,0.9,"error 1"
150,0,0,0,0,0,0,0,0,-50,"error 2"
52,2,2,2,2,2,2,12,12,12,"multi-line

complex ""test"" here."`

	distributions, err := ReadDistributions(strings.NewReader(in))

	if err != nil {
		panic(err)
	}

	if len(distributions) != 2 {
		t.Errorf("expected 2 distributions, saw %d", len(distributions))
	}

	truths := make([][10]int, 2)
	truths[0] = [10]int{100, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	truths[1] = [10]int{52, 2, 2, 2, 2, 2, 2, 12, 12, 12}

	if distributions[0].Total() != 100 {
		t.Errorf("expected 100 soldiers, saw %d", distributions[0].Total())
	} else {
		for n, sd := range distributions {
			for i := 0; i < 10; i++ {
				if sd.Soldiers(i) != truths[n][i] {
					t.Errorf("expected %d soldiers in castle %d, saw %d", truths[n][i], i, sd.Soldiers(i))
				}
			}
		}
	}
}
