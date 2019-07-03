package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sort"

	"./scrabbler"
)

var generations = flag.Int("steps", 1000, "number of annealing steps")
var nSamples = flag.Int("n", 1000, "number of simulations to run")
var seed = flag.Int64("seed", 8675309, "random seed")
var report = flag.Int("report", 1000, "report every n generations")
var starttemp = flag.Float64("starttemp", 200000., "starting annealing temperature")
var starttemp = flag.Float64("endtemp", 2., "ending annealing temperature")
var startingBoard = flag.String("startboard", "", "starting board (defaults to random based on seed)")

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
