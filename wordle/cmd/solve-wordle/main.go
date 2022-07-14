package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"

	"github.com/reallyasi9/riddler/wordle/pkg/wordle"
)

func main() {

	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatal("requres two positional arguments: a file containing a list of possible solutions and a file containing a list of possible guesses")
	}

	solnFile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	solutions := readWords(solnFile)
	solnFile.Close()

	guessFile, err := os.Open(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	guessables := readWords(guessFile)
	guessFile.Close()

	log.Println("Minimizing initial entropy (this may take a few minutes)")

	for _, guess := range guessables {
		representativeSolution := make(map[uint64]*wordle.PlayStatus)
		possibleSolutionGroups := make(map[uint64][]wordle.Word)
		for _, soln := range solutions {
			ws := guess.Compare(soln)
			stat := wordle.NewStatus()
			stat.UpdateWithGuess(guess, ws)
			hash := stat.Hash()
			representativeSolution[hash] = stat
			if grp, ok := possibleSolutionGroups[hash]; ok {
				grp = append(grp, soln)
				possibleSolutionGroups[hash] = grp
			} else {
				possibleSolutionGroups[hash] = []wordle.Word{soln}
			}
		}

	}
}

func readWords(r io.Reader) []wordle.Word {
	words := make([]wordle.Word, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		words = append(words, wordle.NewWord(scanner.Bytes()))
	}
	return words
}
