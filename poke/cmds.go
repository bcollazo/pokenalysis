package poke

import (
	"fmt"
	"strconv"

	"github.com/raphamorim/go-rainbow"
)

func buildBar(amount int) string {
	var toReturn string
	for i := 0; i < amount; i++ {
		toReturn += "#"
	}
	return toReturn
}

func typeLabel(t Type) string {
	return rainbow.Bold(rainbow.Hex("#FFFFFF", t.Name+":"))
}

func printHisto(histo map[Type]int) {
	for _, t := range TypeArr {
		bar := buildBar(histo[t])
		fmt.Printf("%s %s (%d)\n",
			typeLabel(t),
			rainbow.Hex(t.HexColor, bar),
			histo[t])
	}
}

func printRatios(ratios map[Type][2]int) {
	for _, t := range TypeArr {
		fmt.Printf("%s %s / %s (%f)\n",
			typeLabel(t),
			rainbow.Hex(t.HexColor, strconv.Itoa(ratios[t][0])),
			rainbow.Hex(t.HexColor, strconv.Itoa(ratios[t][1])),
			float64(ratios[t][0])/float64(ratios[t][1]))
	}
}

func Histo(list []Pokemon) {
	histo := make(map[Type]int)
	for _, p := range list {
		for _, t := range p.Types {
			histo[t] += 1
		}
	}

	printHisto(histo)
}

func SuperEffectiveHisto(list []Pokemon) {
	histo := make(map[Type]int)

	for _, pokemon := range list {
		for _, t := range TypeArr {
			// Check if super-effective.  If so, add
			if EffectiveMultiplier(t, pokemon.Types) >= 2.0 {
				histo[t] += 1
			}
		}
	}

	printHisto(histo)
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
func GoodRatios(list []Pokemon) {
	ratios := make(map[Type][2]int)
	for _, t := range TypeArr {
		pokemonsItKills := 0
		pokemonsThatKillIt := 0
		for _, pokemon := range list {
			// We are good against this pokemon
			if EffectiveMultiplier(t, pokemon.Types) >= 2.0 {
				pokemonsItKills += 1
			}

			// At least one of its type is good against us...
			for _, tt := range pokemon.Types {
				if EffectiveMultiplier(tt, []Type{t}) >= 2.0 {
					pokemonsThatKillIt += 1
					break
				}
			}
		}

		ratios[t] = [2]int{pokemonsItKills, pokemonsThatKillIt}
	}

	printRatios(ratios)
}
