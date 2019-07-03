package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"

	"./scrabbler"
)

var generations = flag.Int("steps", 1000000, "maximum number of annealing steps")
var nSamples = flag.Int("n", 50, "number of simulations to run")
var seed = flag.Int64("seed", 8675309, "random seed")
var report = flag.Int("report", 0, "report every n generations (0 turns off reporting)")
var startTemp = flag.Float64("starttemp", 200., "starting annealing temperature")
var tempStep = flag.Float64("logstep", 0.01, "annealing stepping temperature (log scale)")
var patience = flag.Int("patience", 10000, "number of steps of no change before ending annealing")

type byScore []scrabbler.Generation

func (a byScore) Len() int           { return len(a) }
func (a byScore) Less(i, j int) bool { return a[i].Board.Score < a[j].Board.Score }
func (a byScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
	log.Println("Starting")
	flag.Parse()
	rand.Seed(*seed)

	if flag.NArg() > 0 {
		log.Println("Evaluating boards:")
		for _, b := range flag.Args() {
			board := scrabbler.MakeBoard(b)
			fmt.Printf("%s\n", *board)
			for word, score := range board.ScoreWords() {
				fmt.Printf("%s: %d\n", word, score)
			}
		}
		log.Println("Done")
		return
	}

	gs := make(byScore, *nSamples)
	var wg sync.WaitGroup
	for i := 0; i < *nSamples; i++ {
		wg.Add(1)
		rng := rand.New(rand.NewSource(rand.Int63()))
		gs[i] = *scrabbler.NewGeneration(*generations, *startTemp, *tempStep, *report, *patience, rng)
		go func(i int) {
			gs[i].Anneal()
			wg.Done()
		}(i)
	}
	wg.Wait()

	sort.Sort(sort.Reverse(gs))

	nth := 9
	if nth >= *nSamples {
		nth = *nSamples
	}
	fmt.Printf("Top %d boards:\n", nth+1)
	for ; nth >= 0; nth-- {
		fmt.Printf("%2d: %s\n", nth+1, *gs[nth].Board)
	}
	fmt.Println("Words from top board:")
	for word, score := range gs[0].Board.ScoreWords() {
		fmt.Printf("%s: %d\n", word, score)
	}

	log.Println("Done")
}
