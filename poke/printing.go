package poke

import (
	"fmt"
	"github.com/raphamorim/go-rainbow"
	"strconv"
)

const LONGEST_POKEMON_NAME = 14
const LONGEST_TYPE_NAME = 8
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
	n := LONGEST_TYPE_NAME - len(t.Name)
	txt := strRepeat(n, " ") + t.Name + ":"
	return rainbow.Bold(rainbow.Hex(WHITE, txt))
}

func combiLabel(c TypeCombination) string {
	n := LONGEST_TYPE_NAME*2 + 1 - len(c.FirstSlot.Name) - len(c.SecondSlot.Name)
	var txt string
	if c.SecondSlot.Name != "" {
		txt = strRepeat(n, " ") + c.FirstSlot.Name + "-" + c.SecondSlot.Name + ":"
	} else {
		txt = strRepeat(n+1, " ") + c.FirstSlot.Name + ":"
	}
	return rainbow.Bold(rainbow.Hex(WHITE, txt))
}

func PrintHisto(histo map[Type]int, sorted [18]Type) {
	for _, t := range sorted {
		bar := strRepeat(histo[t], "#")
		fmt.Printf("%s %s (%d)\n",
			typeLabel(t),
			rainbow.Hex(t.HexColor, bar),
			histo[t])
	}
}

func PrintRatios(ratios map[Type][2]int, sorted [18]Type) {
	for _, t := range sorted {
		fmt.Printf("%s %s / %s (%f)\n",
			typeLabel(t),
			rainbow.Hex(t.HexColor, strconv.Itoa(ratios[t][0])),
			rainbow.Hex(t.HexColor, strconv.Itoa(ratios[t][1])),
			float64(ratios[t][0])/float64(ratios[t][1]))
	}
}

func PrintCombiHisto(histo map[TypeCombination]int, sorted []TypeCombination) {
	for _, c := range sorted {
		bar := strRepeat(histo[c], "#")
		fmt.Printf("%s %s (%d)\n",
			combiLabel(c),
			rainbow.Hex(FAV_BLUE, bar),
			histo[c])
	}
}

func PrintBattlePokemon(r BestMoveSetResult) {
	paddedName := strRepeat(LONGEST_POKEMON_NAME - len(r.PokemonName), " ") + r.PokemonName
	s := rainbow.Hex("#ffffff", paddedName+": [")
	s += MoveSetToString(r.MoveSet)
	fmt.Println(s + rainbow.Hex("#ffffff", "] ") + strconv.Itoa(r.TotalKt))
}

func MoveSetToString(moveSet [4]Move) string {
	s := ""
	for i, m := range moveSet {
		if i != 0 {
			s += ", "
		}
		var sOrP string
		if m.isPhysical {
			sOrP = "P"
		} else {
			sOrP = "S"
		}
		s += rainbow.Hex(m.Type.HexColor, fmt.Sprintf("%s(%s%d|%d)", m.Name, sOrP, m.Power, m.Accuracy))
	}
	return s
}

func PrintBestPokemonResults(results []BestMoveSetResult) {
	for _, r := range results {
		PrintBattlePokemon(r)
	}
}
