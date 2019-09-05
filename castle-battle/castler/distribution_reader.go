package castler

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
)

func ReadDistributions(f io.Reader) ([]*SoldierDistribution, []string, error) {
	// f, err := os.Open(file)
	// if err != nil {
	// 	return nil, err
	// }

	// defer f.Close()

	r := csv.NewReader(f)

	// first line is the header
	header, err := r.Read()
	if err != nil {
		return nil, nil, err
	}

	// parse strings into numbers
	sds := make([]*SoldierDistribution, 0)
	var nrec [NCastles]int

LOOP:
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		for i, r := range record[:NCastles] {
			var err2 error
			nrec[i], err2 = strconv.Atoi(r)
			if err2 != nil || nrec[i] < 0 {
				log.Printf("Smart-ass thought the rules didn't apply: %v\n", record)
				continue LOOP
			}
		}

		sds = append(sds, NewSoldierDistribution(nrec))
	}

	return sds, header, nil
}
