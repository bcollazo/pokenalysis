package poke

import (
	"fmt"
	"github.com/raphamorim/go-rainbow"
	"gopkg.in/cheggaaa/pb.v1"
	"strconv"
)

const LONGEST_TYPE_NAME_LEN = 8
const WHITE = "#FFFFFF"
const FAV_BLUE = "#42b3f4"

func strRepeat(amount int, str string) string {
	var toReturn string
	for i := 0; i < amount; i++ {
		toReturn += str
	}
	return toReturn
}

func typeLabel(t Type) string {
	n := LONGEST_TYPE_NAME_LEN - len(t.Name)
	txt := strRepeat(n, " ") + t.Name + ":"
	return rainbow.Bold(rainbow.Hex(WHITE, txt))
}

func combiLabel(c TypeCombination) string {
	n := LONGEST_TYPE_NAME_LEN*2 + 1 - len(c.FirstSlot.Name) - len(c.SecondSlot.Name)
	var txt string
	if c.SecondSlot.Name != "" {
		txt = strRepeat(n, " ") + c.FirstSlot.Name + "-" + c.SecondSlot.Name + ":"
	} else {
		txt = strRepeat(n+1, " ") + c.FirstSlot.Name + ":"
	}
	return rainbow.Bold(rainbow.Hex(WHITE, txt))
}

// Ensure all types have a 0 entry.
func emptyHisto() map[Type]int {
	histo := make(map[Type]int)
	for _, t := range TypeArr {
		histo[t] = 0
	}
	return histo
}

func printHisto(histo map[Type]int, sorted [18]Type) {
	for _, t := range sorted {
		bar := strRepeat(histo[t], "#")
		fmt.Printf("%s %s (%d)\n",
			typeLabel(t),
			rainbow.Hex(t.HexColor, bar),
			histo[t])
	}
}

func printRatios(ratios map[Type][2]int, sorted [18]Type) {
	for _, t := range sorted {
		fmt.Printf("%s %s / %s (%f)\n",
			typeLabel(t),
			rainbow.Hex(t.HexColor, strconv.Itoa(ratios[t][0])),
			rainbow.Hex(t.HexColor, strconv.Itoa(ratios[t][1])),
			float64(ratios[t][0])/float64(ratios[t][1]))
	}
}

func printCombiHisto(histo map[TypeCombination]int, sorted []TypeCombination) {
	for _, c := range sorted {
		bar := strRepeat(histo[c], "#")
		fmt.Printf("%s %s (%d)\n",
			combiLabel(c),
			rainbow.Hex(FAV_BLUE, bar),
			histo[c])
	}
}

func Histo(list []Pokemon, sortDir int) {
	histo := emptyHisto()
	for _, p := range list {
		for _, t := range p.Types {
			histo[t] += 1
		}
	}

	sortedTypes := GetSortedIntTypes(histo, sortDir)
	printHisto(histo, sortedTypes)
}

// Number of pokemons such type is good against.
func SuperEffectiveHisto(list []Pokemon, sortDir int) {
	histo := emptyHisto()

	for _, pokemon := range list {
		for _, t := range TypeArr {
			// Check if super-effective.  If so, add
			if TypeEffectiveness(t, pokemon.Types) >= 2.0 {
				histo[t] += 1
			}
		}
	}

	sortedTypes := GetSortedIntTypes(histo, sortDir)
	printHisto(histo, sortedTypes)
}

// For each type, take the ratio of
// how many pokemons are weak against it (compounded type is strong) vs
// how many pokemons are strong against it (have at least 1 super effective type)
// This does not worry about pokemon that can learn a move from another type
// and make it super effective.  e.g. a Gardevoir with Leaf Blade, makes
// Blastoise vulnerable to it, but such configuration are not consider here...
// Later, maybe we can make the 'vulnerable' definition to be pokemons
// that learn a strong (>60?) attack of a type that is super effective.
// TODO: Take type-combinations instead (+ single types too).
func GoodRatios(list []Pokemon, sortDir int) {
	ratios := make(map[Type][2]int)
	for _, t := range TypeArr {
		pokemonsItKills := 0
		pokemonsThatKillIt := 0
		for _, pokemon := range list {
			// We are good against this pokemon
			if TypeEffectiveness(t, pokemon.Types) >= 2.0 {
				pokemonsItKills += 1
			}

			// At least one of its type is good against us...
			for _, tt := range pokemon.Types {
				if TypeEffectiveness(tt, []Type{t}) >= 2.0 {
					pokemonsThatKillIt += 1
					break
				}
			}
		}

		ratios[t] = [2]int{pokemonsItKills, pokemonsThatKillIt}
	}

	sortedTypes := GetSortedRatioTypes(ratios, sortDir)
	printRatios(ratios, sortedTypes)
}

// For each combination of type, check how many pokemons in list,
// have a type that is at least super-effective against it.
func BestTypeComb(list []Pokemon, sortDir int) {
	histo := make(map[TypeCombination]int)

	for _, c := range TypeCombinations {
		consumeToHisto(c, list, histo)
	}

	sortedCombis := GetSortedIntCombis(histo, sortDir)
	printCombiHisto(histo, sortedCombis)
}

func consumeToHisto(combi TypeCombination, list []Pokemon, histo map[TypeCombination]int) {
	for _, p := range list {
		for _, t := range p.Types {
			s := combi.toSlice()
			if TypeEffectiveness(t, s) >= 2.0 {
				histo[combi] += 1
			}
		}
	}
}

type BestMoveSetResult struct {
	PokemonId   int     `json:"id"`
	PokemonName string  `json:"name"`
	MoveSet     [4]Move `json:"move_set"`
	TotalKt     int     `json:"kt"`
}

func BestPokemons(list []Pokemon, sortDir int) []BestMoveSetResult {
	fmt.Println("Analyzing optimal move sets...")
	bar := pb.StartNew(len(list))

	c := make(chan BestMoveSetResult, len(list))
	for _, p := range list {
		go func(p Pokemon) {
			moveSet, totalKt := BestMoveSet(p, list)
			bar.Increment()
			c <- BestMoveSetResult{p.Id, p.Name, moveSet, totalKt}
		}(p)
	}

	names := make(map[int]string)
	moveSets := make(map[int][4]Move)
	totalKts := make(map[int]int)
	for _, _ = range list {
		r := <-c
		names[r.PokemonId] = r.PokemonName
		moveSets[r.PokemonId] = r.MoveSet
		totalKts[r.PokemonId] = r.TotalKt
	}
	bar.FinishPrint("Finished finding best move sets.")

	sortedPokemonsIds := GetSortedPokemonIds(totalKts, sortDir)
	res := []BestMoveSetResult{}
	for _, pId := range sortedPokemonsIds {
		res = append(res, BestMoveSetResult{pId, names[pId], moveSets[pId], totalKts[pId]})
		PrintBattlePokemon(names[pId], moveSets[pId])
	}
	return res
}
