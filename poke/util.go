package poke

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// Files and System
func EnsureDir(path string) {
	_ = os.Mkdir(path, 0700)
}

type DataPoint struct {
	Type   string `json:"type"`
	Number int    `json:"value"`
}

type TwoDataPoint struct {
	Type string `json:"type"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

func SaveHistoFile(histo map[Type]int, name string) {
	EnsureDir("dist")

	points := []DataPoint{}
	for t, i := range histo {
		points = append(points, DataPoint{t.Name, i})
	}

	outJson, err := json.Marshal(points)
	Check(err)
	err = ioutil.WriteFile("dist/"+name+".json", outJson, 0644)
	Check(err)
}

func SaveRatioFile(histo map[Type][2]int, name string) {
	EnsureDir("dist")

	points := []TwoDataPoint{}
	for t, i := range histo {
		points = append(points, TwoDataPoint{t.Name, i[0], i[1]})
	}

	outJson, err := json.Marshal(points)
	Check(err)
	err = ioutil.WriteFile("dist/"+name+".json", outJson, 0644)
	Check(err)
}

// Inclusive of a and b.  Expects a <= b
func IntRange(a, b int) []int {
	r := make([]int, b-a+1)
	for i := range r {
		r[i] = a + i
	}
	return r
}

// ===== Sorting:
type TypeInt struct {
	Key   Type
	Value int
}

type TypeFloat struct {
	Key   Type
	Value float64
}

type IntInt struct {
	Key   int
	Value int
}

type TypeCombinationInt struct {
	Key   TypeCombination
	Value int
}

func GetSortedIntTypes(histo map[Type]int, dir int) [18]Type {
	if dir == 0 {
		return TypeArr
	}

	var ss []TypeInt
	for k, v := range histo {
		ss = append(ss, TypeInt{k, v})
	}

	if dir < 0 {
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})
	} else if dir > 0 {
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value < ss[j].Value
		})
	}

	var res [18]Type
	for i, pair := range ss {
		res[i] = pair.Key
	}
	return res
}

func GetSortedRatioTypes(histo map[Type][2]int, dir int) [18]Type {
	if dir == 0 {
		return TypeArr
	}

	// Create Pair Slice
	var ss []TypeFloat
	for k, v := range histo {
		ss = append(ss, TypeFloat{k, float64(v[0]) / float64(v[1])})
	}

	if dir < 0 {
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})
	} else if dir > 0 {
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value < ss[j].Value
		})
	}

	var res [18]Type
	for i, pair := range ss {
		res[i] = pair.Key
	}
	return res
}

func GetSortedPokemonIds(totalKt map[int]int, dir int) []int {
	var ss []IntInt
	for id, kt := range totalKt {
		ss = append(ss, IntInt{id, kt})
	}

	if dir < 0 {
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value < ss[j].Value
		})
	} else if dir > 0 {
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})
	} else { // dir == 0, normal ordering.
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Key < ss[j].Key
		})
	}

	var res []int
	for _, pair := range ss {
		res = append(res, pair.Key)
	}
	return res
}

func GetSortedIntCombis(histo map[TypeCombination]int, dir int) []TypeCombination {
	if dir == 0 {
		return TypeCombinations
	}

	var ss []TypeCombinationInt
	for k, v := range histo {
		ss = append(ss, TypeCombinationInt{k, v})
	}

	if dir < 0 {
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})
	} else if dir > 0 {
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value < ss[j].Value
		})
	}

	var res []TypeCombination
	for _, pair := range ss {
		res = append(res, pair.Key)
	}
	return res
}
