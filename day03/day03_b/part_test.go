package main

import (
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))`

	println(small)

	input, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	solution := solution(input, logger)
	const correct_answer = 161

	if solution != correct_answer {
		t.Fatal("Expected answer to be", correct_answer, "but it was", solution)
	}
}

func TestWithDont(t *testing.T) {
	small := `xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`

	println(small)

	input, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", input)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	solution := solution(input, logger)
	const correct_answer = 48

	if solution != correct_answer {
		t.Fatal("Expected answer to be", correct_answer, "but it was", solution)
	}
}
