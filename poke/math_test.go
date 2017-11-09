package poke

import "testing"

const TWO_TO_THE_27 = 134217728

type Case struct {
	n   int
	m   int
	res int
}

type SetBitCase struct {
	n   int
	r   int
}

func TestGenerateCombinations(t *testing.T) {
	cases := []Case{
		Case{4, 4, 1},
		Case{14, 4, 1001},
		Case{17, 4, 2380},
		Case{27, 4, 17550},
	}

	for _, c := range cases {
		x := len(GenerateCombinations(c.n, c.m))
		if x != c.res {
			t.Errorf("Wrong: len-generatecomb(%d, %d) != %d", c.n, c.m, x)
		}
	}
}

func TestCountSetBits(t *testing.T) {
	cases := []SetBitCase{
		SetBitCase{2,  1},
		SetBitCase{1,  1},
		SetBitCase{3,  2},
		SetBitCase{20, 2},
	}

	for _, c := range cases {
		if countSetBits(c.n) != c.r {
			t.Errorf("Wrong: countsetbits(%d) != %d", c.n, c.r)
		}
	}
}

func TestCreateBoolVec(t *testing.T) {
	res := createBoolVector(2, 2)
	if res[0] != true || res[1] != false {
		t.Errorf("Wrong: createvec\n")
	}
}

func BenchmarkCountSetBits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		countSetBits(TWO_TO_THE_27)
	}
}

func BenchmarkCreateBoolVec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createBoolVector(TWO_TO_THE_27 - 1, 27)
	}
}

func BenchmarkGenerateCombinations(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateCombinations(27, 4)
	}
}
