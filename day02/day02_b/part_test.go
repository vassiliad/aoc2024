package main

import (
	"puzzle/util"
	"reflect"
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
	const correct_answer = 4

	if solution != correct_answer {
		t.Fatal("Expected answer to be", correct_answer, "but it was", solution)
	}
}

func TestSkip(t *testing.T) {
	solution := apply_skip([]int{8, 6, 4, 4, 1}, 2)
	correct_answer := []int{8, 6, 4, 1}

	if !reflect.DeepEqual(solution, correct_answer) {
		t.Fatalf("Expected answer to be %+v but it was %+v", correct_answer, solution)
	}

	if is_safe([]int{8, 6, 4, 4, 1}, 2) == false {
		t.Fatal("Line is safe")
	}
}
