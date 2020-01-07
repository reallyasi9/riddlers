package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
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
		a := avgValue(c, maxValue)
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

func avgValue(coins []int, maxValue int) float64 {
	m := coinChangeMatrix(coins, maxValue)
	return mean(m[len(coins)][1:])
}

var nCoins int
var maxValue int

type coinValues []int

var manualCoinValues coinValues

// String implements flag.Value interface
func (c *coinValues) String() string {
	return fmt.Sprint(*c)
}

// Set implements flag.Value interface
func (c *coinValues) Set(value string) error {
	for _, v := range strings.Split(value, ",") {
		coin, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		*c = append(*c, coin)
	}
	return nil
}

func init() {
	flag.IntVar(&maxValue, "m", 100, "Value for which to make change")
	flag.IntVar(&nCoins, "n", 0, "Number of coins to optimize")
	flag.Var(&manualCoinValues, "v", "Comma-separated values of coins for manual calculation")
}

func main() {

	flag.Parse()
	doManual := false

	if nCoins < 1 && manualCoinValues == nil {
		panic("number of coins cannot be less than 1")
	}

	if manualCoinValues != nil && nCoins != 0 {
		panic("specify only one of flag -n or -v")
	}

	if nCoins == 0 {
		doManual = true
	}

	if doManual {
		smallestV := -1
		duplicateV := make(map[int]bool)
		for _, v := range manualCoinValues {
			if v <= 0 {
				panic("coin values must be greater than 0")
			}
			if _, isSet := duplicateV[v]; isSet {
				panic("coin values must be unique")
			}
			duplicateV[v] = true
			if smallestV < 0 || v < smallestV {
				smallestV = v
			}
		}
		if smallestV != 1 {
			fmt.Println("WARNING: Smallest coin value must be 1.  Value of 1 will be added to coin list automatically.")
			manualCoinValues = append(manualCoinValues, 1)
		}
	} else {
		smallestSum := 0
		for i := 1; i <= nCoins; i++ {
			smallestSum += i
		}
		if maxValue < smallestSum {
			panic("maximum value too small for given number of coins")
		}
	}

	if doManual {
		a := avgValue(manualCoinValues, maxValue)
		fmt.Printf("Average number of coins required for %v to make change for %d:\n%v\n", manualCoinValues, maxValue, a)
	} else {
		best := findBest(nCoins, maxValue)
		fmt.Printf("Best overall %d coins for making change for %d:\n%v\n", nCoins, maxValue, best)
	}

}
