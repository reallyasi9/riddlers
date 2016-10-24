package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// DAGTrie is a combination prefix and suffix tree.
type DAGTrie struct {
	prefixNodes  [26]*DAGTrie
	suffixNodes  [26]*DAGTrie
	prefixParent *DAGTrie
	suffixParent *DAGTrie
	word         bool
	length       int
	value        string
}

type stack []*DAGTrie

func (s *stack) Push(v *DAGTrie) {
	*s = append(*s, v)
}

func (s *stack) PushAll(v []*DAGTrie) {
	for _, w := range v {
		if w != nil {
			*s = append(*s, w)
		}
	}
}

func (s *stack) Pop() *DAGTrie {
	if len(*s) == 0 {
		return nil
	}
	ret := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return ret
}

type stringStack []string

func (s *stringStack) Push(v string) {
	*s = append(*s, v)
}

func (s *stringStack) Pop() string {
	if len(*s) == 0 {
		return ""
	}
	ret := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return ret
}

// Insert a string into the trie.  Must be an ASCII uppercase string, else everything breaks!
func (trie *DAGTrie) Insert(s string) {
	endTrie := trie.insertPrefix(s)
	trie.backfillSuffix(endTrie, s)
}

// Delete by simply marking the given node as not a word
func (trie *DAGTrie) Delete(s string) {
	node := trie.Get(s)
	if node != nil {
		node.word = false
		node.value = ""
	}
}

// Dump produces a string representation of the trie
func (trie *DAGTrie) Dump() string {
	return fmt.Sprintf("*  (%p)\n%s", trie, trie.dumpLevel(0))
}

// dumpLevel dumps  a level of a trie with a given indentation
func (trie *DAGTrie) dumpLevel(indent int) string {
	var b bytes.Buffer
	for i, v := range trie.prefixNodes {
		if v == nil {
			continue
		}
		b.WriteString(strings.Repeat(" ", indent*3))
		b.WriteString("-> ")
		b.WriteRune(rune('A' + i))
		if v.word {
			b.WriteString("+ ")
		} else {
			b.WriteString("  ")
		}
		b.WriteString(fmt.Sprintf("(%p)\n", v))
		b.WriteString(v.dumpLevel(v.length))
	}
	for i, v := range trie.suffixNodes {
		if v == nil {
			continue
		}
		b.WriteString(strings.Repeat(" ", indent*3))
		b.WriteString("<- ")
		b.WriteRune(rune('A' + i))
		if v.word {
			b.WriteString("+ ")
		} else {
			b.WriteString("  ")
		}
		b.WriteString(fmt.Sprintf("(%p)\n", v))
		b.WriteString(v.dumpLevel(v.length))
	}
	return b.String()
}

// Get an entry in the trie if it exists, else nil
func (trie *DAGTrie) Get(s string) *DAGTrie {
	last := trie
	for _, r := range s {
		val := r - 'A'
		last = last.prefixNodes[val]
		if last == nil {
			return nil
		}
	}
	return last
}

func (trie *DAGTrie) buildPrefix(s string) *DAGTrie {
	last := trie
	for i, r := range s {
		val := r - 'A'
		if last.prefixNodes[val] == nil {
			last.prefixNodes[val] = &DAGTrie{length: i + 1, prefixParent: last}
		}
		last = last.prefixNodes[val]
	}
	return last
}

func (trie *DAGTrie) insertPrefix(s string) *DAGTrie {
	last := trie.buildPrefix(s)
	last.word = true
	last.value = s
	return last
}

func (trie *DAGTrie) backfillSuffix(endTrie *DAGTrie, s string) {
	if len(s) <= 1 {
		return
	}

	val := s[0] - 'A'
	remainder := s[1:]
	suffix := trie.buildPrefix(remainder)
	suffix.suffixNodes[val] = endTrie
	endTrie.suffixParent = suffix
	trie.backfillSuffix(suffix, remainder)
}

// LongestChain uses a depth-first search to find the longest word made from a chain of prefix/suffixes in the trie
func (trie *DAGTrie) LongestChain(startLength int) *DAGTrie {
	s := make(stack, 0)
	s.PushAll(trie.prefixNodes[:])
	s.PushAll(trie.suffixNodes[:])

	var longest *DAGTrie

	for v := s.Pop(); v != nil; v = s.Pop() {
		if v.length < startLength || v.word {
			s.PushAll(v.prefixNodes[:])
			s.PushAll(v.suffixNodes[:])
		}

		if v.word && (longest == nil || v.length >= longest.length) {
			longest = v
		}

	}

	return longest
}

// Trace tells you how to get from the root of the trie to the given node following words all the way down
func (trie *DAGTrie) Trace(s string, startLength int) []string {
	var ss stringStack

	last := trie.Get(s)
	ss.Push(last.value)

	if last.prefixParent != nil && last.prefixParent.findRoot(&ss, startLength) {
		return ss
	} else if last.suffixParent != nil && last.suffixParent.findRoot(&ss, startLength) {
		return ss
	}
	return nil
}

func (trie *DAGTrie) findRoot(ss *stringStack, startLength int) bool {
	if !trie.word {
		return false
	}
	ss.Push(trie.value)
	if trie.length <= startLength {
		return true
	}
	if trie.prefixParent != nil && trie.prefixParent.findRoot(ss, startLength) {
		return true
	} else if trie.suffixParent != nil && trie.suffixParent.findRoot(ss, startLength) {
		return true
	}
	ss.Pop()
	return false
}

// DumpGraphML will dump a trie to GML format for visualization
func (trie *DAGTrie) DumpGraphML(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<graphml xmlns="http://graphml.graphdrawing.org/xmlns"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns
http://graphml.graphdrawing.org/xmlns/1.0/graphml.xsd">
`)

	file.WriteString("<graph id=\"trie\" edgedefault=\"directed\">\n")
	ss := make(stack, 0)
	trie.collectNodes(&ss)
	// Nodes first
	for _, v := range ss {
		fmt.Fprintf(file, "<node id=\"%p\">", v)
		fmt.Fprintf(file, "<data key=\"word\">%t</data>", v.word)
		fmt.Fprintf(file, "<data key=\"value\">%s</data>", v.value)
		fmt.Fprintf(file, "</node>\n")
	}
	// Edges second
	for i, v := range ss {
		for j, e := range v.prefixNodes {
			if e == nil {
				continue
			}
			fmt.Fprintf(file, "<edge id=\"e%d-%d\" source=\"%p\" target=\"%p\">", i, j, v, e)
			fmt.Fprintf(file, "<data key=\"type\">p</data><data key=\"letter\">%c</data>", j+'A')
			fmt.Fprint(file, "</edge>\n")
		}
		for j, e := range v.suffixNodes {
			if e == nil {
				continue
			}
			fmt.Fprintf(file, "<edge id=\"e%d-%d\" source=\"%p\" target=\"%p\">", i, j+26, v, e)
			fmt.Fprintf(file, "<data key=\"type\">s</data><data key=\"letter\">%c</data>", j+'A')
			fmt.Fprint(file, "</edge>\n")
		}
	}
	file.WriteString("</graph>\n</graphml>\n")
	return nil
}

func (trie *DAGTrie) collectNodes(ss *stack) {
	for _, v := range trie.prefixNodes {
		if v == nil {
			continue
		}
		ss.Push(v)
		v.collectNodes(ss)
	}
}
