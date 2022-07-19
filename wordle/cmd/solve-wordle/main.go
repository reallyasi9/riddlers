package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/reallyasi9/riddler/wordle/pkg/wordle"
)

func main() {

	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("requres at two positional arguments: a file containing a list of possible solutions and a file containing a list of possible guesses")
	}
	if flag.NArg()%2 != 0 {
		log.Fatal("additional positional arguments must be word-status pairs")
	}

	solnFile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	initialSolutions := readWords(solnFile)
	solnFile.Close()

	guessFile, err := os.Open(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	guessables := readWords(guessFile)
	guessFile.Close()

	startingStatus := wordle.NewPlayStatus()
	for iarg := 2; iarg < flag.NArg(); iarg += 2 {
		guess := flag.Arg(iarg)
		if len(guess) != 5 {
			log.Fatal("guesses can only be 5 letters")
		}
		word := wordle.NewWordFromString(strings.ToLower(guess))
		stat := wordle.NewWordStatus(flag.Arg(iarg + 1))
		startingStatus.UpdateWithGuess(word, stat)
	}

	solutions := make(map[wordle.Word]struct{})
	for _, soln := range initialSolutions {
		if !startingStatus.Possible(soln) {
			continue
		}
		solutions[soln] = struct{}{}
	}

	if len(solutions) == 1 {
		fmt.Print("There is only one possible soution remaining: ")
		for soln := range solutions {
			fmt.Print("%s\n", soln)
		}
		return
	}

	if len(solutions) <= 10 {
		fmt.Printf("There are only %d solutions remaining:\n", len(solutions))
		for soln := range solutions {
			fmt.Printf("%s\n", soln)
		}
	}
	if len(solutions) == 2 {
		return
	}

	log.Println("Minimizing initial entropy (this may take a few minutes)")

	entropy := make([]EntropyWord, len(guessables))
	nsolns := float64(len(solutions))
	for iguess, guess := range guessables {
		possibleSolutionGroups := make(map[uint64][]wordle.Word)
		for soln := range solutions {
			ws := guess.Compare(soln)
			stat := startingStatus.Clone()
			stat.UpdateWithGuess(guess, ws)
			hash := stat.Hash()
			if grp, ok := possibleSolutionGroups[hash]; ok {
				grp = append(grp, soln)
				possibleSolutionGroups[hash] = grp
			} else {
				possibleSolutionGroups[hash] = []wordle.Word{soln}
			}
		}
		var eta float64
		for _, grp := range possibleSolutionGroups {
			l := float64(len(grp)) / nsolns
			eta += l * math.Log2(l)
		}
		_, isSoln := solutions[guess]
		entropy[iguess] = EntropyWord{Entropy: eta, Word: guess, IsSolution: isSoln, SolutionGroups: len(possibleSolutionGroups)}
	}

	sort.Sort(ByEntropy(entropy))
	fmt.Printf("Best guess: %s (entropy: %f, solutions remaining: %d)\n", entropy[0].Word, entropy[0].Entropy, entropy[0].SolutionGroups)

}

type EntropyWord struct {
	Entropy        float64
	Word           wordle.Word
	IsSolution     bool
	SolutionGroups int
}

type ByEntropy []EntropyWord

func (a ByEntropy) Len() int {
	return len(a)
}

func (a ByEntropy) Swap(x, y int) {
	a[x], a[y] = a[y], a[x]
}

func (a ByEntropy) Less(x, y int) bool {
	if a[x].Entropy != a[y].Entropy {
		return a[x].Entropy < a[y].Entropy
	}
	if a[x].IsSolution != a[y].IsSolution {
		return a[x].IsSolution
	}
	if a[x].SolutionGroups != a[y].SolutionGroups {
		return a[x].SolutionGroups > a[y].SolutionGroups
	}
	return a[x].Word[0] < a[y].Word[0]
}

func readWords(r io.Reader) []wordle.Word {
	words := make([]wordle.Word, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		words = append(words, wordle.NewWord(scanner.Bytes()))
	}
	return words
}
