package main

import (
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `
Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0
`

	// println(small)

	input, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	solution := solution(input, logger)
	const correctAnswer = "4,6,3,5,6,3,5,2,1,0"

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}
}

func TestTiny1(t *testing.T) {
	input := util.Puzzle{
		Registers: [3]int64{0, 0, 9},
		Program:   []int64{2, 6},
	}
	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	solution(&input, logger)
	solution := input.Registers[1]
	const correctAnswer int64 = 1

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}
}

func TestTiny2(t *testing.T) {
	input := util.Puzzle{
		Registers: [3]int64{10, 0, 0},
		Program:   []int64{5, 0, 5, 1, 5, 4},
	}
	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	solution := solution(&input, logger)
	const correctAnswer = "0,1,2"

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}
}

func TestTiny3(t *testing.T) {
	input := util.Puzzle{
		Registers: [3]int64{2024, 0, 0},
		Program:   []int64{0, 1, 5, 4, 3, 0},
	}
	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	solution := solution(&input, logger)
	const correctAnswer = "4,2,5,6,7,7,7,7,3,1,0"

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}

	if input.Registers[0] != 0 {
		t.Fatal("Expected register A to be  0 but it was", input.Registers[0])
	}
}

func TestTiny4(t *testing.T) {
	input := util.Puzzle{
		Registers: [3]int64{0, 29, 0},
		Program:   []int64{1, 7},
	}
	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	solution(&input, logger)
	solution := input.Registers[1]
	const correctAnswer int64 = 26

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}
}

func TestTiny5(t *testing.T) {
	input := util.Puzzle{
		Registers: [3]int64{0, 2024, 43690},
		Program:   []int64{4, 0},
	}
	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	solution(&input, logger)
	solution := input.Registers[1]
	const correctAnswer int64 = 44354

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}
}
