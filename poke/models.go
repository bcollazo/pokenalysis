package poke

type Type struct {
	Name     string `json:"type"`
	HexColor string `json:"hex"`
}

var TypeArr = [18]Type{
	Type{"Normal", "#A8A77A"},
	Type{"Fire", "#EE8130"},
	Type{"Water", "#6390F0"},
	Type{"Electric", "#F7D02C"},
	Type{"Grass", "#7AC74C"},
	Type{"Ice", "#96D9D6"},
	Type{"Fighting", "#C22E28"},
	Type{"Poison", "#A33EA1"},
	Type{"Ground", "#E2BF65"},
	Type{"Flying", "#A98FF3"},
	Type{"Psychic", "#F95587"},
	Type{"Bug", "#A6B91A"},
	Type{"Rock", "#B6A136"},
	Type{"Ghost", "#735797"},
	Type{"Dragon", "#6F35FC"},
	Type{"Dark", "#705746"},
	Type{"Steel", "#B7B7CE"},
	Type{"Fairy", "#D685AD"},
}

type TypeCombination struct {
	FirstSlot  Type
	SecondSlot Type
}

func (combi TypeCombination) toSlice() []Type {
	s := []Type{}
	if combi.FirstSlot.Name != "" {
		s = append(s, combi.FirstSlot)
	}
	if combi.SecondSlot.Name != "" {
		s = append(s, combi.SecondSlot)
	}
	return s
}

var TypeCombinations []TypeCombination

type Pokemon struct {
	Id             int
	Name           string
	Weight         int
	Types          []Type
	LearnableMoves []Move
	BaseStats      Stats
}

type Move struct {
	Name       string `json:"name"`
	Power      int    `json:"power"`
	Accuracy   int    `json:"accuracy"`
	Type       Type   `json:"type"`
	isPhysical bool   `json:"is_physical"` // else isSpecial
}

type Stats struct {
	Hp             int
	Attack         int
	Defense        int
	SpecialAttack  int
	SpecialDefense int
	Speed          int
}

type PokemonData struct {
	Responses []PokemonApiResponse
	Moves     []MoveApiResponse
}

type PokemonApiResponse struct {
	Id     int    `json:id`
	Name   string `json:name`
	Weight int    `json:weight`
	Types  []struct {
		Slot int `json:slot`
		Type struct {
			Name string `json:name`
		} `json:type`
	} `json:types`
	Stats []struct {
		Stat struct {
			Name string `json:name`
		} `json:stat`
		Effort   int `json:effort`
		BaseStat int `json:"base_stat"`
	} `json:stats`
	Moves []struct {
		VersionGroupDetails []struct {
			MoveLearnMethod struct {
				Name string `json:name`
			} `json:"move_learn_method"`
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
			} `json:"version_group"`
		} `json:"version_group_details"`
		Move struct {
			Name string `json:name`
		} `json:move`
	} `json:moves`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:sprites`
	BaseExperience int `json:"base_experience"`
}

func (r PokemonApiResponse) ToPokemon(movesMap map[string]Move) Pokemon {
	var types []Type
	for _, t := range r.Types {
		types = append(types, ApiToType[t.Type.Name])
	}

	var learnable []Move
	for _, m := range r.Moves {
		for _, v := range m.VersionGroupDetails {
			if v.VersionGroup.Name == "x-y" &&
				v.MoveLearnMethod.Name == "level-up" {
				learnable = append(learnable, movesMap[m.Move.Name])
			}
		}
	}

	var baseStats Stats
	for _, s := range r.Stats {
		switch s.Stat.Name {
		case "hp":
			baseStats.Hp = s.BaseStat
		case "attack":
			baseStats.Attack = s.BaseStat
		case "defense":
			baseStats.Defense = s.BaseStat
		case "special-attack":
			baseStats.SpecialAttack = s.BaseStat
		case "special-defense":
			baseStats.SpecialDefense = s.BaseStat
		case "speed":
			baseStats.Speed = s.BaseStat
		}
	}

	return Pokemon{r.Id, r.Name, r.Weight, types, learnable, baseStats}
}

type MoveApiResponse struct {
	Id       int    `json:id`
	Name     string `json:name`
	Power    int    `json:power`
	Accuracy int    `json:accuracy`
	Type     struct {
		Name string `json:name`
	} `json:type`
	DamageClass struct {
		Name string `json:name`
	} `json:"damage_class"`
}

func (r MoveApiResponse) ToMove() Move {
	isPhysical := r.DamageClass.Name == "physical"
	return Move{r.Name, r.Power, r.Accuracy, ApiToType[r.Type.Name], isPhysical}
}

// Goes from PokeApi string to Type.
var ApiToType = map[string]Type{
	"normal":   Type{"Normal", "#A8A77A"},
	"fighting": Type{"Fighting", "#C22E28"},
	"flying":   Type{"Flying", "#A98FF3"},
	"poison":   Type{"Poison", "#A33EA1"},
	"ground":   Type{"Ground", "#E2BF65"},
	"rock":     Type{"Rock", "#B6A136"},
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

var EffectMap = map[string]map[string]float64{
	"Normal": map[string]float64{
		"Normal":   1.0,
		"Fire":     1.0,
		"Water":    1.0,
		"Electric": 1.0,
		"Grass":    1.0,
		"Ice":      1.0,
		"Fighting": 1.0,
		"Poison":   1.0,
		"Ground":   1.0,
		"Flying":   1.0,
		"Psychic":  1.0,
		"Bug":      1.0,
		"Rock":     0.5,
		"Ghost":    0.0,
		"Dragon":   1.0,
		"Dark":     1.0,
		"Steel":    0.5,
		"Fairy":    1.0,
	},
	"Fire": map[string]float64{
		"Normal":   1.0,
		"Fire":     0.5,
		"Water":    0.5,
		"Electric": 1.0,
		"Grass":    2.0,
		"Ice":      2.0,
		"Fighting": 1.0,
		"Poison":   1.0,
		"Ground":   1.0,
		"Flying":   1.0,
		"Psychic":  1.0,
		"Bug":      2.0,
		"Rock":     0.5,
		"Ghost":    1.0,
		"Dragon":   0.5,
		"Dark":     1.0,
		"Steel":    2.0,
		"Fairy":    1.0,
	},
	"Water": map[string]float64{
		"Normal":   1.0,
		"Fire":     2.0,
		"Water":    0.5,
		"Electric": 1.0,
		"Grass":    0.5,
		"Ice":      1.0,
		"Fighting": 1.0,
		"Poison":   1.0,
		"Ground":   2.0,
		"Flying":   1.0,
		"Psychic":  1.0,
		"Bug":      1.0,
		"Rock":     2.0,
		"Ghost":    1.0,
		"Dragon":   0.5,
		"Dark":     1.0,
		"Steel":    1.0,
		"Fairy":    1.0,
	},
	"Electric": map[string]float64{
		"Normal":   1.0,
		"Fire":     1.0,
		"Water":    2.0,
		"Electric": 0.5,
		"Grass":    0.5,
		"Ice":      1.0,
		"Fighting": 1.0,
		"Poison":   1.0,
		"Ground":   0.0,
		"Flying":   2.0,
		"Psychic":  1.0,
		"Bug":      1.0,
		"Rock":     1.0,
		"Ghost":    1.0,
		"Dragon":   0.5,
		"Dark":     1.0,
		"Steel":    1.0,
		"Fairy":    1.0,
	},
	"Grass": map[string]float64{
		"Normal":   1.0,
		"Fire":     0.5,
		"Water":    2.0,
		"Electric": 1.0,
		"Grass":    0.5,
		"Ice":      1.0,
		"Fighting": 1.0,
		"Poison":   0.5,
		"Ground":   2.0,
		"Flying":   0.5,
		"Psychic":  1.0,
		"Bug":      0.5,
		"Rock":     2.0,
		"Ghost":    1.0,
		"Dragon":   0.5,
		"Dark":     1.0,
		"Steel":    0.5,
		"Fairy":    1.0,
	},
	"Ice": map[string]float64{
		"Normal":   1.0,
		"Fire":     0.5,
		"Water":    0.5,
		"Electric": 1.0,
		"Grass":    2.0,
		"Ice":      0.5,
		"Fighting": 1.0,
		"Poison":   1.0,
		"Ground":   2.0,
		"Flying":   2.0,
		"Psychic":  1.0,
		"Bug":      1.0,
		"Rock":     1.0,
		"Ghost":    1.0,
		"Dragon":   2.0,
		"Dark":     1.0,
		"Steel":    0.5,
		"Fairy":    1.0,
	},
	"Fighting": map[string]float64{
		"Normal":   2.0,
		"Fire":     1.0,
		"Water":    1.0,
		"Electric": 1.0,
		"Grass":    1.0,
		"Ice":      2.0,
		"Fighting": 1.0,
		"Poison":   0.5,
		"Ground":   1.0,
		"Flying":   0.5,
		"Psychic":  0.5,
		"Bug":      0.5,
		"Rock":     2.0,
		"Ghost":    0.0,
		"Dragon":   1.0,
		"Dark":     2.0,
		"Steel":    2.0,
		"Fairy":    0.5,
	},
	"Poison": map[string]float64{
		"Normal":   1.0,
		"Fire":     1.0,
		"Water":    1.0,
		"Electric": 1.0,
		"Grass":    2.0,
		"Ice":      1.0,
		"Fighting": 1.0,
		"Poison":   0.5,
		"Ground":   0.5,
		"Flying":   1.0,
		"Psychic":  1.0,
		"Bug":      1.0,
		"Rock":     0.5,
		"Ghost":    0.5,
		"Dragon":   1.0,
		"Dark":     1.0,
		"Steel":    0.0,
		"Fairy":    2.0,
	},
	"Ground": map[string]float64{
		"Normal":   1.0,
		"Fire":     2.0,
		"Water":    1.0,
		"Electric": 2.0,
		"Grass":    0.5,
		"Ice":      1.0,
		"Fighting": 1.0,
		"Poison":   2.0,
		"Ground":   1.0,
		"Flying":   0.0,
		"Psychic":  1.0,
		"Bug":      0.5,
		"Rock":     2.0,
		"Ghost":    1.0,
		"Dragon":   1.0,
		"Dark":     1.0,
		"Steel":    2.0,
		"Fairy":    1.0,
	},
	"Flying": map[string]float64{
		"Normal":   1.0,
		"Fire":     1.0,
		"Water":    1.0,
		"Electric": 0.5,
		"Grass":    2.0,
		"Ice":      1.0,
		"Fighting": 2.0,
		"Poison":   1.0,
		"Ground":   1.0,
		"Flying":   1.0,
		"Psychic":  1.0,
		"Bug":      2.0,
		"Rock":     0.5,
		"Ghost":    1.0,
		"Dragon":   1.0,
		"Dark":     1.0,
		"Steel":    0.5,
		"Fairy":    1.0,
	},
	"Psychic": map[string]float64{
		"Normal":   1.0,
		"Fire":     1.0,
		"Water":    1.0,
		"Electric": 1.0,
		"Grass":    1.0,
		"Ice":      1.0,
		"Fighting": 2.0,
		"Poison":   2.0,
		"Ground":   1.0,
		"Flying":   1.0,
		"Psychic":  0.5,
		"Bug":      1.0,
		"Rock":     1.0,
		"Ghost":    1.0,
		"Dragon":   1.0,
		"Dark":     0.0,
		"Steel":    0.5,
		"Fairy":    1.0,
	},
	"Bug": map[string]float64{
		"Normal":   1.0,
		"Fire":     0.5,
		"Water":    1.0,
		"Electric": 1.0,
		"Grass":    2.0,
		"Ice":      1.0,
		"Fighting": 0.5,
		"Poison":   0.5,
		"Ground":   1.0,
		"Flying":   0.5,
		"Psychic":  2.0,
		"Bug":      1.0,
		"Rock":     1.0,
		"Ghost":    0.5,
		"Dragon":   1.0,
		"Dark":     2.0,
		"Steel":    0.5,
		"Fairy":    0.5,
	},
	"Rock": map[string]float64{
		"Normal":   1.0,
		"Fire":     2.0,
		"Water":    1.0,
		"Electric": 1.0,
		"Grass":    1.0,
		"Ice":      2.0,
		"Fighting": 0.5,
		"Poison":   1.0,
		"Ground":   0.5,
		"Flying":   2.0,
		"Psychic":  1.0,
		"Bug":      2.0,
		"Rock":     1.0,
		"Ghost":    1.0,
		"Dragon":   1.0,
		"Dark":     1.0,
		"Steel":    0.5,
		"Fairy":    1.0,
	},
	"Ghost": map[string]float64{
		"Normal":   0.0,
		"Fire":     1.0,
		"Water":    1.0,
		"Electric": 1.0,
		"Grass":    1.0,
		"Ice":      1.0,
		"Fighting": 1.0,
		"Poison":   1.0,
		"Ground":   1.0,
		"Flying":   1.0,
		"Psychic":  2.0,
		"Bug":      1.0,
		"Rock":     1.0,
		"Ghost":    2.0,
		"Dragon":   1.0,
		"Dark":     0.5,
		"Steel":    1.0,
		"Fairy":    1.0,
	},
	"Dragon": map[string]float64{
		"Normal":   1.0,
		"Fire":     1.0,
		"Water":    1.0,
		"Electric": 1.0,
		"Grass":    1.0,
		"Ice":      1.0,
		"Fighting": 1.0,
		"Poison":   1.0,
		"Ground":   1.0,
		"Flying":   1.0,
		"Psychic":  1.0,
		"Bug":      1.0,
		"Rock":     1.0,
		"Ghost":    1.0,
		"Dragon":   2.0,
		"Dark":     1.0,
		"Steel":    0.5,
		"Fairy":    0.0,
	},
	"Dark": map[string]float64{
		"Normal":   1.0,
		"Fire":     1.0,
		"Water":    1.0,
		"Electric": 1.0,
		"Grass":    1.0,
		"Ice":      1.0,
		"Fighting": 0.5,
		"Poison":   1.0,
		"Ground":   1.0,
		"Flying":   1.0,
		"Psychic":  2.0,
		"Bug":      1.0,
		"Rock":     1.0,
		"Ghost":    2.0,
		"Dragon":   1.0,
		"Dark":     0.5,
		"Steel":    1.0,
		"Fairy":    0.5,
	},
	"Steel": map[string]float64{
		"Normal":   1.0,
		"Fire":     0.5,
		"Water":    0.5,
		"Electric": 0.5,
		"Grass":    1.0,
		"Ice":      2.0,
		"Fighting": 1.0,
		"Poison":   1.0,
		"Ground":   1.0,
		"Flying":   1.0,
		"Psychic":  1.0,
		"Bug":      1.0,
		"Rock":     2.0,
		"Ghost":    1.0,
		"Dragon":   1.0,
		"Dark":     1.0,
		"Steel":    0.5,
		"Fairy":    2.0,
	},
	"Fairy": map[string]float64{
		"Normal":   1.0,
		"Fire":     0.5,
		"Water":    1.0,
		"Electric": 1.0,
		"Grass":    1.0,
		"Ice":      1.0,
		"Fighting": 2.0,
		"Poison":   0.5,
		"Ground":   1.0,
		"Flying":   1.0,
		"Psychic":  1.0,
		"Bug":      1.0,
		"Rock":     1.0,
		"Ghost":    1.0,
		"Dragon":   2.0,
		"Dark":     2.0,
		"Steel":    0.5,
		"Fairy":    1.0,
	},
}

func init() {
	var combi TypeCombination
	for i := 0; i < len(TypeArr); i++ {
		a := TypeArr[i]
		// Single-type combinations.
		combi.FirstSlot = a
		TypeCombinations = append(TypeCombinations, combi)

		for j := i + 1; j < len(TypeArr); j++ {
			b := TypeArr[j]
			// Two-type combinations.
			combi.SecondSlot = b
			TypeCombinations = append(TypeCombinations, combi)
		}
	}
}
