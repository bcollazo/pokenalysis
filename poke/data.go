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

func downloadData() PokemonList {
	list := PokemonList{}
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

		pokemon := apiRes.ToPokemon()
		list.Pokemons = append(list.Pokemons, pokemon)
	}
	return list
}

func MaybeDownloadData() PokemonList {
	// Check cached location.
	var list PokemonList
	bytes, err := ioutil.ReadFile(CACHED_FILE_LOCATION)
	if err != nil {
		list = downloadData()

		bytes, err = json.Marshal(list)
		Check(err)
		err = ioutil.WriteFile(CACHED_FILE_LOCATION, bytes, 0644)
		Check(err)
	} else {
		list = PokemonList{}
		err = json.Unmarshal(bytes, &list)
		Check(err)
	}
	return list
}
