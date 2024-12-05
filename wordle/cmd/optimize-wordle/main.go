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
	"strings"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"github.com/gonum/stat/combin"
	"github.com/kelindar/bitmap"
	"github.com/reallyasi9/riddler/wordle/pkg/wordle"
)

var nGuesses = flag.Int("g", 2, "(exact) number of guesses to optimize after starting set")
var forceDisjoint = flag.Bool("d", false, "force all words in all guesses to have mutually unique letters")
var optVariable = flag.String("v", "entropy", "optimization variable (one of 'entropy', 'probability', or 'deduced')")

func init() {
	log.SetOutput(os.Stdout)
	flag.CommandLine.Usage = func() {
		name, _ := os.Executable()
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <solutions-file> <guesses-file> [-d] [-g N] [-v {entropy,probability,deduced}] [GUESS [GUESS...]]\n", filepath.Base(name))
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

	if *nGuesses > wordle.MAX_GUESSES {
		log.Fatalf("a maximum of %d guesses are allowed", wordle.MAX_GUESSES)
	}

	if *optVariable != "entropy" && *optVariable != "probability" && *optVariable != "deduced" {
		log.Fatalf("optimization variable not allowed")
	}

	solns := readWords(solnFile)
	guesses := readWords(guessFile)

	start := make([]wordle.Word, flag.NArg()-2)
	for i, arg := range flag.Args()[2:] {
		start[i] = wordle.NewWordFromString(arg)
	}

	wordle := wordle.NewWordle(solns)

	combinations := wordCombinations(guesses, start, *nGuesses)          // produce combinations
	unfiltered := calculateEntropy(wordle, *forceDisjoint, combinations) // multi-thread calculate solutions
	filtered := filterBest(unfiltered, *optVariable)                     // merge

	for cp := range filtered {
		log.Print(cp)
	}
}

type ComboProb struct {
	Combination []wordle.Word
	Entropy     float64
	Probability float64
	Deduced     int
}

func (cp ComboProb) Compare(other ComboProb) int {
	if cp.Entropy < other.Entropy {
		return 1
	} else if cp.Entropy > other.Entropy {
		return -1
	} else if cp.Probability > other.Probability {
		return 1
	} else if cp.Probability < other.Probability {
		return -1
	} else if cp.Deduced > other.Deduced {
		return 1
	} else if cp.Deduced > other.Deduced {
		return -1
	}
	return 0
}

func (cp ComboProb) String() string {
	words := make([]string, len(cp.Combination))
	for i, w := range cp.Combination {
		words[i] = w.String()
	}
	joined := strings.Join(words, " + ")
	return fmt.Sprintf("%s = %f (p=%f, %d deduced)", joined, cp.Entropy, cp.Probability, cp.Deduced)
}

func calculateEntropy(wdl *wordle.Wordle, disjoint bool, in <-chan []wordle.Word) <-chan ComboProb {
	out := make(chan ComboProb, 1024)
	filter := func(words []wordle.Word) bool {
		return true
	}
	if disjoint {
		filter = func(words []wordle.Word) bool {
			return disjointLetters(words)
		}
	}
	go func() {
		var wg sync.WaitGroup
		for words := range in {
			wg.Add(1)
			go func(words []wordle.Word) {
				defer wg.Done()

				if !filter(words) {
					return
				}

				entropy, probability, deduced := wdl.Try(words)
				out <- ComboProb{Combination: words, Entropy: entropy, Probability: probability, Deduced: deduced}
			}(words)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func filterBest(in <-chan ComboProb, variable string) <-chan ComboProb {
	out := make(chan ComboProb, 1024)
	go func() {
		best := ComboProb{Combination: make([]wordle.Word, 0), Entropy: math.Inf(1)}
		for cp := range in {
			var better bool
			if variable == "entropy" {
				better = cp.Entropy < best.Entropy
			} else if variable == "probability" {
				better = cp.Probability > best.Probability
			} else if variable == "deduced" {
				better = cp.Deduced > best.Deduced
			}
			if better {
				if len(best.Combination) != len(cp.Combination) {
					best.Combination = make([]wordle.Word, len(cp.Combination))
				}
				copy(best.Combination, cp.Combination)
				best.Entropy = cp.Entropy
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
		words = append(words, wordle.NewWordFromString(string(scanner.Bytes())))
	}
	return words
}

func disjointLetters(ws []wordle.Word) bool {
	if len(ws) == 0 {
		return true
	}
	letters := bitmap.Bitmap{(1 << wordle.N_LETTERS)}
	for _, w := range ws {
		for _, l := range w {
			if letters.Contains(uint32(l)) {
				return false
			}
			letters.Set(uint32(l))
		}
	}
	return true
}

func wordCombinations(ws []wordle.Word, start []wordle.Word, n int) <-chan []wordle.Word {
	out := make(chan []wordle.Word, 1024)
	go func() {
		numComb := combin.Binomial(len(ws), n)
		bar := pb.ProgressBarTemplate(pb.Full).Start(numComb)
		gen := combin.NewCombinationGenerator(len(ws), n)
		comb := make([]int, n)
		for gen.Next() {
			gen.Combination(comb)
			words := make([]wordle.Word, n+len(start))
			copy(words[:len(start)], start)
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
