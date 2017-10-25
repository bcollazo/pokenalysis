package poke

import (
	"encoding/json"
	"fmt"
	"gopkg.in/cheggaaa/pb.v1"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const POKEMON_API = "https://pokeapi.co/api/v2/pokemon/"
const MOVE_API = "https://pokeapi.co/api/v2/move/"
const NUM_POKEMONS = 802
const NUM_MOVES = 639

var DATA_DIR = filepath.Join(os.TempDir(), "pokemon_data")
var POKEMON_DATA_DIR = filepath.Join(DATA_DIR, "pokemons")
var MOVES_DATA_DIR = filepath.Join(DATA_DIR, "moves")

// Assumes directories in path exist.
func maybeDownloadResource(baseUrl string, id int, path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		res, err := http.Get(baseUrl + strconv.Itoa(id))
		Check(err)

		bytes, err := ioutil.ReadAll(res.Body)
		Check(err)
		defer res.Body.Close()

		err = ioutil.WriteFile(path, bytes, 0644)
		Check(err)
	}
	return bytes
}

func pokemonPath(i int) string {
	return filepath.Join(POKEMON_DATA_DIR, strconv.Itoa(i)+".json")
}

func movePath(i int) string {
	return filepath.Join(MOVES_DATA_DIR, strconv.Itoa(i)+".json")
}

func MaybeDownloadData() []Pokemon {
	fmt.Printf("Downloading data to %s\n", DATA_DIR)
	// Ensure data directories exist.
	_ = os.MkdirAll(DATA_DIR, 0700)
	_ = os.MkdirAll(POKEMON_DATA_DIR, 0700)
	_ = os.MkdirAll(MOVES_DATA_DIR, 0700)

	// Download any missing pokemon.
	pokeBar := pb.StartNew(NUM_POKEMONS)
	for i := 1; i <= NUM_POKEMONS; i++ { // For range.
		maybeDownloadResource(POKEMON_API, i, pokemonPath(i))
		pokeBar.Increment()
	}
	pokeBar.FinishPrint("Finished downloading pokemons.")

	// Download any missing moves.
	movesBar := pb.StartNew(NUM_MOVES)
	for i := 1; i <= NUM_MOVES; i++ {
		maybeDownloadResource(MOVE_API, i, movePath(i))
		movesBar.Increment()
	}
	movesBar.FinishPrint("Finished downloading moves.")
}

func ReadDataFromLocal() []Pokemon {
	// Read from files into memory.
	data := PokemonData{}
	for i := 1; i <= NUM_POKEMONS; i++ {
		bytes, err := ioutil.ReadFile(pokemonPath(i))
		apiRes := PokemonApiResponse{}
		err = json.Unmarshal(bytes, &apiRes)
		Check(err)
		data.Responses = append(data.Responses, apiRes)
	}
	for i := 1; i <= NUM_MOVES; i++ {
		bytes, err := ioutil.ReadFile(movePath(i))
		apiRes := MoveApiResponse{}
		err = json.Unmarshal(bytes, &apiRes)
		Check(err)
		data.Moves = append(data.Moves, apiRes)
	}
	return toPokemonsArray(data)
}

func toPokemonsArray(data PokemonData) []Pokemon {
	movesMap := make(map[string]Attack)
	for _, v := range data.Moves {
		movesMap[v.Name] = v.ToAttack()
	}

	res := make([]Pokemon, len(data.Responses))
	for i, v := range data.Responses {
		res[i] = v.ToPokemon(movesMap)
	}
	return res
}
