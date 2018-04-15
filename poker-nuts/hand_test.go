package main

import (
	"testing"
)

func checkSpan(t *testing.T, h hand, el, eh int) {
	s := h.SpanAcesLow()
	if s != el {
		t.Errorf("%v expected ace-low span %d, saw %d", h, el, s)
	}
	s = h.SpanAcesHigh()
	if s != eh {
		t.Errorf("%v expected ace-high span %d, saw %d", h, eh, s)
	}
}

func TestSpanAcesHigh(t *testing.T) {
	// empty hand
	checkSpan(t, hand{}, 0, 0)

	// single card
	for _, c := range deck {
		h := hand{c}
		checkSpan(t, h, 0, 0)
	}

	// 2 cards, same suite
	h := hand{*makeCard(0, 2), *makeCard(0, 3)}
	checkSpan(t, h, 1, 1)

	// 2 cards, different suite
	h = hand{*makeCard(0, 2), *makeCard(1, 3)}
	checkSpan(t, h, 1, 1)

	// many cards, different suites
	h = hand{*makeCard(0, 2), *makeCard(1, 4), *makeCard(3, 12)}
	checkSpan(t, h, 10, 10)

	// many cards, different suites, aces
	h = hand{*makeCard(0, 1), *makeCard(1, 4), *makeCard(3, 12)}
	checkSpan(t, h, 11, 10)
}

func checkStraightFlushNuts(t *testing.T, h hand, n, en []hand) {
	if len(n) != len(en) {
		t.Errorf("%v should have exactly %d nuts, saw %d", h, len(en), len(n))
	}

	for _, nut := range en {
		found := false
		for _, x := range n {
			if x.IsEqual(nut) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected but did not find nut %v for %v", en, h)
		}
	}

	for _, x := range n {
		found := false
		for _, nut := range en {
			if x.IsEqual(nut) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("calculated nut %v for %v not expected", x, h)
		}
	}
}

func TestStraightFlushNuts(t *testing.T) {
	for _, s := range suites {
		// 1 missing on top
		h := hand{*makeCard(s, 13), *makeCard(s, 12), *makeCard(s, 11), *makeCard(s, 10)}
		ex := []hand{hand{*makeCard(s, 14)}}
		checkStraightFlushNuts(t, h, h.StraightFlushNuts(), ex)

		// 1 missing on bottom
		h = hand{*makeCard(s, 14), *makeCard(s, 13), *makeCard(s, 12), *makeCard(s, 11)}
		ex = []hand{hand{*makeCard(s, 10)}}
		checkStraightFlushNuts(t, h, h.StraightFlushNuts(), ex)
	}
}
