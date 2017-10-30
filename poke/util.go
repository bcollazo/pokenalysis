package poke

import (
	"sort"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
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

type BattlePokemonFloat struct {
	Key   BattlePokemon
	Value float64
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

func GetSortedPokemon(totalEkts map[string]float64, versions map[string]BattlePokemon, dir int) []BattlePokemon {
	var ss []BattlePokemonFloat
	for name, totalEkt := range totalEkts {
		ss = append(ss, BattlePokemonFloat{versions[name], totalEkt})
	}

	if dir < 0 {
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})
	} else if dir > 0 {
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value < ss[j].Value
		})
	} else { // dir == 0, normal ordering.
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Key.Id < ss[j].Key.Id
		})
	}

	var res []BattlePokemon
	for _, pair := range ss {
		res = append(res, pair.Key)
	}
	return res
}
