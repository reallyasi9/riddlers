package scrabbler

import (
	"bufio"
	"log"
	"net/http"
)

const dictionaryURL = "https://norvig.com/ngrams/enable1.txt"

// ScoreTrie is a trie of legal Scrabble words and their scores.
var ScoreTrie ScrabbleTrie

// BigramTrie is a lookup table of bigram counts in legal Scrabble words
var BigramTrie ScrabbleBigrams

// BuildDictionary fills the score and bigram tries.
func init() {
	resp, err := http.Get(dictionaryURL)
	if err != nil {
		log.Panicf("Error with request: %s", err)
	}
	defer resp.Body.Close()

	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		ScoreTrie.Insert(s.Text())
		BigramTrie.Add(s.Text())
	}

	if err := s.Err(); err != nil {
		log.Panicf("Error reading response: %s", err)
	}
}
