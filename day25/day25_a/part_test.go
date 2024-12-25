package main

import (
	"log/slog"
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `
#####
.####
.####
.####
.#.#.
.#...
.....

#####
##.##
.#.##
...##
...#.
...#.
.....

.....
#....
#....
#...#
#.#.#
#.###
#####

.....
.....
#.#..
###..
###.#
###.#
#####

.....
.....
.....
#....
#.#..
#.#.#
#####
`

	// println(small)

	puzzle, err := util.ReadString(small)

	slog.Info("Input is", "input", puzzle)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle)
	const correctAnswer = 3

	if sol != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", sol)
	}
}

func TestTiny(t *testing.T) {
	small := `
#####
.####
.####
.####
.#.#.
.#...
.....

.....
.....
.....
#....
#.#..
#.#.#
#####
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
