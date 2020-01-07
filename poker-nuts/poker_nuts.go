package main

import (
	"fmt"
	"sort"
)

func main() {

	hands := generateHands()
	bests := gimmeNuts(hands)

	for best := range bests {
		fmt.Printf("Best %v\n", best)
	}

}

func generateHands() <-chan Hand {
	handChan := make(chan Hand)
	go func(ch chan<- Hand) {
		defer close(ch)
		for c1 := 0; c1 < 1; c1++ {
			for c2 := c1 + 1; c2 < 2; c2++ {
				for c3 := c2 + 1; c3 < 3; c3++ {
					for c4 := c3 + 1; c4 < 4; c4++ {
						for c5 := c4 + 1; c5 < 52; c5++ {
							ch <- Hand{StandardDeck[c1],
								StandardDeck[c2],
								StandardDeck[c3],
								StandardDeck[c4],
								StandardDeck[c5]}
						}
					}
				}
			}
		}
	}(handChan)
	return handChan
}

func gimmeNuts(hands <-chan Hand) <-chan Hand {
	bestChan := make(chan Hand)
	go func(in <-chan Hand, out chan<- Hand) {
		defer close(bestChan)
		for h := range in {
			for _, c := range h {
				fmt.Printf("%s ", c.Name)
			}
			fmt.Printf("\n")
			out <- h[3:5]
		}
		fmt.Printf("Ending nuts\n")
	}(hands, bestChan)
	return bestChan
}

func findStraightFlushNuts(h Hand) ([]Hand, bool) {
	best := make(Hand, 0)
	nuts := make([]Hand, 0)
	found := false

	// First look for ace-high
	sort.Sort(sort.Reverse(byAcesHighSuited(h)))

	// Find the first run of at least 3 with no more than 2 gaps
	for suitedHand := range h.BySuite() {
		if !suitedHand.CanBeStraight() {
			continue
		}
		// we can make a straight flush with at most 2 more cards
		if len(best) == 0 || (suitedHand[0].Value > best[0].Value) {
			best = suitedHand
			found = true
		}
	}

	// Do the same with ace-low
	sort.Sort(sort.Reverse(byAcesLowSuited(h)))

	// Find the first run of at least 3 with no more than 2 gaps
	for suitedHand := range h.BySuite() {
		if !suitedHand.CanBeStraight() {
			continue
		}
		// we can make a straight flush with at most 2 more cards
		if len(best) == 0 || (suitedHand[0].Value > best[0].Value) {
			best = suitedHand
			found = true
		}
	}

	if !found {
		return nuts, found
	}

	// if the best is a royal flush, there are no nuts
	if len(best) == 5 && best[0].Value == 14 {
		return nuts, false
	}

	// otherwise, figure out how to fill in the gaps

	return nuts, found
}
