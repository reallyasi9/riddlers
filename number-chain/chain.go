package main

import (
	"fmt"
	"os"
	"strconv"
)

type numberNet struct {
	n       int
	adjList map[int][]int
}

func makeNumberNet(n int) *numberNet {
	nn := new(numberNet)
	nn.n = n
	nn.adjList = multiplesAndFactors(n)
	for i := 1; i <= n; i++ {
		nn.adjList[0] = append(nn.adjList[0], i) // start here
	}
	return nn
}

type dfsHelper struct {
	nn      *numberNet
	visited map[int]bool
}

func makeDFSHelper(nnt *numberNet) *dfsHelper {
	h := new(dfsHelper)
	h.nn = nnt
	h.visited = make(map[int]bool)
	for key := range nnt.adjList {
		h.visited[key] = false
	}
	return h
}

func main() {
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	nn := makeNumberNet(n)
	fmt.Printf("Multiples and Factors:\n%v\n", nn.adjList)

	dh := makeDFSHelper(nn)

	longest := make([]int, 0)
	checked := make([]bool, n+1)

	for i := 1; i <= n; i++ {
		if checked[i] {
			continue
		}

		test := dh.dfs(i)
		if len(test) > len(longest) {
			longest = make([]int, len(test))
			copy(longest, test)
			fmt.Printf("Best so far: %v\n", longest)
		}
		checked[i] = true
		checked[longest[0]] = true // elements are pushed from the front
	}

	fmt.Printf("Longest found: %v\n", longest)
	fmt.Printf("(length %d)\n", len(longest))
}

func multiplesAndFactors(n int) map[int][]int {
	m := make(map[int][]int)
	for i := 1; i <= n; i++ {
		m[i] = allFactors(i)
		m[i] = append(m[i], allMultiples(i, n)...)
	}
	return m
}

func allFactors(n int) []int {
	f := make([]int, 0)
	for i := 1; i < n; i++ {
		if (n % i) == 0 {
			f = append(f, i)
		}
	}
	return f
}

func allMultiples(n int, max int) []int {
	f := make([]int, 0)
	for i := 2 * n; i <= max; i += n {
		f = append(f, i)
	}
	return f
}

func (dh *dfsHelper) dfs(n int) []int {
	dh.visited[n] = true

	best := make([]int, 0)

	for _, m := range dh.nn.adjList[n] {
		if dh.visited[m] {
			continue
		}
		test := dh.dfs(m)
		if len(test) > len(best) {
			best = make([]int, len(test))
			copy(best, test)
		}
	}

	dh.visited[n] = false
	return append(best, n)
}
