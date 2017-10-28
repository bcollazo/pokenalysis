package poke

// Returns the best move to use, along with expected damage.
func BestMove(a BattlePokemon, b Pokemon) (bestMove Move, expectedDamage float64) {
	var best Move
	bestDamage := -1.0

	// Check expected damage from each move.
	for _, move := range a.Moves {
		var attOrSpa float64
		var defOrSpd float64
		if move.isPhysical {
			attOrSpa = float64(a.BaseStats.Attack)
			defOrSpd = float64(b.BaseStats.Defense)
		} else {
			attOrSpa = float64(a.BaseStats.SpecialAttack)
			defOrSpd = float64(b.BaseStats.SpecialDefense)
		}

		stab := 1.0 // same-type attack bonus
		for _, t := range a.Types {
			if move.Type == t {
				stab = 1.5
			}
		}
		t := TypeEffectiveness(move.Type, b.Types)
		damage := t * stab * expectedPower(move) * attOrSpa / defOrSpd

		if damage > bestDamage {
			best = move
			bestDamage = damage
		}
	}

	return best, bestDamage
}

func expectedPower(move Move) float64 {
	return float64(move.Power) * float64(move.Accuracy)
}

// PROBLEM: When you fight; do you select best attack at the moment, or
// start with a plan and stick to it?  Maybe a combination: start with
// a plan and adjust appropriately (re-compute)?
