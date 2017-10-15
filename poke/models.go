package poke

type Type struct {
	Name     string
	HexColor string
}

// Goes from PokeApi string to Type.
var Types = map[string]Type{
	"normal":   Type{"Normal", "#A8A77A"},
	"fighting": Type{"Fighting", "#C22E28"},
	"flying":   Type{"Flying", "#A98FF3"},
	"poison":   Type{"Poison", "#A33EA1"},
	"ground":   Type{"Ground", "#E2BF65"},
	"rock":     Type{"Ground", "#B6A136"},
	"bug":      Type{"Bug", "#A6B91A"},
	"ghost":    Type{"Ghost", "#735797"},
	"steel":    Type{"Steel", "#B7B7CE"},
	"fire":     Type{"Fire", "#EE8130"},
	"water":    Type{"Water", "#6390F0"},
	"grass":    Type{"Grass", "#7AC74C"},
	"electric": Type{"Electric", "#F7D02C"},
	"psychic":  Type{"Psychic", "#F95587"},
	"ice":      Type{"Ice", "#96D9D6"},
	"dragon":   Type{"Dragon", "#6F35FC"},
	"dark":     Type{"Dark", "#705746"},
	"fairy":    Type{"Fairy", "#D685AD"},
}

type Pokemon struct {
	Name   string
	Weight int
	Type   Type
}

type PokemonList struct {
	Pokemons []Pokemon
}

type PokemonApiResponse struct {
	Name   string `json:name`
	Weight int    `json:weight`
	Types  []struct {
		Slot int `json:slot`
		Type struct {
			Url  string `json:url`
			Name string `json:name`
		} `json:type`
	} `json:types`
}

func (r PokemonApiResponse) ToPokemon() Pokemon {
	var t Type
	if len(r.Types) == 1 {
		t = Types[r.Types[0].Type.Name]
	}
	return Pokemon{r.Name, r.Weight, t}
}
