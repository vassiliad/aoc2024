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

func canDesign(design string, towels []string, book map[string]bool) bool {
	if feasible, seen := book[design]; seen {
		return feasible
	}

	for _, towel := range towels {
		remaining := strings.TrimPrefix(design, towel)

		if remaining == "" || (remaining != design && canDesign(remaining, towels, book)) {
			book[design] = true
			return true
		} else {
			book[remaining] = false
		}
	}

	book[design] = false

	return false
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	feasible := 0

	book := map[string]bool{}

	for _, design := range puzzle.Designs {
		if canDesign(design, puzzle.Towels, book) {
			feasible++
		}
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
