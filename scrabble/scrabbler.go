package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	radix "github.com/armon/go-radix"
)

var runeScores = map[rune]int{
	'?': 0,
	'e': 1,
	'a': 1,
	'i': 1,
	'o': 1,
	'n': 1,
	'r': 1,
	't': 1,
	'l': 1,
	's': 1,
	'u': 1,
	'd': 2,
	'g': 2,
	'b': 3,
	'c': 3,
	'm': 3,
	'p': 3,
	'f': 4,
	'h': 4,
	'v': 4,
	'w': 4,
	'y': 4,
	'k': 5,
	'j': 8,
	'x': 8,
	'q': 10,
	'z': 10,
}

var scoreTrie = radix.New()

func scoreWord(word string) (int, error) {
	score := 0
	for _, r := range word {
		value, ok := runeScores[r]
		if !ok {
			err := fmt.Errorf("Rune %c not a recognized scrabble letter", r)
			log.Printf("Error looking up word: %s", err)
			return 0, err
		}
		score += value
	}
	return score, nil
}

var dictionaryURL = "https://norvig.com/ngrams/enable1.txt"

func buildDictionary(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error with request: %s", err)
		return err
	}
	defer resp.Body.Close()

	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		word := s.Text()
		score, err := scoreWord(word)
		if err != nil {
			log.Printf("Error scoring word: %s", err)
			return err
		}
		scoreTrie.Insert(word, score)
	}

	if err := s.Err(); err != nil {
		log.Printf("Error reading response: %s", err)
		return err
	}
	return nil
}

var letterCollection = []rune("??aaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyyz")

func generateBoard() []rune {
	out := make([]rune, len(letterCollection))
	copy(out, letterCollection)
	rand.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})
	return out
}

func replaceQMs(board []rune) string {
	for i, r := range board {
		if r == '?' {
			board[i] = rune(rand.Int31n(122-97) + 97)
		}
	}
	return string(board)
}

type prefixWalker struct {
	isPrefix bool
}

func (p *prefixWalker) walk(s string, v interface{}) bool {
	p.isPrefix = true
	return false // stop iterating immediately
}

func scoreBoard(board string) int {
	score := 0
	foundTrie := radix.New()
	i := 0
	j := 1
	for i < len(board) {
		sub := board[i:j]

		if val, ok := scoreTrie.Get(sub); ok {
			if _, alreadyFound := foundTrie.Get(sub); !alreadyFound {
				// Is a word: score it
				score += val.(int)
				foundTrie.Insert(sub, true)
			}
		}

		var walker prefixWalker
		scoreTrie.WalkPrefix(sub, walker.walk)
		if walker.isPrefix {
			// Is a prefix: increase j
			j++
			if j > len(board) {
				// No more letters
				break
			}
		} else {
			// Scoot forward
			i++
			if i == j {
				j++
			}
		}
	}
	return score
}

func mutateBoard(board []rune) {
	i := rand.Intn(len(board))
	j := rand.Intn(len(board))
	board[i], board[j] = board[j], board[i]
}

func generateOffspring(p1 []rune, p2 []rune) []rune {
	child := make([]rune, len(p1))
	for i := 0; i < len(child); i++ {
		if rand.Float32() < .5 {
			child[i] = p1[i]
		} else {
			child[i] = p2[i]
		}
	}
	return child
}

const generations = 1000
const perGeneration = 1000

func main() {
	log.Println("Starting")

	err := buildDictionary(dictionaryURL)
	if err != nil {
		log.Fatalln(err)
	}

	rawBoard1 := generateBoard()
	rawBoard2 := generateBoard()
	cleanBoard1 := replaceQMs(rawBoard1)
	cleanBoard2 := replaceQMs(rawBoard2)
	boardScore1 := scoreBoard(cleanBoard1)
	boardScore2 := scoreBoard(cleanBoard2)
	log.Printf("Sample raw boards: \n%v\n%v\n", rawBoard1, rawBoard2)
	log.Printf("Sample clean boards: \n%s %d\n%s %d\n", cleanBoard1, boardScore1, cleanBoard2, boardScore2)

	mutateBoard(rawBoard1)
	cleanMutant1 := replaceQMs(rawBoard1)
	mutantScore1 := scoreBoard(cleanMutant1)
	offspring := generateOffspring(rawBoard1, rawBoard2)
	cleanOffspring := replaceQMs(offspring)
	offspringScore := scoreBoard(cleanOffspring)
	log.Printf("Sample mutant: \n%v %d\n", cleanMutant1, mutantScore1)
	log.Printf("Sample offspring: \n%v %d\n", cleanOffspring, offspringScore)

	log.Println("Done")
}
