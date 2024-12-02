package main

import (
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9

`

	println(small)

	input, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	solution := solution(input, logger)
	const correct_answer = 2

	if solution != correct_answer {
		t.Fatal("Expected answer to be", correct_answer, "but it was", solution)
	}
}
