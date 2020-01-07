package main

import (
	"sort"
)

// Hand is a simple slice of Cards
type Hand []Card

// StandardDeck is a standard American deck of 52 playing cards.
var StandardDeck Hand

func init() {
	StandardDeck = make(Hand, 52)

	for s := 0; s < 4; s++ {
		for j := 1; j <= 13; j++ {
			c := makeCard(Suite(s), j)
			StandardDeck[c.Index] = *c
		}
	}
}

// IsEqual determines if one hand has equal value to another.
func (h Hand) IsEqual(o Hand) bool {
	if len(h) != len(o) {
		return false
	}
	q := make(map[Card]bool)
	for _, c := range h {
		q[c] = true
	}
	for _, c := range o {
		if _, ok := q[c]; !ok {
			return false
		}
	}
	return true
}

// BySuite returns subhands of the given Hand where each subhand consists
// of cards of the same suite.  The subahnds are returned in random order.
func (h Hand) BySuite() <-chan Hand {
	out := make(chan Hand)
	suiteMap := make(map[Suite]Hand)
	for _, c := range h {
		if _, ok := suiteMap[c.Suite]; !ok {
			suiteMap[c.Suite] = Hand{c}
		} else {
			suiteMap[c.Suite] = append(suiteMap[c.Suite], c)
		}
	}
	go func(o chan<- Hand, sm map[Suite]Hand) {
		defer close(o)
		for _, v := range sm {
			o <- v
		}
	}(out, suiteMap)
	return out
}

// BySuite returns subhands of the given Hand where each subhand consists
// of cards of the same value.  The subahnds are returned in random order.
func (h Hand) ByValue() <-chan Hand {
	out := make(chan Hand)
	valueMap := make(map[int]Hand)
	for _, c := range h {
		if _, ok := valueMap[c.Value]; !ok {
			valueMap[c.Value] = Hand{c}
		} else {
			valueMap[c.Value] = append(valueMap[c.Value], c)
		}
	}
	go func(o chan<- Hand, sm map[int]Hand) {
		defer close(o)
		for _, v := range sm {
			o <- v
		}
	}(out, valueMap)
	return out
}

// byAcesLowSuited implements sort.Interface for Hand based on
// the Suite and Value fields (technically the Index field)
type byAcesLowSuited Hand

func (a byAcesLowSuited) Len() int      { return len(a) }
func (a byAcesLowSuited) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byAcesLowSuited) Less(i, j int) bool {
	if a[i].Suite < a[j].Suite {
		return true
	}
	return a[i].AcesLowValue() < a[j].AcesLowValue()
}

// byAcesHighSuited implements sort.Interface for Hand based on
// the Suite and Value fields (technically the Index field)
type byAcesHighSuited Hand

func (a byAcesHighSuited) Len() int      { return len(a) }
func (a byAcesHighSuited) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byAcesHighSuited) Less(i, j int) bool {
	if a[i].Suite < a[j].Suite {
		return true
	}
	return a[i].AcesHighValue() < a[j].AcesHighValue()
}

// byAcesLowUnsuited implements sort.Interface for Hand based on
// the Suite and Value fields (technically the Index field)
type byAcesLowUnsuited Hand

func (a byAcesLowUnsuited) Len() int           { return len(a) }
func (a byAcesLowUnsuited) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAcesLowUnsuited) Less(i, j int) bool { return a[i].AcesLowValue() < a[j].AcesLowValue() }

// byAcesHighUnsuited implements sort.Interface for Hand based on
// the Suite and Value fields (technically the Index field)
type byAcesHighUnsuited Hand

func (a byAcesHighUnsuited) Len() int           { return len(a) }
func (a byAcesHighUnsuited) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAcesHighUnsuited) Less(i, j int) bool { return a[i].AcesHighValue() < a[j].AcesHighValue() }

func (h Hand) CanBeStraight() bool {
	if len(h) < 3 {
		// can't make a straight with 2 addtional cards
		return false
	}
	// sort guaranteed above
	if h[0].Value-h[len(h)-1].Value >= len(h)+2 {
		// can't fill in the gaps with 2 additional cards
		return false
	}
	// we can make a straight flush with at most 2 more cards
	return true
}

func (h Hand) IsStraight() bool {
	sort.Sort(byAcesHighUnsuited(h))
	for i := 0; i < len(h)-1; i++ {
		if h[i+1].Value-h[i].Value != 1 {
			return false
		}
	}
	sort.Sort(byAcesLowUnsuited(h))
	for i := 0; i < len(h)-1; i++ {
		if h[i+1].Value-h[i].Value != 1 {
			return false
		}
	}
	return true
}

func (h Hand) IsFlush() bool {
	if len(h) == 0 {
		return false
	}
	s := h[0].Suite
	for i := 1; i < len(h); i++ {
		if h[i].Suite != s {
			return false
		}
	}
	return true
}

func (h Hand) SpanAcesHigh() int {
	if len(h) < 2 {
		return 0
	}
	sort.Sort(byAcesHighUnsuited(h))
	v1 := h[len(h)-1].Value
	if v1 == 1 {
		v1 = 14
	}
	return v1 - h[0].Value
}

func (h Hand) SpanAcesLow() int {
	if len(h) < 2 {
		return 0
	}
	sort.Sort(byAcesLowUnsuited(h))
	return h[len(h)-1].Value - h[0].Value
}

func (h Hand) StraightFlushNuts() []Hand {
	// sorting in span call
	nHoles1 := h.SpanAcesHigh() - len(h) + 1
	nHoles2 := h.SpanAcesLow() - len(h) + 1
	nHoles := nHoles1
	if nHoles2 < nHoles {
		nHoles = nHoles2
	}

	sort.Sort(sort.Reverse(byAcesHighSuited(h)))

	s := h[0].Suite // special for flushes

	nut := make(Hand, 0)

	// Can add something to the ends to make it better.
	// Start at the top
	// Note:  if there are 2 holes, we fall through
	for i := h[0].AcesHighValue() + 1; i <= 14 && len(nut) < 2-nHoles && len(nut)+len(h) < 5; i++ {
		nut = append(nut, *makeCard(s, i))
	}

	// Fill in the bottom
	// Note:  if there are 2 holes, we fall through
	for i := h[len(h)-1].AcesHighValue() - 1; i > 0 && len(nut) < 2-nHoles && len(nut)+len(h) < 5; i-- {
		nut = append(nut, *makeCard(s, i))
	}

	// Fill in any holes
	for i := 0; i < len(h)-1; i++ {
		for j := 0; j < h[i].Value-h[i+1].Value-1; j++ {
			nut = append(nut, *makeCard(s, h[i].Value-1))
		}
	}

	return []Hand{nut}
}
