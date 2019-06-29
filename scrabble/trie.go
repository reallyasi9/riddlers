package main

import (
	"log"
	"math/rand"
)

// In rune order
var runeScores = [26]int{1, 3, 3, 2, 1, 4, 2, 4, 1, 8, 5, 1, 3, 1, 1, 3, 10, 1, 1, 1, 1, 4, 4, 8, 4, 10}

// ScrabbleTrie represents a scrabble word list with associated word values.
type ScrabbleTrie struct {
	children [26]*ScrabbleTrie
	Score    int
}

// Insert a word into the trie in O(n) time.
// Word score is calculated automatically and inserted.
// Returns true if a new child was inserted, false otherwise.
func (s *ScrabbleTrie) Insert(key string) {
	s.recursiveInsert([]byte(key), 0)
}

// idx is the index to insert (rune - 'a'), remainder is what remains, and value is the value up to this node (not including idx)
func (s *ScrabbleTrie) recursiveInsert(remainder []byte, value int) {

	if len(remainder) == 0 {
		s.Score = value
		return
	}

	idx := remainder[0] - 'a'
	if idx >= 26 {
		log.Panicf("letter '%c' not a scrabble tile", idx)
	}
	value += runeScores[idx]
	child := s.children[idx]

	if child == nil {
		child = &ScrabbleTrie{}
		s.children[idx] = child
	}

	child.recursiveInsert(remainder[1:], value)
}

// Get a node representing a prefix in the trie in O(n) time.
// Returns nil if no such prefix exists in the trie.
func (s *ScrabbleTrie) Get(prefix string) *ScrabbleTrie {
	return s.recursiveGet([]byte(prefix))
}

func (s *ScrabbleTrie) recursiveGet(prefix []byte) *ScrabbleTrie {
	if len(prefix) == 0 {
		return s
	}

	child := s.Step(prefix[0])

	if child == nil {
		return nil
	}

	return child.recursiveGet(prefix[1:])
}

// Step gets the next branch in the trie if it exists, otherwise returns nil.
func (s *ScrabbleTrie) Step(b byte) *ScrabbleTrie {
	idx := b - 'a'
	if idx >= 26 {
		log.Panicf("letter '%c' not a scrabble tile", idx)
	}
	return s.children[idx]
}

// BigramTrie represents a simple table of bigrams for fast frequency lookups
type BigramTrie [26][26]int

// NormedBigramTrie is a normalized version of BigramTrie
type NormedBigramTrie [26][26]float64

// Add a word's bigrams to the trie
func (t *BigramTrie) Add(word string) {
	for i := 0; i < len(word)-1; i++ {
		t.Increment(word[i], word[i+1])
	}
}

// Increment adds to the bigram count for bigram a,b
func (t *BigramTrie) Increment(a, b byte) {
	(*t)[a-'a'][b-'a']++
}

// Normalize makes all the values proper CDFs by first element.  Adds 1 to zeros to make sure those have some probability.
func (t BigramTrie) Normalize() NormedBigramTrie {
	var out NormedBigramTrie
	for i := 0; i < 26; i++ {
		sum := 0.
		for j := 0; j < 26; j++ {
			if t[i][j] == 0 {
				t[i][j] = 1
			}
			sum += float64(t[i][j])
		}
		out[i][0] = float64(t[i][0]) / sum
		for j := 1; j < 26; j++ {
			out[i][j] = out[i][j-1] + float64(t[i][j])/sum
		}
	}
	return out
}

// Random picks a random next letter in an MCMC fashion.
func (n NormedBigramTrie) Random(a byte, rng *rand.Rand) byte {
	p := rng.Float64()
	for i, c := range n[a] {
		if c > p {
			return byte(i) + 'a'
		}
	}
	return 'z' // just in case
}
