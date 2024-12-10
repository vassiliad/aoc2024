package main

import (
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732
`

	println(small)

	input, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	solution := solution(input, logger)
	const correctAnswer = 36

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}
}

func TestTiny(t *testing.T) {
	small := `
0123
1234
8765
9876
`

	println(small)

	input, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	solution := solution(input, logger)
	const correctAnswer = 1

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}
}
