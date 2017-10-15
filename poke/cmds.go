package poke

import (
	"fmt"

	"github.com/raphamorim/go-rainbow"
)

func buildBar(amount int) string {
	var toReturn string
	for i := 0; i < amount; i++ {
		toReturn += "#"
	}
	return toReturn
}

func Histo(list PokemonList) {
	histo := make(map[Type]int)
	for _, p := range list.Pokemons {
		histo[p.Type] += 1
	}
	for _, t := range Types {
		bar := buildBar(histo[t])
		fmt.Printf("%s %s\n",
			rainbow.Bold(rainbow.Hex("#FFFFFF", t.Name+":")),
			rainbow.Hex(t.HexColor, bar))
	}
}
