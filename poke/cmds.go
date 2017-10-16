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

func printHisto(histo map[Type]int) {
	for _, t := range TypeArr {
		bar := buildBar(histo[t])
		fmt.Printf("%s %s (%d)\n",
			rainbow.Bold(rainbow.Hex("#FFFFFF", t.Name+":")),
			rainbow.Hex(t.HexColor, bar),
			histo[t])
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
