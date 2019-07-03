package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sort"

	"./scrabbler"
)

var generations = flag.Int("generations", 1000, "number of generations")
var perGeneration = flag.Int("size", 1000, "number of permutations per generation")
var survivors = flag.Int("survivors", 100, "number of survivors to mutate per generation")
var spawn = flag.Int("spawn", 500, "number of new permutations to spawn each generation")
var seed = flag.Int64("seed", 8675309, "random seed")
var report = flag.Int("report", 1000, "report every n generations")
var temperature = flag.Float64("temperature", 200., "randomness (the higher the temperature, the more random the shuffles)")
var startingBoard = flag.String("start", "", "starting board")

func main() {
	log.Println("Starting")
	flag.Parse()
	rand.Seed(*seed)

	var board *scrabbler.Board
	if *startingBoard == "" {
		board = scrabbler.NewBoard()
	} else {
		board = scrabbler.MakeBoard(*startingBoard)
	}
	gen := scrabbler.MakeGeneration(*perGeneration, board, *temperature)
	for i := 0; i < *generations; i++ {
		gen.Iterate(*survivors, *spawn, *temperature)
		if i%*report == 0 {
			log.Printf("Generation %d: %v\n", i, gen[0])
		}
	}

	sort.Sort(sort.Reverse(gen))
	fmt.Printf("Best permutation found:\n%v\n", gen[0])
	fmt.Println("Words and scores:")
	for word, score := range gen[0].ScoreWords() {
		fmt.Printf("%s: %d\n", word, score)
	}

	log.Println("Done")
}
