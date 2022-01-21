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
	"github.com/kelindar/bitmap"
	"gonum.org/v1/gonum/stat/combin"
)

var nGuesses = flag.Int("g", 2, "(exact) number of guesses to optimize after starting guesses")
var forceDisjoint = flag.Bool("d", false, "force all words in guess to have unique letters, including forbidding guesses to have duplicate letters")
var startingWords = flag.String("s", "", "comma-separated list of starting guesses")

const alphabetSize = 26
const wordSize = 5
const firstLetter = 'a'

type Word [wordSize]byte

func makeWord(bs []byte) Word {
	var w Word
	for i, b := range bs {
		if i >= wordSize {
			break
		}
		w[i] = b - firstLetter
	}
	return w
}

func makeWordFromString(s string) Word {
	return makeWord([]byte(s))
}

func (w Word) String() string {
	b := make([]byte, wordSize)
	for i := 0; i < wordSize; i++ {
		b[i] = w[i] + firstLetter
	}
	return string(b)
}

type Wordle struct {
	words []Word

	wordsContainingLetter           [alphabetSize]bitmap.Bitmap
	wordsContainingLetterByPosition [wordSize][alphabetSize]bitmap.Bitmap
	wordsNotContainingLetter        [alphabetSize]bitmap.Bitmap
}

func NewWordle(words []Word) *Wordle {
	ws := make([]Word, len(words))
	copy(ws, words)

	wordle := Wordle{words: ws}
	for i, word := range words {
		for letter := byte(0); letter < alphabetSize; letter++ {
			found := false
			for pos, character := range word {
				if character == letter {
					wordle.wordsContainingLetterByPosition[pos][character].Set(uint32(i))
					wordle.wordsContainingLetter[character].Set(uint32(i))
					found = true
				}
			}
			if !found {
				wordle.wordsNotContainingLetter[letter].Set(uint32(i))
			}
		}
	}

	return &wordle
}

func (w *Wordle) GetWord(i int) Word {
	return w.words[i]
}

func (w *Wordle) NWords() int {
	return len(w.words)
}

const (
	MISSING int = iota
	PRESENT
	CORRECT
)

func (w Word) Compare(soln Word) [5]int {
	var status [5]int

OUTER:
	for i, c := range w {
		for j, x := range soln {
			if c == x {
				if i == j {
					status[i] = CORRECT
					continue OUTER
				} else {
					status[i] = PRESENT
				}
			}
		}
	}

	return status
}

func (w *Wordle) Ambiguities(guesses []Word, soln Word) bitmap.Bitmap {
	var possible bitmap.Bitmap
	possible.Grow(uint32(len(w.words)))
	possible.Ones()

	for _, guess := range guesses {
		status := guess.Compare(soln)
		for letter, stat := range status {
			switch stat {
			case CORRECT:
				possible.And(w.wordsContainingLetterByPosition[letter][guess[letter]])
			case PRESENT:
				possible.And(w.wordsContainingLetter[guess[letter]])
			case MISSING:
				possible.And(w.wordsNotContainingLetter[guess[letter]])
			default:
			}
		}
	}
	return possible
}

type ComboProb struct {
	Combination []Word
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

	start := []Word{}
	for _, word := range strings.Split(*startingWords, ",") {
		start = append(start, makeWordFromString(word))
	}

	solns := readWords(solnFile)
	guesses := readWords(guessFile)

	wordle := NewWordle(solns)

	combinations := wordCombinations(guesses, start, *nGuesses)                       // produce combinations
	unfiltered := calculateProbabilities(wordle, solns, *forceDisjoint, combinations) // multi-thread calculate solutions
	filtered := filterBest(unfiltered)                                                // merge

	for cp := range filtered {
		log.Print(cp)
	}
}

func calculateProbabilities(wordle *Wordle, solns []Word, disjoint bool, in <-chan []Word) <-chan ComboProb {
	out := make(chan ComboProb, 1000)
	filter := func(words []Word) bool {
		return true
	}
	if disjoint {
		filter = func(words []Word) bool {
			return disjointLetters(words)
		}
	}
	nsolns := float64(len(solns))
	go func() {
		var wg sync.WaitGroup
		for words := range in {
			wg.Add(1)
			go func(words []Word) {
				defer wg.Done()

				if !filter(words) {
					return
				}

				prob := 0.
				deduced := 0
				for _, solution := range solns {
					ambiguities := wordle.Ambiguities(words, solution).Count()
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
		best := ComboProb{Combination: make([]Word, 0)}
		for cp := range in {
			if cp.Probability > best.Probability {
				if len(best.Combination) != len(cp.Combination) {
					best.Combination = make([]Word, len(cp.Combination))
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

func readWords(r io.Reader) []Word {
	words := make([]Word, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		words = append(words, makeWord(scanner.Bytes()))
	}
	return words
}

func disjointLetters(ws []Word) bool {
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

func wordCombinations(ws []Word, start []Word, n int) <-chan []Word {
	if len(start)+n > 6 {
		panic("a maximum of only 6 guesses are allowed!")
	}
	out := make(chan []Word, 1000)
	go func() {
		numComb := combin.Binomial(len(ws), n)
		bar := pb.ProgressBarTemplate(pb.Full).Start(numComb)
		gen := combin.NewCombinationGenerator(len(ws), n)
		comb := make([]int, n)
		for gen.Next() {
			gen.Combination(comb)
			words := make([]Word, n+len(start))
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
