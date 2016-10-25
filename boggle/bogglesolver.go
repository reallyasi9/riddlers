package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const minWordLength = 3

var score = []int{0, 0, 0, 1, 1, 2, 3, 5, 11, 11, 11, 11, 11, 11, 11, 11, 11}

type boggleSolver struct {
	rows       int
	cols       int
	adjList    [][]int
	dictionary OptimizedTrie
}

func buildAdjList(rows, cols int) [][]int {
	ret := make([][]int, rows*cols)
	for i := 0; i < rows*cols; i++ {
		ret[i] = make([]int, 0)
		r := i % cols
		c := i / cols
		for deltar := range []int{-1, 0, 1} {
			targetr := r + deltar
			if targetr < 0 || targetr >= rows {
				continue
			}
			for deltac := range []int{-1, 0, 1} {
				targetc := c + deltac
				if targetc < 0 || targetc >= cols {
					continue
				}
				ret[i] = append(ret[i], targetr*cols+targetc)
			}
		}
	}
	return ret
}

func newSolver(rows, cols int, dictfile string) *boggleSolver {
	maxWordLength := rows * cols

	adjList := buildAdjList(rows, cols)

	file, err := os.Open(dictfile)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	solver := boggleSolver{rows: rows,
		cols:       cols,
		adjList:    adjList,
		dictionary: dictionary}
	return &solver
}

func (bs *boggleSolver) score(bb *BoggleBoard) int {
	visited := make([]bool, len(bs.adjList))
	var results OptimizedTrie
	var buf bytes.Buffer
	var score int

	for p := 0; p < len(bs.adjList); p++ {
		score += bs.dfs(bb, &bs.dictionary, p, &visited, &results, &buf)
	}

	return score
}

func (bs *boggleSolver) dfs(bb *BoggleBoard, dictionary *OptimizedTrie, p int, visited *[]bool, results *OptimizedTrie, sb *bytes.Buffer) int {

	if (*visited)[p] {
		return 0
	}

	letter := (*bb)[p/bs.rows][p%bs.cols]
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

func main() {

	seed := time.Now().UnixNano()
	rand.Seed(seed)

	bs := newSolver(4, 4, "dictionaries/dictionary-enable1.txt")

	topScore := 0
	i := 0
	for {
		i++
		board := NewBoggleBoardRandom(4, 4)
		score := bs.score(board)
		if score > topScore {
			topScore = score
			fmt.Printf("Score %d found at iteration %d\n%s\n", score, i, board)
		}
	}
}
