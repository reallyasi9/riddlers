package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// Boggler is an interface to a boggle board
type Boggler interface {
	Rows() int
	Cols() int
	Get(int, int) rune
	GetLinear(int) rune
}

// BoggleBoard is a struct that defines a Boggle playspace
type BoggleBoard struct {
	rows  int
	cols  int
	board [][]rune
}

// Rows implements Boggler's interface
func (bb *BoggleBoard) Rows() int {
	return bb.rows
}

// Cols implements Boggler's interface
func (bb *BoggleBoard) Cols() int {
	return bb.cols
}

// Get implements Boggler's interface
func (bb *BoggleBoard) Get(i int, j int) rune {
	return bb.board[i][j]
}

// GetLinear implements Boggler's interface
func (bb *BoggleBoard) GetLinear(k int) rune {
	return bb.board[k/bb.cols][k%bb.cols]
}

// DiceBoard stores the actual dice used to make the board (instead of storing the runes)
type DiceBoard struct {
	rows int
	cols int
	dice []string
	die  [][]int
	face [][]int
}

// Rows implements Boggler's interface
func (bb *DiceBoard) Rows() int {
	return bb.rows
}

// Cols implements Boggler's interface
func (bb *DiceBoard) Cols() int {
	return bb.cols
}

// Get implements Boggler's interface
func (bb *DiceBoard) Get(i int, j int) rune {
	return rune(bb.dice[bb.die[i][j]][bb.face[i][j]])
}

// GetLinear implements Boggler's interface
func (bb *DiceBoard) GetLinear(k int) rune {
	return bb.Get(k/bb.cols, k%bb.cols)
}

// the 16 Boggle dice (1992 version)
var boggle1992 = []string{
	"LRYTTE", "VTHRWE", "EGHWNE", "SEOTIS",
	"ANAEEG", "IDSYTT", "OATTOW", "MTOICU",
	"AFPKFS", "XLDERI", "HCPOAS", "ENSIEU",
	"YLDEVR", "ZNRNHL", "NMIQHU", "OBBAOJ",
}

// the 16 Boggle dice (1983 version)
var boggle1983 = []string{
	"AACIOT", "ABILTY", "ABJMOQ", "ACDEMP",
	"ACELRS", "ADENVZ", "AHMORS", "BIFORX",
	"DENOSW", "DKNOTU", "EEFHIY", "EGINTV",
	"EGKLUY", "EHINPS", "ELPSTU", "GILRUW",
}

// the 25 Boggle Master / Boggle Deluxe dice
var boggleMaster = []string{
	"AAAFRS", "AAEEEE", "AAFIRS", "ADENNN", "AEEEEM",
	"AEEGMU", "AEGMNN", "AFIRSY", "BJKQXZ", "CCNSTW",
	"CEIILT", "CEILPT", "CEIPST", "DDLNOR", "DHHLOR",
	"DHHNOT", "DHLNOR", "EIIITT", "EMOTTT", "ENSSSU",
	"FIPRSY", "GORRVW", "HIPRRY", "NOOTUW", "OOOTTU",
}

// the 25 Big Boggle dice
var boggleBig = []string{
	"AAAFRS", "AAEEEE", "AAFIRS", "ADENNN", "AEEEEM",
	"AEEGMU", "AEGMNN", "AFIRSY", "BJKQXZ", "CCENST",
	"CEIILT", "CEILPT", "CEIPST", "DDHNOT", "DHHLOR",
	"DHLNOR", "DHLNOR", "EIIITT", "EMOTTT", "ENSSSU",
	"FIPRSY", "GORRVW", "IPRRRY", "NOOTUW", "OOOTTU",
}

// letters in the English alphabet
const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var frequencies = []float64{
	0.08167, 0.01492, 0.02782, 0.04253, 0.12703, 0.02228,
	0.02015, 0.06094, 0.06966, 0.00153, 0.00772, 0.04025,
	0.02406, 0.06749, 0.07507, 0.01929, 0.00095, 0.05987,
	0.06327, 0.09056, 0.02758, 0.00978, 0.02360, 0.00150,
	0.01974, 0.00074,
}

func inPlaceShuffle(a *[]string) {
	for i := range *a {
		j := rand.Intn(i + 1)
		(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
	}
}

func shuffledInts(n int) []int {
	a := make([]int, n)
	for i := range a {
		a[i] = i
	}
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func newBoggleBoard(rows int, cols int, dice []string) *BoggleBoard {
	inPlaceShuffle(&dice)
	board := make([][]rune, rows)
	for i := 0; i < rows; i++ {
		board[i] = make([]rune, cols)
		for j := 0; j < cols; j++ {
			letters := dice[cols*i+j]
			r := rand.Intn(len(letters))
			board[i][j] = rune(letters[r])
		}
	}
	return &BoggleBoard{rows: rows, cols: cols, board: board}
}

func newDiceBoard(rows int, cols int, dice []string) *DiceBoard {
	diceorder := shuffledInts(len(dice))
	die := make([][]int, rows)
	face := make([][]int, rows)
	l := len(dice[0])
	for i := 0; i < rows; i++ {
		die[i] = make([]int, cols)
		face[i] = make([]int, len(dice[0]))
		for j := 0; j < cols; j++ {
			die[i][j] = diceorder[cols*i+j]
			face[i][j] = rand.Intn(l)
		}
	}
	return &DiceBoard{rows: rows, cols: cols, dice: dice, die: die, face: face}
}

// NewBoggleBoard initializes a random 4-by-4 board by rolling the Hasbro dice.
func NewBoggleBoard() *BoggleBoard {
	return newBoggleBoard(4, 4, boggle1992)
}

// NewBoggleBoard1983 initializes a random 4-by-4 board by rolling the 1983 Hasbro dice.
// This function is not threadsafe.
func NewBoggleBoard1983() *BoggleBoard {
	return newBoggleBoard(4, 4, boggle1983)
}

// NewBoggleBoardMaster initializes a random 5-by-5 board by rolling the Boggle Master/Boggle Deluxe dice.
// This function is not threadsafe.
func NewBoggleBoardMaster() *BoggleBoard {
	return newBoggleBoard(5, 5, boggleMaster)
}

// NewBoggleBoardBig initializes a random 5-by-5 board by rolling the Big Boggle dice.
// This function is not threadsafe.
func NewBoggleBoardBig() *BoggleBoard {
	return newBoggleBoard(5, 5, boggleBig)
}

func (bb *BoggleBoard) String() string {
	var bf bytes.Buffer
	bf.WriteString(fmt.Sprintf("%d %d\n", bb.rows, bb.cols))
	for _, br := range bb.board {
		for _, bc := range br {
			bf.WriteRune(bc)
			if bc == 'Q' {
				bf.WriteString("u ")
			} else {
				bf.WriteString("  ")
			}
		}
		bf.WriteString("\n")
	}
	return strings.TrimSpace(bf.String())
}

func (bb *DiceBoard) String() string {
	var bf bytes.Buffer
	bf.WriteString(fmt.Sprintf("%d %d\n", bb.rows, bb.cols))
	for i := 0; i < bb.Rows(); i++ {
		for j := 0; j < bb.Cols(); j++ {
			bc := bb.Get(i, j)
			bf.WriteRune(bc)
			if bc == 'Q' {
				bf.WriteString("u ")
			} else {
				bf.WriteString("  ")
			}
		}
		bf.WriteString("\n")
	}
	return strings.TrimSpace(bf.String())
}

// MarshalText implements text marshalling for TextMarshaler interface.
func (bb *BoggleBoard) MarshalText() ([]byte, error) {
	return []byte(bb.String()), nil
}

// MarshalText implements text marshalling for TextMarshaler interface.
func (bb *DiceBoard) MarshalText() ([]byte, error) {
	return []byte(bb.String()), nil
}

// UnmarshalText implements text unmarshaling for TextMarshaler interface.
func (bb *BoggleBoard) UnmarshalText(text []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(text))
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	rows, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return err
	}
	scanner.Scan()
	cols, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return err
	}

	board := make([][]rune, rows)
	for i := 0; i < rows; i++ {
		board[i] = make([]rune, cols)
		for j := 0; j < cols; j++ {
			if !scanner.Scan() {
				return errors.New("ran out of letters when scanning text")
			}

			letter := strings.ToUpper(scanner.Text())

			if len(letter) != 1 {
				return fmt.Errorf("invalid character: %s", letter)
			}
			if strings.Index(alphabet, letter) == -1 {
				return fmt.Errorf("invalid character: %s", letter)
			}

			if letter == "QU" {
				board[i][j] = 'Q'
			} else {
				board[i][j] = rune(letter[0])
			}
		}
	}

	bb.rows = rows
	bb.cols = cols
	bb.board = board

	return nil
}

// ReadBoggleBoard reads a board from a file
func ReadBoggleBoard(filename string) (*BoggleBoard, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bb BoggleBoard
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = bb.UnmarshalText(bs)
	if err != nil {
		return nil, err
	}

	return &bb, nil
}

func randomIndex(p []float64) (int, error) {
	sum := 0.
	for _, v := range p {
		if v < 0. {
			return -1, fmt.Errorf("value %f < 0", v)
		}
		sum += v
	}
	for {
		r := rand.Float64()
		s2 := 0.
		for i, v := range p {
			s2 += v / sum
			if s2 > r {
				return i, nil
			}
		}
	}
}

// NewBoggleBoardRandom creates a random M-by-N board according to the frequency of letters in the English language
func NewBoggleBoardRandom(rows int, cols int) *BoggleBoard {
	board := make([][]rune, rows)
	for i := 0; i < rows; i++ {
		board[i] = make([]rune, cols)
		for j := 0; j < cols; j++ {
			idx, _ := randomIndex(frequencies)
			board[i][j] = rune(alphabet[idx])
		}
	}

	return &BoggleBoard{rows: rows, cols: cols, board: board}
}

// DictShuffle shuffles the dice according to the 2-letter occurance frequencies
func (bb *DiceBoard) DictShuffle(adjList [][]int, f2 [][]float64) {
	weights := make([]float64, bb.rows*bb.cols)
	for i, adjl := range adjList {
		for _, adj := range adjl {
			r1 := bb.GetLinear(i)
			r2 := bb.GetLinear(adj)
			weights[i] += f2[r1-'A'][r2-'A']
		}
		weights[i] = 1. - weights[i]/float64(len(adjl))
	}

	i1, err := randomIndex(weights)
	if err != nil {
		panic(err)
	}
	i2, err := randomIndex(weights)
	if err != nil {
		panic(err)
	}

	r1 := i1 / bb.Cols()
	c1 := i1 % bb.Cols()
	r2 := i2 / bb.Cols()
	c2 := i2 % bb.Cols()

	// Flip?
	if rand.Float32() < .5 {
		bb.die[r1][c1], bb.die[r2][c2] = bb.die[r2][c2], bb.die[r1][c1]
	}

	// Roll
	l := len(bb.dice[0])
	bb.face[r1][c1] = rand.Intn(l)
	bb.face[r2][c2] = rand.Intn(l)
}

// NewBoggleBoardArray Initialize board from the given 2D character array.
func NewBoggleBoardArray(board [][]rune) (*BoggleBoard, error) {
	rows := len(board)
	bb := make([][]rune, rows)
	cols := len(board[0])
	for i, bc := range board {
		if len(bc) != cols {
			return nil, errors.New("array is ragged")
		}
		bb[i] = make([]rune, cols)
		for j, br := range bc {
			if strings.Index(alphabet, string(br)) == -1 {
				return nil, fmt.Errorf("invalid character: %c", bc)
			}
			bb[i][j] = rune(br)
		}
	}
	return &BoggleBoard{rows: rows, cols: cols, board: bb}, nil
}
