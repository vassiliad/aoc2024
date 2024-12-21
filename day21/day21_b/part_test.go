package main

import (
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `
029A
980A
179A
456A
379A
`

	// println(small)

	puzzle, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", puzzle)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, 2, logger)
	const correctAnswer = 126384

	if sol != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", sol)
	}
}

func TestTiny(t *testing.T) {
	small := `
179A
`
	puzzle, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", puzzle)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, 2, logger)

	const correctAnswer = 68 * 179

	if sol != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", sol)
	}
}
