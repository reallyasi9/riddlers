package main

import (
	"bytes"
	"fmt"
)

const radix = 26

type node struct {
	val  int
	next [radix]*node
}

// OptimizedTrie an optimized symbol table trie for Boggle
type OptimizedTrie struct {
	root *node
}

// RootValue returns either the value at the root or zero
func (ot *OptimizedTrie) RootValue() int {
	if ot.root == nil {
		return 0
	}
	return ot.root.val
}

// Get returns the value associated with the given key
func (ot *OptimizedTrie) Get(key string) int {
	x := ot.get(ot.root, key, 0)
	if x == nil {
		return 0
	}
	return x.val
}

// Subtrie returns the subtrie at the given key or nil if no such key can be found
func (ot *OptimizedTrie) Subtrie(key string) *OptimizedTrie {
	n := ot.get(ot.root, key, 0)
	if n == nil {
		return nil
	}
	return &OptimizedTrie{root: n}
}

// SubtrieR returns the subtrie at the given rune key or nil if no such key can be found
func (ot *OptimizedTrie) SubtrieR(key rune) *OptimizedTrie {
	n := ot.getR(ot.root, key)
	if n == nil {
		return nil
	}
	return &OptimizedTrie{root: n}
}

// Has returns true if the key is in the trie
func (ot *OptimizedTrie) Has(key string) bool {
	return ot.get(ot.root, key, 0) != nil
}

func (ot *OptimizedTrie) getR(x *node, key rune) *node {
	if key == 'Q' {
		return ot.get(x, "QU", 0)
	}
	if x == nil {
		return nil
	}
	return x.next[key-'A']
}

func (ot *OptimizedTrie) get(x *node, key string, d int) *node {
	if x == nil {
		return nil
	}
	if d == len(key) {
		return x
	}
	c := key[d] - 'A'
	return ot.get(x.next[c], key, d+1)
}

// Insert puts a value into the trie
func (ot *OptimizedTrie) Insert(key string, val int) {
	ot.root = ot.put(ot.root, key, val, 0)
}

func (ot *OptimizedTrie) put(x *node, key string, val int, d int) *node {
	if x == nil {
		x = &node{}
	}
	if d == len(key) {
		x.val = val
		return x
	}
	c := key[d] - 'A'
	x.next[c] = ot.put(x.next[c], key, val, d+1)
	return x
}

// DFS performs a depth-first visit of all key:value paris in the Trie.
func (ot *OptimizedTrie) DFS(visitor func(string, int) error) error {
	var buf bytes.Buffer
	node := ot.root
	err := node.dfs(&buf, visitor)
	return err
}

func (n *node) dfs(buf *bytes.Buffer, visitor func(string, int) error) error {
	if n.val != 0 {
		err := visitor(buf.String(), n.val)
		if err != nil {
			return err
		}
	}
	for i, c := range n.next {
		if c == nil {
			continue
		}
		buf.WriteRune(rune('A' + i))
		err := c.dfs(buf, visitor)
		if err != nil {
			return err
		}
		buf.Truncate(buf.Len() - 1)
	}
	return nil
}

type stringifier struct {
	buf bytes.Buffer
}

func (s *stringifier) visitor(key string, val int) error {
	s.buf.WriteString(fmt.Sprintf("%s: %d\n", key, val))
	return nil
}

func (ot *OptimizedTrie) String() string {
	var s stringifier
	ot.DFS(s.visitor)
	return s.buf.String()
}
