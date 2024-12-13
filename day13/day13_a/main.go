package main

import (
	"os"
	"puzzle/util"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

func gamble(gamba util.Gamba) (int, int) {
	pressA, pressB := 0, 0

	x1 := gamba.ButtonA.X
	y1 := gamba.ButtonA.Y

	x2 := gamba.ButtonB.X
	y2 := gamba.ButtonB.Y

	x := gamba.Prize.X
	y := gamba.Prize.Y

	divident := x1*y - y1*x
	divisor := x1*y2 - x2*y1
	pressB = divident / divisor

	if divisor != 0 && (divident%divisor != 0) {
		return 0, 0
	}

	divident = y - pressB*y2

	if y1 != 0 && divident%y1 != 0 {
		return 0, 0
	}

	pressA = divident / y1

	return pressA, pressB
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	tokens := 0

	for _, gamba := range puzzle.Gambas {
		pressA, pressB := gamble(gamba)

		t := pressA*3 + pressB*1
		if pressA > 0 && pressB > 0 {
			logger.Printf("It would take %d A and %d B buttons for a total of %d tokens to win %+v\n", pressA, pressB, t, gamba)
			tokens += t
		} else {
			logger.Printf("Cannot beat the claw for %+v\n", gamba)
		}

	}

	return tokens
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
