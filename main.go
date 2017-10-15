package main

import (
	"flag"
	"fmt"

	"github.com/bcollazo/pokenalysis/poke"
)

var command string

func main() {
	flag.StringVar(&command, "command", "histo", "Command to compute.")
	flag.Parse()

	isValid := map[string]bool{
		"histo": true,
	}
	if !isValid[command] {
		panic("Bad Command")
	}

	list := poke.MaybeDownloadData()
	if command == "histo" {
		fmt.Printf("Running histo\n")
		poke.Histo(list)
	}
}
