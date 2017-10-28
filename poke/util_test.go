package poke

import "testing"

type Case struct {
	n   int
	m   int
	res int
}

type Zero struct {
	in  string
	n   int
	out string
}

type Hamming struct {
	in  string
	out int
}

func TestNChooseM(t *testing.T) {
	cases := []Case{
		Case{10, 0, 1},
		Case{10, 10, 1},
		Case{10, 1, 10},
		Case{3, 2, 3},
		Case{14, 4, 1001},
	}

	for _, c := range cases {
		if nChooseM(c.n, c.m) != c.res {
			t.Errorf("Wrong: nChooseM(%d, %d) = %d", c.n, c.m, c.res)
		}
	}
}

func TestPadWithZeros(t *testing.T) {
	cases := []Zero{
		Zero{"1010", 5, "01010"},
		Zero{"1010", 10, "0000001010"},
		Zero{"1010", 4, "1010"},
		Zero{"1111", 6, "001111"},
		Zero{"0000", 6, "000000"},
	}

	for _, c := range cases {
		if padWithZeros(c.in, c.n) != c.out {
			t.Errorf("Wrong: padWithZeros(%s, %d) = %s", c.in, c.n, c.out)
		}
	}
}

func TestHammingWeight(t *testing.T) {
	cases := []Hamming{
		Hamming{"10001", 2},
		Hamming{"11111", 5},
		Hamming{"11001", 3},
		Hamming{"10011110010", 6},
		Hamming{"0", 0},
	}

	for _, c := range cases {
		if hammingWeight(c.in) != c.out {
			t.Errorf("Wrong: hammingWeight(%s) = %d", c.in, c.out)
		}
	}
}
