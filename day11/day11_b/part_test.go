package main

import (
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `125 17`

	println(small)

	input, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	solution := solution(input, 25, logger)
	const correctAnswer = 55312

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}
}

func TestTiny1(t *testing.T) {
	small := `0 1 10 99 999`

	println(small)

	input, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	solution := solution(input, 1, logger)
	const correctAnswer = 7

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}
}

func TestTiny2(t *testing.T) {
	small := `125 17`

	println(small)

	input, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	solution := solution(input, 6, logger)
	const correctAnswer = 22

	if solution != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", solution)
	}
}

func TestSplit(t *testing.T) {
	tests := map[int][]int{
		1000: {4, 10, 0},
		1010: {4, 10, 10},
		10:   {2, 1, 0},
		99:   {2, 9, 9},
		1:    {1},
		123:  {3},
	}

	for number, parts := range tests {
		digits := util.NumberOfDigits(number)
		if digits != parts[0] {
			t.Fatal(number, "expected digits to be", parts[0], "but it was", digits)
		}

		if digits%2 == 0 {
			left, right := util.SplitNumber(number, digits/2)
			if left != parts[1] {
				t.Fatal(number, parts[1], parts[2], "expected left to be", parts[1], "but it was", left)
			}
			if right != parts[2] {
				t.Fatal(number, parts[1], parts[2], "expected right to be", parts[2], "but it was", right)
			}
		}

	}
}
