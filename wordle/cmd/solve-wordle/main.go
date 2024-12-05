package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/reallyasi9/riddler/wordle/pkg/wordle"
)

func init() {
	flag.CommandLine.Usage = func() {
		name, _ := os.Executable()
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <solutions-file> <guesses-file> [GUESS1 RESULT1 [GUESS2 RESULT2...]]\n", filepath.Base(name))
		flag.PrintDefaults()
	}
}

func main() {

	flag.Parse()
	if flag.NArg() < 2 {
		flag.CommandLine.Usage()
		name, _ := os.Executable()
		log.Fatalf("%s requres two positional arguments: a file containing a list of possible solutions and a file containing a list of possible guesses", filepath.Base(name))
	}
	if (flag.NArg()-2)%2 != 0 {
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

	initialGuesses := make([]wordle.Word, (flag.NArg()-2)/2)
	initialStatus := make([]wordle.WordStatus, len(initialGuesses))
	for iarg := 2; iarg < flag.NArg(); iarg += 2 {
		guess := flag.Arg(iarg)
		status := flag.Arg(iarg + 1)
		if len(guess) != wordle.WORD_SIZE {
			log.Fatalf("guesses must be %d letters", wordle.WORD_SIZE)
		}
		if len(status) != wordle.WORD_SIZE {
			log.Fatalf("status results must be %d characters", wordle.WORD_SIZE)
		}
		initialGuesses[iarg/2-1] = wordle.NewWordFromString(strings.ToLower(guess))
		initialStatus[iarg/2-1] = wordle.NewWordStatus(status)
	}

	wdl := wordle.NewWordle(initialSolutions)
	wdl.Filter(initialGuesses, initialStatus)

	if wdl.Len() == 1 {
		fmt.Print("There is only one possible soution remaining.\n")
	} else {
		fmt.Printf("There are %d solutions remaining.\n", wdl.Len())
	}
	for i, soln := range wdl.Solutions() {
		fmt.Printf("%s\n", soln)
		if i == 9 {
			break
		}
	}
	if wdl.Len() > 10 {
		fmt.Println("...")
	}
	if wdl.Len() <= 2 {
		return
	}

	log.Println("Finding the guess that minimizes entropy")
	allSolutions := make(map[wordle.Word]struct{})
	for _, soln := range wdl.Solutions() {
		allSolutions[soln] = struct{}{}
	}
	bestGuesses := []EntropyWord{}
	bestEntropy := math.Inf(1)
	for _, guess := range guessables {
		entropy, probability, deduced := wdl.Try([]wordle.Word{guess})
		if entropy <= bestEntropy {
			_, isSoln := allSolutions[guess]
			log.Printf("%s = %f (p=%f, %d deduced, solution %t)\n", guess, entropy, probability, deduced, isSoln)
			if entropy < bestEntropy {
				bestGuesses = bestGuesses[:0]
			}
			bestEntropy = entropy
			bestGuesses = append(bestGuesses, EntropyWord{Word: guess, Entropy: entropy, Probability: probability, Deduced: deduced, IsSolution: isSoln})
		}
	}

	fmt.Println("Best guesses:")
	sort.Sort(ByEntropy(bestGuesses))
	for i, bg := range bestGuesses {
		fmt.Printf("%d: %s = %f (p=%f, %d deduced, solution %t)\n", i+1, bg.Word, bg.Entropy, bg.Probability, bg.Deduced, bg.IsSolution)
		if i == 9 {
			break
		}
	}
	if len(bestGuesses) > 10 {
		fmt.Println("...")
	}
}

type EntropyWord struct {
	Word        wordle.Word
	Entropy     float64
	Probability float64
	Deduced     int
	IsSolution  bool
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
		return a[x].Entropy > a[y].Entropy // more entropy is worse
	}
	if a[x].Probability != a[y].Probability {
		return a[x].Probability < a[y].Probability
	}
	if a[x].Deduced != a[y].Deduced {
		return a[x].Deduced < a[y].Deduced
	}
	if a[x].IsSolution != a[y].IsSolution {
		return a[x].IsSolution
	}
	return false
}

func readWords(r io.Reader) []wordle.Word {
	words := make([]wordle.Word, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		words = append(words, wordle.NewWordFromString(string(scanner.Bytes())))
	}
	return words
}
