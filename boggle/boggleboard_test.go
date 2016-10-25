package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestBoggleBoard(t *testing.T) {
	entries, err := ioutil.ReadDir("test")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fn := filepath.Join("test", e.Name())
		_, err := ReadBoggleBoard(fn)
		if err != nil {
			t.Fatal(err)
		}
		// buf, err := b.MarshalText()
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// bufgold, err := ioutil.ReadFile(fn)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// if len(buf) != len(bufgold) {
		// 	t.Errorf("length of file %s (%d) not equal to produced board (%d)", fn, len(bufgold), len(buf))
		// 	t.FailNow()
		// }
		// for i := range bufgold {
		// 	if buf[i] != bufgold[i] {
		// 		t.Errorf("rune %d in file %s (%c) != produced rune (%c)", i, fn, bufgold[i], buf[i])
		// 	}
		// }
	}
}

func BenchmarkShuffle(b *testing.B) {
	dictfile := filepath.Join("dictionaries", "dictionary-enable1.txt")
	adjList := buildAdjList(4, 4)
	f2, err := frequencyCount(dictfile, 16)
	if err != nil {
		b.Fatal(err)
	}
	bb := newDiceBoard(4, 4, boggle1992)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bb.DictShuffle(adjList, f2)
	}
}
