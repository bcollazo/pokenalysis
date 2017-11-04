package main

import (
	"flag"
	"os"
	"strings"

	"github.com/bcollazo/pokenalysis/poke"
)

var command string
var gens string
var sort int

var GEN_BOUNDS = map[string][]int{
	"1": []int{1, 151},
	"2": []int{152, 251},
	"3": []int{252, 386},
	"4": []int{387, 494},
	"5": []int{495, 649},
	"6": []int{650, 721},
	"7": []int{722, 802},
}

func idsFromGens(gens string) []int {
	genKeys := strings.Split(gens, ",")
	ids := []int{}
	for _, k := range genKeys {
		genIds := poke.IntRange(GEN_BOUNDS[k][0], GEN_BOUNDS[k][1])
		ids = append(ids, genIds...)
	}
	return ids
}

func main() {
	flag.StringVar(&command, "command", "histo", "one of either 'histo', 'superhisto', 'goodratio', 'bestpoke'")
	flag.StringVar(&gens, "gens", "1,2,3,4,5,6,7", "comma-separated generations to include")
	flag.IntVar(&sort, "sort", 0, "sort direction. -1, 0, or 1")
	flag.Parse()

	isValid := map[string]bool{
		"clean":      true,
		"histo":      true,
		"superhisto": true,
		"goodratio":  true,
		"typecomb":   true,
		"bestpoke":   true,
	}
	if !isValid[command] {
		panic("Bad Command")
	}

	if command == "clean" {
		os.RemoveAll(poke.DATA_DIR)
		return
	}

	ids := idsFromGens(gens)
	poke.MaybeDownloadData(ids)
	list := poke.ReadDataFromLocal(ids)
	if command == "histo" {
		poke.Histo(list, sort)
	} else if command == "superhisto" {
		poke.SuperEffectiveHisto(list, sort)
	} else if command == "goodratio" {
		poke.GoodRatios(list, sort)
	} else if command == "typecomb" {
		poke.BestTypeComb(list, sort)
	} else if command == "bestpoke" {
		poke.BestPokemons(list, sort)
	}
}
