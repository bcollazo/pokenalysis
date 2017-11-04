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
	"sync"
)

const POKEMON_API = "https://pokeapi.co/api/v2/pokemon/"
const MOVE_API = "https://pokeapi.co/api/v2/move/"
const NUM_POKEMONS = 802
const NUM_MOVES = 639
const NUM_WORKERS = 4

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

func divideWork(r []int, n int) [][]int {
	step := len(r) / n
	res := [][]int{}
	for i := 0; i < n-1; i++ {
		a := i * step
		b := (i + 1) * step
		res = append(res, r[a:b])
	}
	res = append(res, r[(n-1)*step:])
	return res
}

func MaybeDownloadData(ids []int) {
	fmt.Printf("Downloading data to %s\n", DATA_DIR)
	// Ensure data directories exist.
	_ = os.MkdirAll(DATA_DIR, 0700)
	_ = os.MkdirAll(POKEMON_DATA_DIR, 0700)
	_ = os.MkdirAll(MOVES_DATA_DIR, 0700)

	// Download any missing pokemon.
	var wg sync.WaitGroup
	wg.Add(NUM_WORKERS)
	pokeBar := pb.StartNew(len(ids))
	subRanges := divideWork(ids, NUM_WORKERS)
	for _, r := range subRanges { // For range.
		go func(r []int) {
			for _, i := range r {
				maybeDownloadResource(POKEMON_API, i, pokemonPath(i))
				pokeBar.Increment()
			}
			wg.Done()
		}(r)
	}
	wg.Wait()
	pokeBar.FinishPrint("Finished downloading pokemons.")

	// Download any missing moves.
	wg.Add(NUM_WORKERS)
	movesBar := pb.StartNew(NUM_MOVES)
	subRanges = divideWork(IntRange(1, NUM_MOVES), NUM_WORKERS)
	for _, r := range subRanges {
		go func(r []int) {
			for _, i := range r {
				maybeDownloadResource(MOVE_API, i, movePath(i))
				movesBar.Increment()
			}
		}(r)
	}
	wg.Wait()
	movesBar.FinishPrint("Finished downloading moves.")
}

func ReadDataFromLocal(ids []int) []Pokemon {
	// Read from files into memory.
	data := PokemonData{}
	for _, i := range ids {
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
	movesMap := make(map[string]Move)
	for _, v := range data.Moves {
		movesMap[v.Name] = v.ToMove()
	}

	res := make([]Pokemon, len(data.Responses))
	for i, v := range data.Responses {
		res[i] = v.ToPokemon(movesMap)
	}
	return res
}
