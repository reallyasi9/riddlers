package main

import (
	"bufio"
	"log"
	"math/rand"
	"net/http"
	"sort"
)

const dictionaryURL = "https://norvig.com/ngrams/enable1.txt"

func buildDictionary(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error with request: %s", err)
		return err
	}
	defer resp.Body.Close()

	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		word := s.Text()
		score, err := scoreWord(word)
		if err != nil {
			log.Printf("Error scoring word: %s", err)
			return err
		}
		scoreTrie.Insert(word, score)
	}

	if err := s.Err(); err != nil {
		log.Printf("Error reading response: %s", err)
		return err
	}
	return nil
}

type prefixWalker struct {
	isPrefix bool
}

func (p *prefixWalker) walk(s string, v interface{}) bool {
	p.isPrefix = true
	return false // stop iterating immediately
}

func mutateBoard(board []rune) {
	i := rand.Intn(len(board))
	j := rand.Intn(len(board))
	board[i], board[j] = board[j], board[i]
}

func generateOffspring(p1 []rune, p2 []rune) []rune {
	child := make([]rune, len(p1))
	for i := 0; i < len(child); i++ {
		if rand.Float32() < .5 {
			child[i] = p1[i]
		} else {
			child[i] = p2[i]
		}
	}
	return child
}

const generations = 1000
const perGeneration = 1000

func main() {
	log.Println("Starting")

	err := buildDictionary(dictionaryURL)
	if err != nil {
		log.Fatalln(err)
	}

	gen := MakeGeneration(perGeneration)

	sort.Sort(sort.Reverse(gen))
	for i := 0; i < 10; i++ {
		println(gen[i].String())
	}

	log.Println("Done")
}
