package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bcollazo/pokenalysis/poke"
	"github.com/bcollazo/pokenalysis/serve"
	"github.com/pkg/profile"
)

var command string
var gens string
var sort int
var host string
var port int
var machines string

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
	defer profile.Start().Stop()
	flag.StringVar(&command, "command", "histo", "one of either 'histo', 'superhisto', 'goodratio', 'bestpoke', 'work'")
	flag.StringVar(&gens, "gens", "1", "comma-separated generations to include")
	flag.IntVar(&sort, "sort", 0, "sort direction. -1, 0, or 1")

	flag.StringVar(&host, "host", "localhost", "host where this code is running.  Used when command is 'master'")
	flag.IntVar(&port, "port", 3000, "port to use if command is 'master' or 'serve'")
	flag.StringVar(&machines, "machines", "localhost:3000", "comma-separated hostnames")
	flag.Parse()

	// Validate flags.
	isValid := map[string]bool{
		"clean":      true,
		"histo":      true,
		"superhisto": true,
		"goodratio":  true,
		"typecomb":   true,
		"bestpoke":   true,
		"serve":      true,
		"master":     true,
	}
	if !isValid[command] {
		panic("Bad Command")
	}

	// ===== clean command
	if command == "clean" {
		os.RemoveAll(poke.DATA_DIR)
		return
	}

	// ===== mater command
	ids := idsFromGens(gens)
	var stringPort = ":" + strconv.Itoa(port)
	if command == "master" {
		parsed := strings.Split(machines, ",")
		serve.StartMaster(ids, host, stringPort, parsed)
		return
	}

	// ===== the rest of the commands (require data)
	poke.MaybeDownloadData(ids)
	list := poke.ReadDataFromLocal(ids)
	if command == "histo" {
		histo, sortedTypes := poke.Histo(list, sort)
		poke.PrintHisto(histo, sortedTypes)
	} else if command == "superhisto" {
		histo, sortedTypes := poke.SuperEffectiveHisto(list, sort)
		poke.PrintHisto(histo, sortedTypes)
	} else if command == "goodratio" {
		ratios, sortedTypes := poke.GoodRatios(list, sort)
		poke.PrintRatios(ratios, sortedTypes)
	} else if command == "typecomb" {
		histo, sortedCombis := poke.BestTypeComb(list, sort)
		poke.PrintCombiHisto(histo, sortedCombis)
	} else if command == "bestpoke" {
		fmt.Println("Analyzing optimal move sets...")
		results := poke.BestPokemons(list, sort)
		poke.PrintBestPokemonResults(results)
	} else if command == "serve" {
		serve.StartWorker(stringPort)
	}
}
