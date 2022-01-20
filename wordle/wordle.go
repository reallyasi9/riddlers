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

	"gonum.org/v1/gonum/stat/combin"
)

var nGuesses = flag.Int("g", 2, "(exact) number of guesses to optimize")

type Word [5]byte

func makeWord(bs []byte) Word {
	var w Word
	for i, b := range bs {
		if i > 4 {
			break
		}
		w[i] = b - 'a'
	}
	return w
}

func (w Word) String() string {
	b := make([]byte, 5)
	for i := 0; i < 5; i++ {
		b[i] = w[i] + 'a'
	}
	return string(b)
}

type WordSet struct {
	set map[Word]struct{}
}

func NewWordSet(words ...Word) *WordSet {
	set := make(map[Word]struct{})
	for _, word := range words {
		set[word] = struct{}{}
	}
	return &WordSet{set: set}
}

func (s *WordSet) Len() int {
	return len(s.set)
}

func (s *WordSet) Clone() *WordSet {
	set := make(map[Word]struct{})
	for word := range s.set {
		set[word] = struct{}{}
	}
	return &WordSet{set: set}
}

func (s *WordSet) Insert(w Word) {
	s.set[w] = struct{}{}
}

func (s *WordSet) Intersection(other *WordSet) *WordSet {
	out := make(map[Word]struct{})
	for val := range s.set {
		if _, exists := other.set[val]; exists {
			out[val] = struct{}{}
		}
	}
	return &WordSet{set: out}
}

// KeepIntersection is a very special function: it will either drop the non-intersecting values from s that are not in other, or it will replace s with other if s is empty.
func (s *WordSet) KeepIntersection(other *WordSet) {
	if len(s.set) == 0 {
		for word := range other.set {
			s.set[word] = struct{}{}
		}
		return
	}
	for val := range s.set {
		if _, exists := other.set[val]; !exists {
			delete(s.set, val)
		}
	}
}

type Wordle struct {
	byLetter         map[byte]*WordSet
	byLetterPosition [5]map[byte]*WordSet
	byMissing        map[byte]*WordSet
}

func NewWordle(words []Word) *Wordle {
	letters := make(map[byte]*WordSet)
	var matches [5]map[byte]*WordSet
	misses := make(map[byte]*WordSet)

	for l := 0; l < 5; l++ {
		matches[l] = make(map[byte]*WordSet)
	}
	for i := byte(0); i < 26; i++ {
		letters[i] = NewWordSet()
		misses[i] = NewWordSet()
		for l := 0; l < 5; l++ {
			matches[l][i] = NewWordSet()
		}
	}

	for _, word := range words {
		var i byte
		for i = 0; i < 26; i++ {
			found := false
			for l, b := range word {
				if b == i {
					matches[l][b].Insert(word)
					letters[b].Insert(word)
					found = true
				}
			}
			if !found {
				misses[i].Insert(word)
			}
		}
	}

	return &Wordle{byLetter: letters, byLetterPosition: matches, byMissing: misses}
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

func (w *Wordle) Ambiguities(guesses []Word, soln Word) *WordSet {
	possible := NewWordSet()
	for _, guess := range guesses {
		status := guess.Compare(soln)
		for letter, stat := range status {
			switch stat {
			case CORRECT:
				possible.KeepIntersection(w.byLetterPosition[letter][guess[letter]])
			case PRESENT:
				possible.KeepIntersection(w.byLetter[guess[letter]])
			default:
				possible.KeepIntersection(w.byMissing[guess[letter]])
			}

		}
	}
	return possible
}

type ComboProb struct {
	Combination []Word
	Probability float64
}

func (cp *ComboProb) String() string {
	words := make([]string, len(cp.Combination))
	for i, w := range cp.Combination {
		words[i] = w.String()
	}
	joined := strings.Join(words, " + ")
	return fmt.Sprintf("%s = %f", joined, cp.Probability)
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

	solns := readWords(solnFile)
	guesses := readWords(guessFile)

	wordle := NewWordle(solns)

	combinations := wordCombinations(guesses, *nGuesses)              // produce combinations
	unfiltered := calculateProbabilities(wordle, solns, combinations) // multi-thread calculate solutions
	filtered := filterBest(unfiltered)                                // merge

	for cp := range filtered {
		log.Print(cp)
	}
}

func calculateProbabilities(wordle *Wordle, solns []Word, in <-chan []Word) <-chan ComboProb {
	out := make(chan ComboProb, 100)
	go func() {
		var wg sync.WaitGroup
		for words := range in {
			wg.Add(1)
			go func(words []Word) {
				defer wg.Done()

				if !disjointLetters(words) {
					return
				}

				prob := 0.
				for _, solution := range solns {
					ambiguities := wordle.Ambiguities(words, solution).Len()
					prob += 1. / float64(ambiguities)
				}
				out <- ComboProb{Combination: words, Probability: prob}
			}(words)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func filterBest(in <-chan ComboProb) <-chan ComboProb {
	out := make(chan ComboProb, 100)
	go func() {
		best := ComboProb{Combination: make([]Word, *nGuesses), Probability: 0}
		for cp := range in {
			if cp.Probability > best.Probability {
				copy(best.Combination, cp.Combination)
				best.Probability = cp.Probability
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

func wordCombinations(ws []Word, n int) <-chan []Word {
	out := make(chan []Word, 100)
	go func() {
		gen := combin.NewCombinationGenerator(len(ws), n)
		comb := make([]int, n)
		for gen.Next() {
			gen.Combination(comb)
			words := make([]Word, n)
			for i, j := range comb {
				words[i] = ws[j]
			}
			out <- words
		}
		close(out)
	}()
	return out
}
