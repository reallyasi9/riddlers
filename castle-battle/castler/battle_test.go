package castler

import (
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func TestCastles(t *testing.T) {
	d := DefaultSoldierDistribution(100)
	d.castles[0] = 10
	c1 := d.Castles()

	d.castles[0] = 0
	c2 := d.Castles()

	if c1[0] == c2[0] {
		t.Errorf("expected castles to be different, got %v and %v", c1, c2)
	}
}

func TestBattle(t *testing.T) {
	d1 := NewSoldierDistribution([10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	score1, score2 := battle(d1, NewSoldierDistribution([10]int{0, 3, 4, 5, 6, 7, 8, 9, 10, 11}))
	if score1 != 1.0 {
		t.Errorf("expected score of 1, got %f", score1)
	}
	if score2 != 54.0 {
		t.Errorf("expected score of 54, got %f", score2)
	}

	score1, score2 = battle(d1, NewSoldierDistribution([10]int{1, 3, 4, 5, 6, 7, 8, 9, 10, 11}))
	if score1 != 0.5 {
		t.Errorf("expected score of 1, got %f", score1)
	}
	if score2 != 54.5 {
		t.Errorf("expected score of 54.5, got %f", score2)
	}

	score1, score2 = battle(d1, d1)
	if score1 != 27.5 {
		t.Errorf("expected score of 27.5, got %f", score1)
	}
	if score2 != 27.5 {
		t.Errorf("expected score of 27.5, got %f", score2)
	}
}

func BenchmarkBattle(b *testing.B) {
	sd1 := NewSoldierDistribution([10]int{rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int()})
	sd2s := make([]*SoldierDistribution, b.N)
	for i := 0; i < b.N; i++ {
		sd2s[i] = NewSoldierDistribution([10]int{rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int(), rand.Int()})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		battle(sd1, sd2s[i])
	}
}

func TestJitter(t *testing.T) {
	// soldiers := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	total := 100
	sd := DefaultSoldierDistribution(total)
	rng := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	for i := 0; i < 100; i++ {
		sd.Jitter(rng, 5)

		if sd.Total() != total {
			t.Errorf("expected total of %d, got %d", total, sd.Total())
		}

		for i := 0; i < 10; i++ {
			if sd.Soldiers(i) < 0 {
				t.Fatalf("negative soldiers %d in castle %d", sd.Soldiers(i), i)
			}
		}
	}

}

func BenchmarkJitter(b *testing.B) {
	sd := DefaultSoldierDistribution(100)
	rng := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sd.Jitter(rng, 5)
	}
}

func BenchmarkRandomize(b *testing.B) {
	sd := DefaultSoldierDistribution(100)
	rng := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sd.Randomize(rng)
	}
}

func BenchmarkRandomWalk(b *testing.B) {
	sd := DefaultSoldierDistribution(100)
	rng := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sd.RandomWalk(rng, 1)
	}
}
