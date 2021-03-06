package poke

import (
	"math"
)

const LIFE_MULTIPLIER = 100000
const LEVEL = 50.0
const DIED_PENALTY = 100

type BestMoveSetResult struct {
	PokemonId   int     `json:"id"`
	PokemonName string  `json:"name"`
	MoveSet     [4]Move `json:"move_set"`
	TotalKt     int     `json:"kt"`
}

func TypeEffectiveness(attType Type, defTypes []Type) float64 {
	mult := 1.0
	for _, t := range defTypes {
		mult *= EffectMap[attType.Name][t.Name]
	}
	return mult
}

// Analyze all pokemon. Returns its kt.
func BestMoveSet(pokemon Pokemon, list []Pokemon) (best [4]Move, bestKt int) {
	bestKt = 10000000000

	candidates := learnableCandidates(pokemon.LearnableMoves)
	combinations := GenerateCombinations(len(candidates), 4) // TODO: What if learns < 4
	for _, moveVector := range combinations {
		moveSet := intersectMoves(candidates, moveVector)

		// For each 150 pokemon:  Fight and get 'effective damage' of best move
		totalKt := 0
		for _, enemy := range list {
			kt := spamMoveFight(pokemon, moveSet, enemy)
			totalKt += kt
		}

		// Maintain bestKt
		if totalKt < bestKt {
			best = moveSet
			bestKt = totalKt
		}
	}
	return best, bestKt
}

// Skip moves where there exist a strictly better option
func learnableCandidates(moves []Move) []Move {
	res := []Move{}
	for _, m := range moves {
		thereExistBetter := false
		for _, mm := range moves {
			if m == mm {
				continue
			}
			sameCategory := (m.isPhysical && mm.isPhysical) || (!m.isPhysical && !mm.isPhysical)
			sameType := m.Type.Name == mm.Type.Name
			otherIsBetter := expectedPower(m) < expectedPower(mm)
			if sameCategory && sameType && otherIsBetter {
				thereExistBetter = true
				break
			}
		}
		if !thereExistBetter {
			res = append(res, m)
		}
	}
	return res
}

func intersectMoves(moveSpace []Move, moveVector []bool) [4]Move {
	moves := [4]Move{}
	i := 0
	for j, b := range moveVector {
		if b {
			moves[i] = moveSpace[j]
			i++
		}
	}
	return moves
}

// Simulates a fight where pokemons spam the best available move in every turn.
// Returns kt (kill turn).  If enemy wins, kt = 10
func spamMoveFight(a Pokemon, moveSet [4]Move, b Pokemon) (kt int) {
	// Check expected damage from each move.
	_, aDamage := bestMove(a, moveSet[:], b)
	_, bDamage := bestMove(b, b.LearnableMoves, a) // out of all possible moves.

	aKillTurn := int(math.Ceil(float64(b.BaseStats.Hp*LIFE_MULTIPLIER)/aDamage) * 2) // * 2 since have to wait for opponent
	bKillTurn := int(math.Ceil(float64(a.BaseStats.Hp*LIFE_MULTIPLIER)/bDamage) * 2)
	if a.BaseStats.Speed >= b.BaseStats.Speed {
		aKillTurn -= 1
	} else {
		bKillTurn -= 1
	}

	if aKillTurn > bKillTurn {
		aKillTurn += DIED_PENALTY
	}
	return aKillTurn
}

func bestMove(a Pokemon, moves []Move, b Pokemon) (move Move, damage float64) {
	damage = -1.0
	for _, m := range moves {
		d := computeMoveDamage(a, m, b)

		if d > damage {
			move = m
			damage = d
		}
	}
	return move, damage
}

func computeMoveDamage(a Pokemon, move Move, b Pokemon) float64 {
	var A float64
	var D float64
	if move.isPhysical {
		A = float64(a.BaseStats.Attack)
		D = float64(b.BaseStats.Defense)
	} else {
		A = float64(a.BaseStats.SpecialAttack)
		D = float64(b.BaseStats.SpecialDefense)
	}

	stab := 1.0 // same-type attack bonus
	for _, t := range a.Types {
		if move.Type == t {
			stab = 1.5
		}
	}
	levelTerm := ((2.0 * LEVEL) / 5.0) + 2.0
	modifier := TypeEffectiveness(move.Type, b.Types) * stab
	damage := (((levelTerm * expectedPower(move) * A / D) / 50) + 2) * modifier

	if damage <= 0.0 {
		damage = 1.0
	}
	return damage
}

func expectedPower(move Move) float64 {
	return float64(move.Power) * float64(move.Accuracy) / 100.0
}

// PROBLEM: When you fight; do you select best attack at the moment, or
// start with a plan and stick to it?  Maybe a combination: start with
// a plan and adjust appropriately (re-compute)?
