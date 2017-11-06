package main

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/bcollazo/pokenalysis/poke"
	"github.com/bcollazo/pokenalysis/serve"
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
	flag.StringVar(&command, "command", "histo", "one of either 'histo', 'superhisto', 'goodratio', 'bestpoke', 'work'")
	flag.StringVar(&gens, "gens", "1", "comma-separated generations to include")
	flag.IntVar(&sort, "sort", 0, "sort direction. -1, 0, or 1")

	flag.StringVar(&host, "host", "localhost", "host where this code is running.  Used when command is 'master'")
	flag.IntVar(&port, "port", 3000, "port to use if command is 'master' or 'serve'")
	flag.StringVar(&machines, "machines", "localhost:3000", "comma-separated hostnames")
	flag.Parse()

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

	if command == "clean" {
		os.RemoveAll(poke.DATA_DIR)
		return
	}

	ids := idsFromGens(gens)
	var stringPort = ":" + strconv.Itoa(port)
	if command == "master" {
		parsed := strings.Split(machines, ",")
		serve.StartMaster(ids, host, stringPort, parsed)
		return
	}

	// Commands that need data ready.
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
	} else if command == "serve" {
		serve.StartWorker(stringPort)
	}
}
