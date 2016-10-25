package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const minWordLength = 3

type boggleSolver struct {
	rows       int
	cols       int
	adjList    [][]int
	dictionary OptimizedTrie
}

func buildAdjList(rows, cols int) [][]int {
	ret := make([][]int, rows*cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			i := r*cols + c
			ret[i] = make([]int, 0)
			for _, deltar := range []int{-1, 0, 1} {
				targetr := r + deltar
				if targetr < 0 || targetr == rows {
					continue
				}
				for _, deltac := range []int{-1, 0, 1} {
					targetc := c + deltac
					if targetc < 0 || targetc == cols {
						continue
					}
					if deltac == 0 && deltar == 0 {
						continue
					}
					ret[i] = append(ret[i], targetr*cols+targetc)
				}
			}
		}
	}
	return ret
}

func buildScore(rows, cols int) []int {
	score := make([]int, rows*cols+1)
	for i := 0; i < rows*cols+1; i++ {
		switch {
		case i < minWordLength:
			score[i] = 0
		case i < 5:
			score[i] = 1
		case i == 5:
			score[i] = 2
		case i == 6:
			score[i] = 3
		case i == 7:
			score[i] = 5
		default:
			score[i] = 11
		}
	}
	return score
}

func newSolver(rows, cols int, dictfile string) (*boggleSolver, error) {
	maxWordLength := rows * cols

	adjList := buildAdjList(rows, cols)
	score := buildScore(rows, cols)

	file, err := os.Open(dictfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var dictionary OptimizedTrie
	for scanner.Scan() {
		w := scanner.Text()
		if len(w) >= minWordLength && len(w) <= maxWordLength {
			dictionary.Insert(strings.ToUpper(w), score[len(w)])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	solver := boggleSolver{
		rows:       rows,
		cols:       cols,
		adjList:    adjList,
		dictionary: dictionary,
	}
	return &solver, nil
}

func (bs *boggleSolver) score(bb Boggler) (int, *OptimizedTrie) {
	visited := make([]bool, len(bs.adjList))
	var results OptimizedTrie
	var buf bytes.Buffer
	score := 0

	for p := 0; p < len(bs.adjList); p++ {
		score += bs.dfs(bb, &bs.dictionary, p, &visited, &results, &buf)
	}

	return score, &results
}

func (bs *boggleSolver) dfs(bb Boggler, dictionary *OptimizedTrie, p int, visited *[]bool, results *OptimizedTrie, sb *bytes.Buffer) int {

	if (*visited)[p] {
		return 0
	}

	letter := bb.GetLinear(p)
	score := 0

	subtrie := dictionary.SubtrieR(letter)
	if subtrie == nil {
		return score
	}

	(*visited)[p] = true
	sb.WriteRune(letter)
	if letter == 'Q' {
		sb.WriteRune('U')
	}

	score += subtrie.RootValue()
	if score > 0 {
		str := sb.String()
		if !results.Has(str) {
			results.Insert(str, score)
		} else {
			score = 0
		}
	}

	for _, p2 := range bs.adjList[p] {
		score += bs.dfs(bb, subtrie, p2, visited, results, sb)
	}

	(*visited)[p] = false
	sb.Truncate(sb.Len() - 1)
	if letter == 'Q' {
		sb.Truncate(sb.Len() - 1)
	}

	return score
}

func frequencyCount(dictfile string, maxWordLength int) ([][]float64, error) {
	file, err := os.Open(dictfile)
	freqs := make([][]float64, 26)
	if err != nil {
		return freqs, err
	}
	defer file.Close()

	sum := 0
	counts := make([][]int, 26)
	for i := range counts {
		counts[i] = make([]int, 26)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w := scanner.Text()
		if len(w) < minWordLength || len(w) > maxWordLength {
			continue
		}
		w = strings.ToUpper(w)
		for l := 0; l < len(w)-1; l++ {
			sum++
			l1 := w[l] - 'A'
			l2 := w[l+1] - 'A'
			counts[l1][l2]++
			counts[l2][l1]++
		}
	}

	if err := scanner.Err(); err != nil {
		return freqs, err
	}

	for i, cnt := range counts {
		freqs[i] = make([]float64, 26)
		for j, c := range cnt {
			freqs[i][j] = float64(c) / float64(sum)
		}
	}

	return freqs, nil
}

type boardScore struct {
	score int
	board []string
}

func solve(dictfile string, best chan int, brd chan boardScore, flip chan bool) {
	bs, err := newSolver(4, 4, dictfile)
	if err != nil {
		panic(err)
	}

	freqs, err := frequencyCount(dictfile, 16)
	if err != nil {
		panic(err)
	}
	board := newDiceBoard(4, 4, boggle1992)

	topscore := 0
	score := 0

	for {
		select {

		case <-flip:
			board = newDiceBoard(4, 4, boggle1992)

		case topscore = <-best:
			// Already did what I wanted to do...

		default:
			last := board.Clone()
			lastScore := score

			board.DictShuffle(bs.adjList, freqs)
			score, _ = bs.score(board)

			if score > topscore {
				topscore = score
				brd <- boardScore{score: topscore, board: board.ArrayLinear()}
			}

			if score <= lastScore && rand.Float64() > float64(score)/float64(lastScore) {
				board = last.(*DiceBoard)
				score = lastScore
			}
		}
	}
}

func main() {

	seed := time.Now().UnixNano()
	rand.Seed(seed)
	dictfile := filepath.Join("dictionaries", "dictionary-enable1.txt")

	nprocs := runtime.GOMAXPROCS(0)
	best := make([]chan int, nprocs)
	brd := make(chan boardScore, nprocs)
	flip := make(chan bool)

	for i := 0; i < nprocs; i++ {
		best[i] = make(chan int, 100)
		go solve(dictfile, best[i], brd, flip)
	}

	i := 0
	topscore := 0
	start := int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)

	for {
		select {
		case b := <-brd:
			if b.score > topscore {
				topscore = b.score
				fmt.Printf("%d,%d,%d,%s\n", i, int64(time.Nanosecond)*time.Now().UnixNano()/int64(time.Millisecond)-start, topscore, strings.Join(b.board, ","))
				for _, bst := range best {
					bst <- topscore
				}
			}
		case <-time.After(time.Minute * 5):
			flip <- true
			i++
		}
	}
}
