package poke

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const CACHED_FILE_LOCATION = "/tmp/pokemon.data"
const POKEMON_API = "https://pokeapi.co/api/v2/pokemon/"
const MOVE_API = "https://pokeapi.co/api/v2/move/"
const NUM_POKEMONS = 150
const NUM_MOVES = 639

func downloadData() PokemonData {
	data := PokemonData{}
	// Keep making API calls and appending to list.
	for i := 1; i <= NUM_POKEMONS; i++ {
		res, err := http.Get(POKEMON_API + strconv.Itoa(i))
		Check(err)
		bytes, err := ioutil.ReadAll(res.Body)
		Check(err)
		res.Body.Close()

		apiRes := PokemonApiResponse{}
		err = json.Unmarshal(bytes, &apiRes)
		Check(err)

		data.Responses = append(data.Responses, apiRes)
	}

	for i := 1; i <= NUM_MOVES; i++ {
		res, err := http.Get(MOVE_API + strconv.Itoa(i))
		Check(err)
		bytes, err := ioutil.ReadAll(res.Body)
		Check(err)
		res.Body.Close()

		apiRes := MoveApiResponse{}
		err = json.Unmarshal(bytes, &apiRes)
		Check(err)

		data.Moves = append(data.Moves, apiRes)
	}
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
