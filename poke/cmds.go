package poke

import (
	"fmt"
	"strconv"

	"github.com/raphamorim/go-rainbow"
)

const LONGEST_TYPE_NAME_LEN = 8

func strRepeat(amount int, str string) string {
	var toReturn string
	for i := 0; i < amount; i++ {
		toReturn += str
	}
	return toReturn
}

func typeLabel(t Type) string {
	padding := strRepeat(LONGEST_TYPE_NAME_LEN-len(t.Name), " ")
	return rainbow.Bold(rainbow.Hex("#FFFFFF", padding+t.Name+":"))
}

func printHisto(histo map[Type]int) {
	for _, t := range TypeArr {
		bar := strRepeat(histo[t], "#")
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

// Number of pokemons such type is good against.
func SuperEffectiveHisto(list []Pokemon) {
	histo := make(map[Type]int)

	for _, pokemon := range list {
		for _, t := range TypeArr {
			// Check if super-effective.  If so, add
			if TypeEffectiveness(t, pokemon.Types) >= 2.0 {
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

	printRatios(ratios)
}

func BestPokemons(list []Pokemon) {
	for _, p := range list {
		version, g := bestVersion(p, list)
		fmt.Println(version.ToString())
	}
}

// Analyze all pokemon. Returns its ekt.
func bestVersion(pokemon Pokemon, list []Pokemon) (best BattlePokemon, goodness float64) {
	bestGoodness := -1.0

	combinations := GenerateCombinations(len(pokemon.LearnableMoves), 4)
	for _, moveVector := range combinations {
		// Create battle pokemon
		moves := IntersectMoves(pokemon, moveVector)
		battlePoke := BattlePokemon{pokemon, moves}

		// For each 150 pokemon:  Fight and get 'effective damage' of best move
		damagePerPokemon := make(map[int]float64)
		for i, enemy := range list {
			_, expectedDamage := BestMove(battlePoke, enemy)
			damagePerPokemon[i] = expectedDamage
		}
		// Every damage output is weighted by the health of opposing pokemon.
		// Idea is that its important to deal high damage with high hp opponents.
		// a.k.a. it hurts more to be weak against high hp opponents.
		goodness := weightByHealth(damagePerPokemon, list)

		// Maintain bestGoodness
		if goodness > bestGoodness {
			best = battlePoke
			bestGoodness = goodness
		}
	}
	return best, bestGoodness
}

func weightByHealth(damagePerPokemon map[int]float64, list []Pokemon) float64 {
	total := 0.0
	totalHealth := 0
	for i, p := range list {
		total += damagePerPokemon[i] * float64(p.BaseStats.Hp)
		totalHealth += p.BaseStats.Hp
	}
	return total / float64(totalHealth)
}
