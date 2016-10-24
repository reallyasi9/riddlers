package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	t := DAGTrie{}

	file, err := os.Open("enable1.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t.Insert(strings.ToUpper(scanner.Text()))
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Try to find some more
	len := 0
	for l := t.LongestChain(2); l.length >= len; l = t.LongestChain(2) {
		fmt.Printf("Longest: %s\n", l.value)
		chain := t.Trace(l.value, 2)
		fmt.Printf("How to get there: %v\n", chain)
		len = l.length
		t.Delete(l.value)
	}

	// Visualize
	err = t.DumpGraphML("test.graphml")
	if err != nil {
		log.Fatal(err)
	}
}
