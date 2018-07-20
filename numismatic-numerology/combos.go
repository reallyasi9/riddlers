package main

func combos(n, k int) <-chan []int {
	out := make(chan []int)
	go genCombo(n, k, out)
	return out
}

func genCombo(n, k int, c chan<- []int) {

	if n <= 0 || k <= 0 || k > n {
		close(c)
		return
	}

	pool := intRange(n)
	indices := intRange(k)
	result := intRange(k)
	resultCopy := copyInc(result)
	c <- resultCopy
	for {
		i := k - 1
		for ; i >= 0 && indices[i] == i+len(pool)-k; i-- {
		}
		if i <= 0 { // Always keep 0
			break
		}
		indices[i]++
		for j := i + 1; j < k; j++ {
			indices[j] = indices[j-1] + 1
		}
		for ; i < len(indices); i++ {
			result[i] = pool[indices[i]]
		}
		resultCopy := copyInc(result)
		c <- resultCopy
	}

	close(c)
}

func intRange(k int) []int {
	out := make([]int, k)
	for i := range out {
		out[i] = i
	}
	return out
}

func copyInc(k []int) []int {
	r := make([]int, len(k))
	for i, v := range k {
		r[i] = v + 1
	}
	return r
}
