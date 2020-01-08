package main

// Suite represents an American playing card suit.
type Suite int

// The four standard American playing card suites.
const (
	Spades   Suite = 0
	Hearts         = 1
	Diamonds       = 2
	Clubs          = 3
)

// Card is a struct representing a single standard American playing card.
type Card struct {
	Suite  Suite
	Value  int
	Index  int
	Name   string
	Symbol rune
}

var suiteSymbol []rune
var valueSymbol []rune

func init() {
	suiteSymbol = []rune{'\u2660', '\u2661', '\u2662', '\u2663'}
	valueSymbol = []rune{'A', '2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K'}
}

// AcesHighValue returns a comparible value of the card assuming the ace is the
// most valuable card in a suite.
func (c Card) AcesHighValue() int {
	if c.Value == 1 {
		return 14
	}
	return c.Value
}

// AcesLowValue returns a comparible value of the card assuming the ace is the
// least valuable card in the suite.
func (c Card) AcesLowValue() int {
	return c.Value
}

func (c Card) String() string {
	return c.Name
}

func makeCard(s Suite, v int) *Card {
	if v == 14 {
		v = 1
	}
	index := int(s)*13 + v - 1
	ucard := 0x1F0A1 + int(s)*16 + v - 1
	// skip over the knight of a Suite
	if v > 10 {
		ucard++
	}
	return &Card{Suite: s,
		Value:  v,
		Index:  index,
		Name:   string(valueSymbol[v-1]) + string(suiteSymbol[s]),
		Symbol: rune(ucard),
	}
}
