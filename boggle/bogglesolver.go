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

var adjList = make([][]int, 16)

const minWordLength = 3
const maxWordLength = 16

var dictionary = OptimizedTrie{}

var score = []int{0, 0, 0, 1, 1, 2, 3, 5, 11, 11, 11, 11, 11, 11, 11, 11, 11}

func init() {
	adjList[0] = []int{1, 4, 5}
	adjList[1] = []int{0, 2, 4, 5, 6}
	adjList[2] = []int{1, 3, 5, 6, 7}
	adjList[3] = []int{2, 6, 7}
	adjList[4] = []int{0, 1, 5, 8, 9}
	adjList[5] = []int{0, 1, 2, 4, 6, 8, 9, 10}
	adjList[6] = []int{1, 2, 3, 5, 7, 9, 10, 11}
	adjList[7] = []int{2, 3, 6, 10, 11}
	adjList[8] = []int{4, 5, 9, 12, 13}
	adjList[9] = []int{4, 5, 6, 8, 10, 12, 13, 14}
	adjList[10] = []int{5, 6, 7, 9, 11, 13, 14, 15}
	adjList[11] = []int{6, 7, 10, 14, 15}
	adjList[12] = []int{8, 9, 13}
	adjList[13] = []int{8, 9, 10, 12, 14}
	adjList[14] = []int{9, 10, 11, 13, 15}
	adjList[15] = []int{10, 11, 14}
}

func (dict *OptimizedTrie) fill(w string) {
	if len(w) < minWordLength || len(w) > maxWordLength {
		return
	}
	dict.Insert(strings.ToUpper(w), score[len(w)])
}

func (bb *BoggleBoard) score() int {
	visited := make([]bool, 16)
	var results OptimizedTrie
	var buf bytes.Buffer
	var score int

	for p := 0; p < 16; p++ {
		score += bb.dfs(&dictionary, p, &visited, &results, &buf)
	}

	return score
}

func (bb *BoggleBoard) dfs(dict *OptimizedTrie, p int, visited *[]bool, results *OptimizedTrie, sb *bytes.Buffer) int {
	if dict == nil {
		return 0
	}

	if (*visited)[p] {
		return 0
	}

	letter := (*bb)[p/len(*bb)][p%len((*bb)[0])]
	score := 0

	subtrie := dict.SubtrieR(letter)
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

	for _, p2 := range adjList[p] {
		score += bb.dfs(subtrie, p2, visited, results, sb)
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

	file, err := os.Open("enable1.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dictionary.fill(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	topScore := 0
	i := 0
	for {
		i++
		board := NewBoggleBoard()
		score := board.score()
		if score > topScore {
			topScore = score
			fmt.Printf("Score %d found at iteration %d\n%s\n", score, i, board)
		}
	}
}
