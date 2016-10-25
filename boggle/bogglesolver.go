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

func main() {

	seed := time.Now().UnixNano()
	rand.Seed(seed)

	bs, err := newSolver(4, 4, "dictionaries/dictionary-enable1.txt")
	if err != nil {
		panic(err)
	}

	topScore := 0
	i := 0
	for {
		i++
		var board Boggler
		board = NewBoggleBoardRandom(4, 4)
		// board = newDiceBoard(4, 4, boggle1992)
		score, _ := bs.score(board)
		if score > topScore {
			topScore = score
			fmt.Printf("Score %d found at iteration %d\n%s\n", score, i, board)
		}
	}
}
