package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/rand"

	"./castler"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fns := []string{"castle-solutions.csv", "castle-solutions-2.csv", "castle-solutions-3.csv"}
	sds := make([]*castler.SoldierDistribution, 0)
	for _, fn := range fns {
		file, err := os.Open(fn)
		check(err)

		fmt.Printf("reading file %s\n", fn)
		sd, err := castler.ReadDistributions(file)
		check(err)

		sds = append(sds, sd...)
	}

	rng := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	test := castler.DefaultSoldierDistribution(100)

	var best [10]int
	var bestRecord [3]int
	var bestScore int
	for i := 0; i < 10000; i++ {
		test.Jitter(rng, 10)
		var record [3]int
		// Bootstrap
		for j := 0; j < len(sds); j++ {
			record[test.MatchUp(sds[rng.Intn(len(sds))])+1]++
		}
		score := record[2]
		if record[0] > bestScore {
			bestScore = score
			bestRecord = record
			for k := 0; k < 10; k++ {
				best[k] = test.Soldiers(k)
			}
		}
		// for l := 0; l < 10; l++ {
		// 	fmt.Printf("%d,", test.Soldiers(l))
		// }
		// for m := 0; m < 3; m++ {
		// 	fmt.Printf("%d,", record[m])
		// }
		// fmt.Printf("%d\n", score)
	}

	fmt.Println(bestRecord)
	fmt.Println(bestScore)
	fmt.Println(best)
}
