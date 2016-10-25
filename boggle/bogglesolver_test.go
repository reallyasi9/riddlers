package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestBoggleSolver(t *testing.T) {
	dictfile := filepath.Join("dictionaries", "dictionary-yawl.txt")

	testPoints := []int{
		0, 1, 2, 3, 4, 5, 100, 200, 300, 400, 500, 750, 1000, 1250, 1500, 2000, 4410, 4527, 4540, 13464, 26539,
	}

	for _, pts := range testPoints {
		fn := filepath.Join("test", fmt.Sprintf("board-points%d.txt", pts))
		board, err := ReadBoggleBoard(fn)
		if err != nil {
			t.Fatal(err)
		}

		bs, err := newSolver(board.rows, board.cols, dictfile)
		if err != nil {
			t.Fatal(err)
		}

		s, _ := bs.score(board)
		//fmt.Printf("%s\n", board)
		if s != pts {
			t.Errorf("score %d != expected %d", s, pts)
		}
	}
}

func BenchmarkBoggleSolver(b *testing.B) {
	dictfile := filepath.Join("dictionaries", "dictionary-enable1.txt")
	bs, err := newSolver(4, 4, dictfile)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// board := NewBoggleBoard()
		board := newDiceBoard(4, 4, boggle1992)
		b.StartTimer()
		bs.score(board)
	}
}

func BenchmarkBoggleSolverRandom(b *testing.B) {
	dictfile := filepath.Join("dictionaries", "dictionary-enable1.txt")
	bs, err := newSolver(4, 4, dictfile)
	if err != nil {
		b.Fatal(err)
	}
	board := newDiceBoard(4, 4, boggle1992)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bs.score(board)
	}
}
