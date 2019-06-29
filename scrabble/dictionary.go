package main

import (
	"bufio"
	"log"
	"net/http"
)

const dictionaryURL = "https://norvig.com/ngrams/enable1.txt"

func buildDictionary(url string) (*ScrabbleTrie, [26][26]float64) {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicf("Error with request: %s", err)
	}
	defer resp.Body.Close()

	trie := &ScrabbleTrie{}
	var bigrams BigramTrie
	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		trie.Insert(s.Text())
		bigrams.Add(s.Text())
	}

	if err := s.Err(); err != nil {
		log.Panicf("Error reading response: %s", err)
	}
	normed := bigrams.Normalize()
	return trie, normed
}
