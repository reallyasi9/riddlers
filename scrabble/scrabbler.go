package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/armon/go-radix"
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

func scoreBoard(board string) int {
	score := 0
	foundTrie := radix.New()
	i := 0
	j := 1
	for i < len(board) {
		log.Printf("Indices %d %d\n", i, j)
		sub := board[i:j]
		log.Printf("Checking %s\n", sub)
		val, ok := scoreTrie.Get(sub)
		log.Printf("Found? %v %v\n", ok, val)
		if ok {
			if _, alreadyFound := foundTrie.Get(sub); !alreadyFound {
				// Is a word: score it
				log.Println("Score!")
				score += val.(int)
				foundTrie.Insert(sub, true)
			}
		}

		if _, _, isPrefix := scoreTrie.LongestPrefix(sub); isPrefix {
			log.Printf("%s is a prefix", sub)
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

func main() {
	log.Println("Starting")

	log.Printf("Rune Scores: %v\n", runeScores)

	err := buildDictionary(dictionaryURL)
	if err != nil {
		log.Fatalln(err)
	}

	rawBoard := generateBoard()
	cleanBoard := replaceQMs(rawBoard)
	boardScore := scoreBoard(cleanBoard)
	log.Printf("Sample raw board: %v\n", rawBoard)
	log.Printf("Sample clean board: %s\n", cleanBoard)
	log.Printf("Score: %d\n", boardScore)

	log.Println("Done")
}
