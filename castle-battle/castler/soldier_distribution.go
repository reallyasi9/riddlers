package castler

import (
	"fmt"
	"sort"

	"golang.org/x/exp/rand"
)

type SoldierDistribution struct {
	castles [10]int
	total   int
}

func NewSoldierDistribution(soldiers [10]int) *SoldierDistribution {
	total := 0
	for i := 0; i < 10; i++ {
		if soldiers[i] < 0 {
			panic(fmt.Sprintf("soldiers must be non-negative, %d given for castle %d", soldiers[i], i))
		}
		total += soldiers[i]
	}
	return &SoldierDistribution{castles: soldiers, total: total}
}

func DefaultSoldierDistribution(total int) *SoldierDistribution {
	if total < 0 {
		panic(fmt.Sprintf("total soliders must be non-negative, %d given", total))
	}
	remainder := total % 10
	per := (total - remainder) / 10
	var castles [10]int
	for i := 0; i < remainder; i++ {
		castles[i] = per + 1
	}
	for i := remainder; i < 10; i++ {
		castles[i] = per
	}
	return &SoldierDistribution{castles: castles, total: total}
}

func (s *SoldierDistribution) Soldiers(castle int) int {
	return s.castles[castle]
}

func (s *SoldierDistribution) Total() int {
	return s.total
}

// GetPartition gets the location of a partition such that, if you numbered the soldiers from 0 to `Total()`, the return value for partition `part` is the index of the first soldier not in castle `part`.
func (s *SoldierDistribution) GetPartition(part int) int {
	total := 0
	for i := 0; i <= part; i++ {
		total += s.castles[i]
	}
	return total
}

// GetPartitions gets all the partitions at once.
func (s *SoldierDistribution) GetPartitions() (partitions [9]int) {
	partitions[0] = s.castles[0]
	for i := 1; i < len(partitions); i++ {
		partitions[i] = s.castles[i] + partitions[i-1]
	}
	return
}

// FindCastle finds the castle for a given soldier.
func (s *SoldierDistribution) FindCastle(soldier int) int {
	for i, c := range s.castles {
		soldier -= c
		if soldier <= 1 {
			return i
		}
	}
	return len(s.castles)
}

func battle(d1, d2 *SoldierDistribution) (score1, score2 float64) {
	for castle := 0; castle < 10; castle++ {
		if d1.Soldiers(castle) > d2.Soldiers(castle) {
			score1 += float64(castle + 1)
		} else if d1.Soldiers(castle) < d2.Soldiers(castle) {
			score2 += float64(castle + 1)
		} else {
			score1 += float64(castle+1) / 2.
			score2 += float64(castle+1) / 2.
		}
	}
	return
}

func (d *SoldierDistribution) MatchUp(other *SoldierDistribution) int {
	s1, s2 := battle(d, other)
	if s1 > s2 {
		return 1
	} else if s1 < s2 {
		return -1
	} else {
		return 0
	}
}

func jittermod(rng *rand.Rand, x int, sigma float64, mod int) int {
	x += int(rng.NormFloat64() * sigma)
	x %= mod
	for x < 0 {
		x += mod
	}
	return x
}

func (s *SoldierDistribution) partitionSoldiers(partitions [9]int) {
	// Sort partitions to find true boundaries
	sort.Ints(partitions[:])
	// Convert boundaries to soldier distribution
	s.castles[0] = partitions[0]
	s.castles[9] = s.Total() - partitions[8]
	for i := 1; i < 9; i++ {
		s.castles[i] = partitions[i] - partitions[i-1]
	}
}

func (s *SoldierDistribution) Jitter(rng *rand.Rand, strength float64) {
	// Jitter partitions
	partitions := s.GetPartitions()
	for i := 0; i < 9; i++ {
		partitions[i] = jittermod(rng, partitions[i], strength, s.Total())
	}

	s.partitionSoldiers(partitions)
}

func (s *SoldierDistribution) Randomize(rng *rand.Rand) {
	var partitions [9]int
	for i := 0; i < 9; i++ {
		partitions[i] = rng.Intn(s.Total())
	}
	s.partitionSoldiers(partitions)
}

func (s *SoldierDistribution) RandomWalk(rng *rand.Rand, n int) {
	for i := 0; i < n; i++ {
		// pick a soldier
		soldier := rng.Intn(s.Total())
		// find the old castle
		oldC := s.FindCastle(soldier)
		// place in new castle
		newC := rng.Intn(10)
		s.castles[oldC]--
		s.castles[newC]++
	}
}
