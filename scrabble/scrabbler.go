package main

import (
	"bufio"
	"flag"
	"fmt"
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

var generations = flag.Int("generations", 1000, "number of generations")
var perGeneration = flag.Int("size", 1000, "number of permutations per generation")
var survivors = flag.Int("survivors", 100, "number of survivors to mutate per generation")
var spawn = flag.Int("spawn", 500, "number of new permutations to spawn each generation")
var temperature = flag.Float64("temperature", 100000., "randomness, scaled by score (the larger the temperature, the more random the mutations)")
var seed = flag.Int64("seed", 8675309, "random seed")

func main() {
	log.Println("Starting")
	flag.Parse()
	rand.Seed(*seed)

	err := buildDictionary(dictionaryURL)
	if err != nil {
		log.Fatalln(err)
	}

	gen := MakeGeneration(*perGeneration)
	for i := 0; i < *generations; i++ {
		gen.Iterate(*survivors, *spawn, *temperature)
		log.Printf("Generation %d: %v\n", i, gen[0])
	}

	sort.Sort(sort.Reverse(gen))
	fmt.Printf("Best permutation found:\n%v\n", gen[0])
	fmt.Println("Words and scores:")
	for _, ws := range gen[0].ScoreWords() {
		fmt.Printf("%s: %d\n", ws.Word, ws.Score)
	}

	log.Println("Done")
}
