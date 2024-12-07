package main

import (
	"container/list"
	"math"
	"os"
	"puzzle/util"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

type State struct {
	opIndex   int
	aggregate int
}

func guessCalibration(equation util.Equation) bool {
	target := equation.Target
	numbers := equation.Numbers

	pending := list.New()

	pending.PushFront(State{opIndex: 1, aggregate: numbers[0]})

	for pending.Front() != nil {
		f := pending.Front()
		pending.Remove(f)

		cur := f.Value.(State)
		{
			plus := cur
			plus.aggregate += numbers[plus.opIndex]

			if plus.aggregate == target && plus.opIndex == len(numbers)-1 {
				return true
			} else if plus.aggregate <= target {
				plus.opIndex++
				if plus.opIndex < len(numbers) {
					pending.PushFront(plus)
				}
			}
		}

		{
			mult := cur
			mult.aggregate *= numbers[mult.opIndex]

			if mult.aggregate == target && mult.opIndex == len(numbers)-1 {
				return true
			} else if mult.aggregate <= target {
				mult.opIndex++
				if mult.opIndex < len(numbers) {
					pending.PushFront(mult)
				}
			}
		}

		{
			join := cur
			digits := util.NumberOfDigits(numbers[join.opIndex])

			join.aggregate = join.aggregate*int(math.Pow10(digits)) + numbers[join.opIndex]

			if join.aggregate == target && join.opIndex == len(numbers)-1 {
				return true
			} else if join.aggregate <= target {
				join.opIndex++
				if join.opIndex < len(numbers) {
					pending.PushFront(join)
				}
			}
		}
	}

	return false
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {

	sol := 0

	for _, equation := range puzzle.Equations {
		if guessCalibration(equation) {
			sol += equation.Target
		}
	}

	return sol
}

func main() {
	logger := SetupLogger()

	logger.Println("Parse input")
	puzzle, err := util.ReadInputFile(os.Args[1])

	// logger.Println("Input was", input)

	if err != nil {
		logger.Fatalln("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, logger)

	logger.Println("Solution is", sol)
}
