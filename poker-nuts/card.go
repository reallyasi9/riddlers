package main

type suite int

const (
	spades   suite = 0
	hearts         = 1
	diamonds       = 2
	clubs          = 3
)

type card struct {
	Suite  suite
	Value  int
	Index  int
	Name   string
	Symbol rune
}

func (c card) AcesHighValue() int {
	if c.Value == 1 {
		return 14
	}
	return c.Value
}

func (c card) AcesLowValue() int {
	return c.Value
}

func (c card) String() string {
	return c.Name
}

func makeCard(s suite, v int) *card {
	if v == 14 {
		v = 1
	}
	index := int(s)*13 + v - 1
	ucard := 0x1F0A1 + int(s)*16 + v - 1
	// skip over the knight of a suite
	if v > 10 {
		ucard++
	}
	return &card{Suite: s,
		Value:  v,
		Index:  index,
		Name:   string(valueSymbol[v-1]) + string(suiteSymbol[s]),
		Symbol: rune(ucard),
	}
}
