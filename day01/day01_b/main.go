package main

import (
	"os"
	"puzzle/util"
	"slices"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

func solution(input *util.Puzzle, logger *log.Logger) int {
	slices.Sort(input.Left)
	slices.Sort(input.Right)

	count := map[int]int{}

	for _, right := range input.Right {
		count[right] = count[right] + 1
	}

	ret := 0
	for _, left := range input.Left {
		ret += left * count[left]
	}

	return ret
}

func main() {
	logger := SetupLogger()

	logger.Println("Parse input")
	input, err := util.ReadInputFile(os.Args[1])

	// logger.Println("Input was", input)

	if err != nil {
		logger.Fatalln("Ran into problems while reading input. Problem", err)
	}

	sol := solution(input, logger)

	logger.Println("Solution is", sol)
}
