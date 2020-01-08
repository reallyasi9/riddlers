package main

import (
	"fmt"
	"sort"
)

var deck hand
var suites []suite
var suiteSymbol []rune
var valueSymbol []rune

func init() {
	deck = make(hand, 52)

	suites = []suite{spades, hearts, diamonds, clubs}

	suiteSymbol = []rune{'\u2660', '\u2661', '\u2662', '\u2663'}

	valueSymbol = []rune{'A', '2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K'}

	for _, suite := range suites {
		for j := 1; j <= 13; j++ {
			c := makeCard(suite, j)
			deck[c.Index] = *c
		}
	}
}

func main() {

	hands := generateHands()
	bests := gimmeNuts(hands)

	for best := range bests {
		fmt.Printf("Best %v\n", best)
	}

}

func generateHands() <-chan hand {
	handChan := make(chan hand)
	go func(ch chan<- hand) {
		defer close(ch)
		for c1 := 0; c1 < 1; c1++ {
			for c2 := c1 + 1; c2 < 2; c2++ {
				for c3 := c2 + 1; c3 < 3; c3++ {
					for c4 := c3 + 1; c4 < 4; c4++ {
						for c5 := c4 + 1; c5 < 52; c5++ {
							ch <- hand{deck[c1],
								deck[c2],
								deck[c3],
								deck[c4],
								deck[c5]}
						}
					}
				}
			}
		}
	}(handChan)
	return handChan
}

func gimmeNuts(hands <-chan hand) <-chan hand {
	bestChan := make(chan hand)
	go func(in <-chan hand, out chan<- hand) {
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

func findStraightFlushNuts(h hand) ([]hand, bool) {
	best := make(hand, 0)
	nuts := make([]hand, 0)
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
