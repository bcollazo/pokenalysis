package poke

import "testing"

func TestBestAttack(t *testing.T) {
	aAttacks := []Attack{
		Attack{"Growl", 0, 100, TypeArr[0]},
		Attack{"Vine Whip", 40, 100, TypeArr[4]},
	}
	a := Pokemon{"Bulbasaur", 1, []Type{TypeArr[4], TypeArr[7]}, aAttacks}
	bAttacks := []Attack{
		Attack{"Water Gun", 30, 100, TypeArr[2]},
		Attack{"Tail Whip", 0, 100, TypeArr[0]},
	}
	b := Pokemon{"Squirtle", 1, []Type{TypeArr[2]}, bAttacks}

	attack := BestAttack(a, b)
	if attack != aAttacks[1] {
		t.Error("Should spam Vine Whip.")
	}

	attack = BestAttack(b, a)
	if attack != bAttacks[0] {
		t.Error("Should spam Water Gun.")
	}
}
