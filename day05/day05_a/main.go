package main

import (
	"image"
	"os"
	"puzzle/util"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

func isInOrder(rules []image.Point, update map[int]int) bool {

	for _, r := range rules {
		left, okl := update[r.X]
		right, okr := update[r.Y]

		if (okl && okr) && left > right {
			return false
		}
	}

	return true
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {

	solution := 0

	for _, update := range puzzle.Updates {
		updateOrder := map[int]int{}
		for i, p := range update {
			updateOrder[p] = i
		}

		if isInOrder(puzzle.Rules, updateOrder) {
			solution += update[len(update)/2]
		}
	}

	return solution
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
