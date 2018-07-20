package main

import (
	"fmt"
	"os"
	"strconv"
)

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func coinChangeMatrix(coins []int, maxValue int) [][]int {
	m := make([][]int, len(coins)+1)
	for ic := 0; ic <= len(coins); ic++ {
		m[ic] = make([]int, maxValue+1)
	}
	for iv := 0; iv < maxValue+1; iv++ {
		m[0][iv] = iv
	}

	for ic := 1; ic <= len(coins); ic++ {
		c := coins[ic-1]
		for iv := 1; iv <= maxValue; iv++ {
			if c == iv {
				m[ic][iv] = 1
			} else if c > iv {
				m[ic][iv] = m[ic-1][iv]
			} else {
				m[ic][iv] = min(m[ic-1][iv], 1+m[ic][iv-c])
			}
		}
	}
	return m
}

func mean(x []int) float64 {
	t := 0.
	for _, val := range x {
		t += float64(val)
	}
	return t / float64(len(x))
}

type bestCoins struct {
	avgCoins float64
	coins    []int
}

func findBest(nCoins, maxValue int) []bestCoins {
	coins := combos(maxValue, nCoins)

	best := bestCoins{avgCoins: float64(maxValue), coins: make([]int, nCoins)}
	out := make([]bestCoins, 0)
	for c := range coins {
		m := coinChangeMatrix(c, maxValue)
		a := mean(m[nCoins][1:])
		if a <= best.avgCoins {
			newBest := bestCoins{avgCoins: a, coins: c}
			if a < best.avgCoins {
				out = out[:0]
			}
			out = append(out, newBest)
			best = newBest
			fmt.Printf("Best coins found so far: %v (average of %0.3f coins per change)\n", c, a)
		}
	}

	return out
}

func main() {

	nCoins, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("unable to parse number of coins")
	}

	maxValue, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("unable to parse maximum value")
	}

	if nCoins < 1 {
		panic("number of coins cannot be less than 1")
	}

	smallestSum := 0
	for i := 1; i <= nCoins; i++ {
		smallestSum += i
	}
	if maxValue < smallestSum {
		panic("maximum value too small for given number of coins")
	}

	best := findBest(nCoins, maxValue)
	fmt.Printf("Best overall:\n%v", best)

}
