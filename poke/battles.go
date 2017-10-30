package poke

import (
	"math"
)

const LEVEL = 50.0

// Returns the best move to use, along with expected kill turn.
func BestMove(a BattlePokemon, b Pokemon) (bestMove Move, ekt float64) {
	var best Move
	bestEkt := 10000.0

	// Check expected damage from each move.
	for _, move := range a.Moves {
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

		numAttacksNeeded := math.Ceil(float64(b.BaseStats.Hp) / damage)
		var ekt float64
		if a.BaseStats.Speed > b.BaseStats.Speed {
			ekt = numAttacksNeeded*2 - 1
		} else {
			ekt = numAttacksNeeded * 2 // Since he always goes second.
		}

		if ekt < bestEkt {
			best = move
			bestEkt = ekt
		}
	}

	return best, bestEkt
}

func expectedPower(move Move) float64 {
	return float64(move.Power) * float64(move.Accuracy) / 100.0
}

// PROBLEM: When you fight; do you select best attack at the moment, or
// start with a plan and stick to it?  Maybe a combination: start with
// a plan and adjust appropriately (re-compute)?
