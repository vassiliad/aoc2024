package main

import (
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `3   4
4   3
2   5
1   3
3   9
3   3
`

	println(small)

	input, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	solution := solution(input, logger)
	const correct_answer = 11

	if solution != correct_answer {
		t.Fatal("Expected answer to be", correct_answer, "but it was", solution)
	}
}
