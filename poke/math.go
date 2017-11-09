package poke

import (
	"math"
)

func countSetBits(n int) int {
	c := 0
	for ; n > 0; n = n & (n - 1) {
		c += 1
	}
	return c
}

func createBoolVector(i, n int) []bool {
	b := make([]bool, n)

	j := 0
	for ; i > 0; i = i >> 1 {
		b[n-1-j] = i & 1 == 1
		j++
	}
	return b
}

// Returns n choose m combinations.
// Array will have many n-len arrays with m True's
func GenerateCombinations(n, m int) [][]bool {
	all := [][]bool{}
	nn := int(math.Pow(2, float64(n)))
	for i := 0; i < nn; i++ {
		if countSetBits(i) == m {
			boolVec := createBoolVector(i, n)
			all = append(all, boolVec)
		}
	}
	return all
}
