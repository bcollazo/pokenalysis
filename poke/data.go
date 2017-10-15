package poke

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

const CACHED_FILE_LOCATION = "/tmp/pokemon.data"
const BASE_API_URL = "https://pokeapi.co/api/v2/pokemon/"
const NUM_POKEMONS = 150

func downloadData() PokemonData {
	data := PokemonData{}
	// Keep making API calls and appending to list.
	for i := 1; i <= NUM_POKEMONS; i++ {
		res, err := http.Get(BASE_API_URL + strconv.Itoa(i))
		Check(err)
		bytes, err := ioutil.ReadAll(res.Body)
		Check(err)
		res.Body.Close()

		apiRes := PokemonApiResponse{}
		err = json.Unmarshal(bytes, &apiRes)
		Check(err)

		data.Responses = append(data.Responses, apiRes)
	}
	return data
}

func MaybeDownloadData() []Pokemon {
	// Check cached location.
	var data PokemonData
	bytes, err := ioutil.ReadFile(CACHED_FILE_LOCATION)
	if err != nil {
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
	res := make([]Pokemon, len(data.Responses))
	for i, v := range data.Responses {
		res[i] = v.ToPokemon()
	}
	return res
}
