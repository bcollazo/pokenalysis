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

func createBoolVector(n int) []bool {
	bn := int(math.Log2(float64(n))) + 1
	b := make([]bool, bn)

	i := 0
	nn := n
	for ; nn > 0; nn = nn >> 1 {
		if nn & 1 == 1 {
			b[bn-1-i] = true
		}
		i++
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
			boolVec := createBoolVector(i)
			all = append(all, boolVec)
		}
	}
	return all
}
