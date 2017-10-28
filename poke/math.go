package poke

import (
	"math"
	"strconv"
)

func fac(n int) int {
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return n * fac(n-1)
	}
}

func nChooseM(n, m int) int {
	if m > n {
		panic("Bad input")
	}
	if m == 0 || m == n {
		return 1
	}
	return fac(n) / (fac(m) * fac(n-m))
}

func padWithZeros(a string, n int) string {
	m := n - len(a)
	if m <= 0 {
		return a
	} else {
		s := ""
		for i := 0; i < m; i++ {
			s += "0"
		}
		return s + a
	}
}

func hammingWeight(a string) int {
	total := 0
	for _, c := range a {
		if c == '1' {
			total += 1
		}
	}
	return total
}

func transformToBool(a string) []bool {
	res := []bool{}
	for _, c := range a {
		if c == '1' {
			res = append(res, true)
		} else {
			res = append(res, false)
		}
	}
	return res
}

// Returns n choose m combinations.
// Array will have many n-len arrays with m True's
func GenerateCombinations(n, m int) [][]bool {
	// b := nChooseM(n, m)
	// fmt.Println(b)
	all := [][]bool{}
	for i := 0; i < int(math.Pow(2, float64(n))); i++ {
		binary := strconv.FormatInt(int64(i), 2)
		if hammingWeight(binary) == m {
			padded := padWithZeros(binary, n)
			boolVec := transformToBool(padded)
			all = append(all, boolVec)
		}
	}
	return all
}
