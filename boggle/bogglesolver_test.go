package main

import (
	"bufio"
	"os"
	"testing"
)

func setup(dictfile string) error {
	file, err := os.Open(dictfile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dictionary.fill(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func TestBoggleSolver(t *testing.T) {
	err := setup("dictionary-yawl.txt")
	if err != nil {
		t.Fatal(err)
	}

	board, err := ReadBoggleBoard("test/board-points0.txt")
	if err != nil {
		t.Fatal(err)
	}

	s := board.score()
	if s != 0 {
		t.Errorf("score %d != 0", s)
	}

	board, err = ReadBoggleBoard("test/board-points4.txt")
	if err != nil {
		t.Fatal(err)
	}

	s = board.score()
	if s != 4 {
		t.Errorf("score %d != 4", s)
	}

	board, err = ReadBoggleBoard("test/board-points4540.txt")
	if err != nil {
		t.Fatal(err)
	}

	s = board.score()
	if s != 4540 {
		t.Errorf("score %d != 4540", s)
	}
}

func BenchmarkBoggleSolver(b *testing.B) {
	err := setup("enable1.txt")
	if err != nil {
		b.Fatal(err)
	}
	board := NewBoggleBoard()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		board.score()
	}
}
