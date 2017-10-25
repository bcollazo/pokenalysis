package poke

import (
	"encoding/json"
	"fmt"
	"gopkg.in/cheggaaa/pb.v1"
	"io/ioutil"
	"net/http"
	"strconv"
)

const CACHED_FILE_LOCATION = "/tmp/pokemon.data"
const POKEMON_API = "https://pokeapi.co/api/v2/pokemon/"
const MOVE_API = "https://pokeapi.co/api/v2/move/"
const NUM_POKEMONS = 150
const NUM_MOVES = 639

func downloadResource(baseUrl string, id int) []byte {
	res, err := http.Get(baseUrl + strconv.Itoa(id))
	Check(err)
	bytes, err := ioutil.ReadAll(res.Body)
	Check(err)
	defer res.Body.Close()
	return bytes
}

func downloadData() PokemonData {
	pokeBar := pb.StartNew(NUM_POKEMONS)
	data := PokemonData{}
	// Keep making API calls and appending to list.
	for i := 1; i <= NUM_POKEMONS; i++ {
		bytes := downloadResource(POKEMON_API, i)

		apiRes := PokemonApiResponse{}
		err := json.Unmarshal(bytes, &apiRes)
		Check(err)
		data.Responses = append(data.Responses, apiRes)

		pokeBar.Increment()
	}
	pokeBar.FinishPrint("Finished downloading pokemons.")

	movesBar := pb.StartNew(NUM_MOVES)
	for i := 1; i <= NUM_MOVES; i++ {
		bytes := downloadResource(MOVE_API, i)

		apiRes := MoveApiResponse{}
		err := json.Unmarshal(bytes, &apiRes)
		Check(err)
		data.Moves = append(data.Moves, apiRes)

		movesBar.Increment()
	}
	movesBar.FinishPrint("Finished downloading moves.")
	return data
}

func MaybeDownloadData() []Pokemon {
	// Check cached location.
	var data PokemonData
	bytes, err := ioutil.ReadFile(CACHED_FILE_LOCATION)
	if err != nil {
		fmt.Println("Downloading data... (this may take a while)")
		data = downloadData()

		bytes, err = json.Marshal(data)
		Check(err)
		err = ioutil.WriteFile(CACHED_FILE_LOCATION, bytes, 0644)
		Check(err)
	} else {
		data = PokemonData{}
		err = json.Unmarshal(bytes, &data)
		Check(err)
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
