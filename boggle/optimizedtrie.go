package main

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
