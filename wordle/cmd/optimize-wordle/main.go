package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"github.com/gonum/stat/combin"
	"github.com/reallyasi9/riddler/wordle/pkg/wordle"
)

var nGuesses = flag.Int("g", 2, "(exact) number of guesses to optimize after starting guesses")
var forceDisjoint = flag.Bool("d", false, "force all words in all guesses to have mutually unique letters")
var startingWords = flag.String("s", "", "comma-separated list of starting guesses")

func init() {
	log.SetOutput(os.Stdout)
}

func main() {

	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatal("wordle requres two positional arguments: a file containing a list of possible solutions and a file containing a list of possible guesses")
	}

	solnFile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer solnFile.Close()

	guessFile, err := os.Open(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	defer guessFile.Close()

	start := []wordle.Word{}
	for _, word := range strings.Split(*startingWords, ",") {
		if len(word) != 5 {
			continue
		}
		start = append(start, wordle.NewWordFromString(word))
	}

	solns := readWords(solnFile)
	guesses := readWords(guessFile)

	wordle := wordle.NewWordle(solns)

	combinations := wordCombinations(guesses, start, *nGuesses)                       // produce combinations
	unfiltered := calculateProbabilities(wordle, solns, *forceDisjoint, combinations) // multi-thread calculate solutions
	filtered := filterBest(unfiltered)                                                // merge

	for cp := range filtered {
		log.Print(cp)
	}
}

type ComboProb struct {
	Combination []wordle.Word
	Probability float64
	Deduced     int
}

func (cp ComboProb) String() string {
	words := make([]string, len(cp.Combination))
	for i, w := range cp.Combination {
		words[i] = w.String()
	}
	joined := strings.Join(words, " + ")
	return fmt.Sprintf("%s = %f (%d deduced)", joined, cp.Probability, cp.Deduced)
}

func calculateProbabilities(wdl *wordle.Wordle, solns []wordle.Word, disjoint bool, in <-chan []wordle.Word) <-chan ComboProb {
	out := make(chan ComboProb, 1000)
	filter := func(words []wordle.Word) bool {
		return true
	}
	if disjoint {
		filter = func(words []wordle.Word) bool {
			return disjointLetters(words)
		}
	}
	nsolns := float64(len(solns))
	go func() {
		var wg sync.WaitGroup
		for words := range in {
			wg.Add(1)
			go func(words []wordle.Word) {
				defer wg.Done()

				if !filter(words) {
					return
				}

				prob := 0.
				deduced := 0
				for _, solution := range solns {
					ambiguities := wdl.Ambiguities(words, solution).Count()
					if ambiguities == 1 {
						deduced++
					}
					prob += 1. / float64(ambiguities)
				}
				out <- ComboProb{Combination: words, Probability: prob / nsolns, Deduced: deduced}
			}(words)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func filterBest(in <-chan ComboProb) <-chan ComboProb {
	out := make(chan ComboProb, 1000)
	go func() {
		best := ComboProb{Combination: make([]wordle.Word, 0)}
		for cp := range in {
			if cp.Probability > best.Probability {
				if len(best.Combination) != len(cp.Combination) {
					best.Combination = make([]wordle.Word, len(cp.Combination))
				}
				copy(best.Combination, cp.Combination)
				best.Probability = cp.Probability
				best.Deduced = cp.Deduced
				out <- cp
			}
		}
		close(out)
	}()
	return out
}

func readWords(r io.Reader) []wordle.Word {
	words := make([]wordle.Word, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		words = append(words, wordle.NewWord(scanner.Bytes()))
	}
	return words
}

func disjointLetters(ws []wordle.Word) bool {
	if len(ws) == 0 {
		return true
	}
	letters := make(map[byte]struct{})
	for i, w := range ws {
		for _, l := range w {
			letters[l] = struct{}{}
		}
		if len(letters) != (i+1)*5 {
			return false
		}
	}
	return true
}

func wordCombinations(ws []wordle.Word, start []wordle.Word, n int) <-chan []wordle.Word {
	if len(start)+n > 6 {
		panic("a maximum of only 6 guesses are allowed!")
	}
	out := make(chan []wordle.Word, 1000)
	go func() {
		numComb := combin.Binomial(len(ws), n)
		bar := pb.ProgressBarTemplate(pb.Full).Start(numComb)
		gen := combin.NewCombinationGenerator(len(ws), n)
		comb := make([]int, n)
		for gen.Next() {
			gen.Combination(comb)
			words := make([]wordle.Word, n+len(start))
			copy(words, start)
			for i, j := range comb {
				words[i+len(start)] = ws[j]
			}
			out <- words
			bar.Increment()
		}
		bar.Finish()
		close(out)
	}()
	return out
}
