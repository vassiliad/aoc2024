package main

import (
	"os"
	"puzzle/util"
	"strings"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

func possibleDesigns(design string, towels []string, book map[string]int) int {
	if numFeasible, seen := book[design]; seen {
		return numFeasible
	}

	numFeasible := 0

	for _, towel := range towels {
		remaining := strings.TrimPrefix(design, towel)

		if remaining == "" {
			numFeasible++
			continue
		}

		// VV: Doesn't have the @towel prefix
		if remaining == design {
			continue
		}

		numFeasible += possibleDesigns(remaining, towels, book)
	}

	book[design] = numFeasible
	return numFeasible
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	feasible := 0

	book := map[string]int{}

	for _, design := range puzzle.Designs {
		feasible += possibleDesigns(design, puzzle.Towels, book)
	}

	return feasible
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
