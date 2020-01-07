package main

import (
	"testing"
)

func checkIsEqual(t *testing.T, h1 Hand, h2 Hand, ex bool) {
	if h1.IsEqual(h2) != ex {
		exStr := "!="
		unexStr := "=="
		if ex {
			exStr = "=="
			unexStr = "!="
		}
		t.Errorf("saw %v %s %v, expected %s", h1, unexStr, h2, exStr)
	}
}

func TestIsEqual(t *testing.T) {
	// empty hands
	h1 := Hand{}
	h2 := Hand{}
	checkIsEqual(t, h1, h2, true)

	// same hands
	checkIsEqual(t, h1, h1, true)

	h1 = append(h1, *makeCard(Clubs, 1))
	checkIsEqual(t, h1, h2, false)
	h2 = append(h2, *makeCard(Clubs, 1))
	checkIsEqual(t, h1, h2, true)

	// same hands out of order
	h1 = append(h1, *makeCard(Diamonds, 2), *makeCard(Spades, 3))
	checkIsEqual(t, h1, h2, false)
	h2 = append(h2, *makeCard(Spades, 3), *makeCard(Diamonds, 2))
	checkIsEqual(t, h1, h2, true)

	// entire deck, shuffled

}

func checkSpan(t *testing.T, h Hand, el, eh int) {
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
	// empty Hand
	checkSpan(t, Hand{}, 0, 0)

	// single Card
	for _, c := range StandardDeck {
		h := Hand{c}
		checkSpan(t, h, 0, 0)
	}

	// 2 cards, same Suite
	h := Hand{*makeCard(0, 2), *makeCard(0, 3)}
	checkSpan(t, h, 1, 1)

	// 2 cards, different Suite
	h = Hand{*makeCard(0, 2), *makeCard(1, 3)}
	checkSpan(t, h, 1, 1)

	// many cards, different suites
	h = Hand{*makeCard(0, 2), *makeCard(1, 4), *makeCard(3, 12)}
	checkSpan(t, h, 10, 10)

	// many cards, different suites, aces
	h = Hand{*makeCard(0, 1), *makeCard(1, 4), *makeCard(3, 12)}
	checkSpan(t, h, 11, 10)
}

func checkStraightFlushNuts(t *testing.T, h Hand, n, en []Hand) {
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
	for s := 0; s < 4; s++ {
		// 1 missing on top
		h := Hand{*makeCard(Suite(s), 13), *makeCard(Suite(s), 12), *makeCard(Suite(s), 11), *makeCard(Suite(s), 10)}
		ex := []Hand{Hand{*makeCard(Suite(s), 14)}}
		checkStraightFlushNuts(t, h, h.StraightFlushNuts(), ex)

		// 1 missing on bottom
		h = Hand{*makeCard(Suite(s), 14), *makeCard(Suite(s), 13), *makeCard(Suite(s), 12), *makeCard(Suite(s), 11)}
		ex = []Hand{Hand{*makeCard(Suite(s), 10)}}
		checkStraightFlushNuts(t, h, h.StraightFlushNuts(), ex)
	}
}
