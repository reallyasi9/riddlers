package main

import (
	"container/heap"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	flag "github.com/spf13/pflag"
	"golang.org/x/exp/rand"

	"./castler"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type result struct {
	distribution [castler.NCastles]int
	record       [3]int // [losses, ties, wins]
}

// implements heap.Interface and holds results.
type results []*result

func (pq results) Len() int { return len(pq) }

func (pq results) Less(i, j int) bool {
	// We want Pop to give us the lowest score for removal
	return pq[i].record[2] < pq[j].record[2]
}

func (pq results) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *results) Push(x interface{}) {
	*pq = append(*pq, x.(*result))
}

func (pq *results) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

var solutionFiles []string
var outFile string
var randomness int
var method string
var iterations int
var keep int
var seed uint64
var bootstrap bool

func init() {
	flag.StringSliceVarP(&solutionFiles, "solution", "s", nil, "previous solution file (required)")
	flag.StringVarP(&outFile, "output", "o", "", "output file (default: stdout)")
	flag.IntVarP(&randomness, "randomness", "r", 10, "strength of randomness for altering solutions (has no effect if method is 'shuffle')")
	flag.StringVarP(&method, "method", "m", "walk", "solution alteration method (possible options are 'walk', 'jitter', and 'shuffle')")
	flag.IntVarP(&iterations, "iterations", "i", 10000, "number of iterations to perform")
	flag.IntVarP(&keep, "keep", "k", 0, "number of top solutions to keep (default: same number as in the solution file(s))")
	flag.BoolVarP(&bootstrap, "bootstrap", "b", false, "bootstrap the previous solutions")
	flag.Uint64Var(&seed, "seed", 0, "random seed (default: use the clock to set the seed)")
}

func randomize(d *castler.SoldierDistribution, rng *rand.Rand) {
	switch method {
	case "walk":
		d.RandomWalk(rng, randomness)
	case "jitter":
		d.Jitter(rng, float64(randomness))
	case "shuffle":
		d.Randomize(rng)
	}
}

func main() {

	flag.Parse()

	if len(solutionFiles) == 0 {
		log.Fatal("must provide at least one solution file")
	}

	sds := make([]*castler.SoldierDistribution, 0)
	header := make([]string, 0)
	for _, fn := range solutionFiles {
		file, err := os.Open(fn)
		check(err)

		log.Printf("reading file %s\n", fn)
		var sd []*castler.SoldierDistribution
		sd, header, err = castler.ReadDistributions(file)
		check(err)

		sds = append(sds, sd...)
	}

	if keep <= 0 {
		keep = len(sds)
		log.Printf("set keep value to %d\n", keep)
	}

	if keep > iterations {
		log.Printf("not enough iterations to keep %d, keeping %d instead\n", keep, iterations)
		keep = iterations
	}

	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
		log.Printf("seed value to %d\n", seed)
	}
	rng := rand.New(rand.NewSource(seed))

	test := castler.DefaultSoldierDistribution(castler.NSoldiers)
	res := &results{}
	heap.Init(res)

	for i := 0; i < iterations; i++ {

		randomize(test, rng)

		var record [3]int
		// Bootstrap?
		for j := 0; j < len(sds); j++ {
			jbs := j
			if bootstrap {
				jbs = rng.Intn(len(sds))
			}
			record[test.MatchUp(sds[jbs])+1]++
		}

		heap.Push(res, &result{distribution: test.Castles(), record: record})
		if res.Len() > keep {
			heap.Pop(res) // don't care about worst element anymore
		}

	}

	var out *os.File
	var err error
	if outFile == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(outFile)
		check(err)
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	// clean header
	header = append(header[:castler.NCastles], "losses", "ties", "wins")
	writer.Write(header)

	outLine := make([]string, castler.NCastles+3)
	for res.Len() > 0 {
		x := heap.Pop(res)
		r := x.(*result)
		for i := 0; i < castler.NCastles; i++ {
			outLine[i] = strconv.Itoa(r.distribution[i])
		}
		for i := 0; i < 3; i++ {
			outLine[castler.NCastles+i] = strconv.Itoa(r.record[i])
		}
		err := writer.Write(outLine)
		check(err)
	}

	writer.Flush()
	check(writer.Error())

	log.Println("Done")
}
