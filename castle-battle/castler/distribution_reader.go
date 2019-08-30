package castler

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

func ReadDistributions(f io.Reader) ([]*SoldierDistribution, error) {
	// f, err := os.Open(file)
	// if err != nil {
	// 	return nil, err
	// }

	// defer f.Close()

	r := csv.NewReader(f)

	// throw away the first line
	_, err := r.Read()
	if err != nil {
		return nil, err
	}

	// parse strings into numbers
	sds := make([]*SoldierDistribution, 0)
	var nrec [10]int

LOOP:
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for i, r := range record[:10] {
			var err2 error
			nrec[i], err2 = strconv.Atoi(r)
			if err2 != nil || nrec[i] < 0 {
				fmt.Printf("Smart-ass thought the rules didn't apply: %v\n", record)
				continue LOOP
			}
		}

		sds = append(sds, NewSoldierDistribution(nrec))
	}

	return sds, nil
}
