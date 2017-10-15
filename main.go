package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bcollazo/pokenalysis/poke"
)

var command string

func main() {
	flag.StringVar(&command, "command", "histo", "Command to compute.")
	flag.Parse()

	isValid := map[string]bool{
		"clean": true,
		"histo": true,
	}
	if !isValid[command] {
		panic("Bad Command")
	}

	list := poke.MaybeDownloadData()
	if command == "clean" {
		os.Remove(poke.CACHED_FILE_LOCATION)
	} else if command == "histo" {
		fmt.Printf("Running histo\n")
		poke.Histo(list)
	}
}
