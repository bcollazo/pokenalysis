package main

import (
	"os"

	"github.com/bcollazo/pokenalysis/poke"
)

func main() {
	command := os.Args[1]

	isValid := map[string]bool{
		"clean":      true,
		"histo":      true,
		"superhisto": true,
		"goodratio":  true,
	}
	if !isValid[command] {
		panic("Bad Command")
	}

	if command == "clean" {
		os.RemoveAll(poke.DATA_DIR)
		return
	}

	poke.MaybeDownloadData()
	list := poke.ReadDataFromLocal()
	if command == "histo" {
		poke.Histo(list)
	} else if command == "superhisto" {
		poke.SuperEffectiveHisto(list)
	} else if command == "goodratio" {
		poke.GoodRatios(list)
	}
}
