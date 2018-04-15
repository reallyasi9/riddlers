package main

import (
	"sort"
)

type hand []card

func (h hand) IsEqual(o hand) bool {
	if len(h) != len(o) {
		return false
	}
	q := make(map[card]bool)
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

func (h hand) BySuite() <-chan hand {
	out := make(chan hand)
	suiteMap := make(map[suite]hand)
	for _, c := range h {
		if _, ok := suiteMap[c.Suite]; !ok {
			suiteMap[c.Suite] = hand{c}
		} else {
			suiteMap[c.Suite] = append(suiteMap[c.Suite], c)
		}
	}
	go func(o chan<- hand, sm map[suite]hand) {
		defer close(o)
		for _, v := range sm {
			o <- v
		}
	}(out, suiteMap)
	return out
}

func (h hand) ByValue() <-chan hand {
	out := make(chan hand)
	valueMap := make(map[int]hand)
	for _, c := range h {
		if _, ok := valueMap[c.Value]; !ok {
			valueMap[c.Value] = hand{c}
		} else {
			valueMap[c.Value] = append(valueMap[c.Value], c)
		}
	}
	go func(o chan<- hand, sm map[int]hand) {
		defer close(o)
		for _, v := range sm {
			o <- v
		}
	}(out, valueMap)
	return out
}

// byAcesLowSuited implements sort.Interface for hand based on
// the Suite and Value fields (technically the Index field)
type byAcesLowSuited hand

func (a byAcesLowSuited) Len() int      { return len(a) }
func (a byAcesLowSuited) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byAcesLowSuited) Less(i, j int) bool {
	if a[i].Suite < a[j].Suite {
		return true
	}
	return a[i].AcesLowValue() < a[j].AcesLowValue()
}

// byAcesHighSuited implements sort.Interface for hand based on
// the Suite and Value fields (technically the Index field)
type byAcesHighSuited hand

func (a byAcesHighSuited) Len() int      { return len(a) }
func (a byAcesHighSuited) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byAcesHighSuited) Less(i, j int) bool {
	if a[i].Suite < a[j].Suite {
		return true
	}
	return a[i].AcesHighValue() < a[j].AcesHighValue()
}

// byAcesLowUnsuited implements sort.Interface for hand based on
// the Suite and Value fields (technically the Index field)
type byAcesLowUnsuited hand

func (a byAcesLowUnsuited) Len() int           { return len(a) }
func (a byAcesLowUnsuited) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAcesLowUnsuited) Less(i, j int) bool { return a[i].AcesLowValue() < a[j].AcesLowValue() }

// byAcesHighUnsuited implements sort.Interface for hand based on
// the Suite and Value fields (technically the Index field)
type byAcesHighUnsuited hand

func (a byAcesHighUnsuited) Len() int           { return len(a) }
func (a byAcesHighUnsuited) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAcesHighUnsuited) Less(i, j int) bool { return a[i].AcesHighValue() < a[j].AcesHighValue() }

func (h hand) CanBeStraight() bool {
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

func (h hand) IsStraight() bool {
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

func (h hand) IsFlush() bool {
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

func (h hand) SpanAcesHigh() int {
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

func (h hand) SpanAcesLow() int {
	if len(h) < 2 {
		return 0
	}
	sort.Sort(byAcesLowUnsuited(h))
	return h[len(h)-1].Value - h[0].Value
}

func (h hand) StraightFlushNuts() []hand {
	// sorting in span call
	nHoles1 := h.SpanAcesHigh() - len(h) + 1
	nHoles2 := h.SpanAcesLow() - len(h) + 1
	nHoles := nHoles1
	if nHoles2 < nHoles {
		nHoles = nHoles2
	}

	sort.Sort(sort.Reverse(byAcesHighSuited(h)))

	s := h[0].Suite // special for flushes

	nut := make(hand, 0)

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

	return []hand{nut}
}
