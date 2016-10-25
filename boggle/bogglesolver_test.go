package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
	err := setup(filepath.Join("dictionaries", "dictionary-yawl.txt"))
	if err != nil {
		t.Fatal(err)
	}

	testPoints := []int{
		0, 1, 2, 3, 4, 5, 100, 200, 300, 400, 500, 750, 1000, 1250, 1500, 2000, 4410, 4527, 4540, 13464, 26539,
	}

	for _, pts := range testPoints {
		fn := filepath.Join("test", fmt.Sprintf("board-points%d.txt", pts))
		board, err := ReadBoggleBoard(fn)
		if err != nil {
			t.Fatal(err)
		}

		s := board.score()
		if s != pts {
			t.Errorf("score %d != expected %d", s, pts)
		}
	}
}

func BenchmarkBoggleSolver(b *testing.B) {
	err := setup(filepath.Join("dictionaries", "enable1.txt"))
	if err != nil {
		b.Fatal(err)
	}
	board := NewBoggleBoard()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		board.score()
	}
}
