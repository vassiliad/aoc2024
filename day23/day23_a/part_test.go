package main

import (
	"log/slog"
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `
kh-tc
qp-kh
de-cg
ka-co
yn-aq
qp-ub
cg-tb
vc-aq
tb-ka
wh-tc
yn-cg
kh-ub
ta-co
de-co
tc-td
tb-wq
wh-td
ta-ka
td-qp
aq-cg
wq-ub
ub-vc
de-ta
wq-aq
wq-vc
wh-yn
ka-de
kh-ta
co-tc
wh-qp
tb-vc
td-yn
`

	// println(small)

	puzzle, err := util.ReadString(small)

	slog.Info("Input is", "input", puzzle)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle)
	const correctAnswer = 7

	if sol != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", sol)
	}
}

func TestTiny(t *testing.T) {
	small := `
ta-co
de-co
de-ta
`

	// println(small)

	puzzle, err := util.ReadString(small)

	slog.Info("Input is", "input", puzzle)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle)
	const correctAnswer = 1

	if sol != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", sol)
	}
}
