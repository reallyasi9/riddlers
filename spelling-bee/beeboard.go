package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var flagWordList string
var flagMin bool

func init() {
	flag.StringVar(&flagWordList, "w", "", "path to word list")
	flag.BoolVar(&flagMin, "m", false, "find least-valuable boards")
}

func main() {
	flag.Parse()

	file, err := os.Open(flagWordList)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	words := make([]string, 0)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) < 4 {
			continue
		}
		words = append(words, word)
	}

	boards, err := NewBoards("abcdefghijklmnopqrtuvwxyz")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("found %d boards", len(boards))

	bestScore := 0
	if flagMin {
		bestScore = 10000000 // bignum!
	}
	var bestBoard *Board
	for i := range boards {
		if i%10000 == 0 {
			log.Printf("processed %d boards so far...", i)
		}
		score := 0
		pangrams := false
		for _, word := range words {
			s, p := boards[i].Score(word)
			score += s
			pangrams = pangrams || p
		}
		if !pangrams {
			// need at least one pangram
			continue
		}
		if (!flagMin && score > bestScore) || (flagMin && score <= bestScore) {
			bestScore = score
			bestBoard = boards[i]
			fmt.Printf("%d: %v\n", bestScore, *bestBoard)
			if flagMin {
				for _, word := range words {
					s, _ := bestBoard.Score(word)
					if s > 0 {
						fmt.Printf("%s = %d\n", word, s)
					}
				}
			}
		}
	}

	for _, word := range words {
		score, _ := bestBoard.Score(word)
		if score > 0 {
			fmt.Printf("%s = %d\n", word, score)
		}
	}

	log.Print("done")
}
