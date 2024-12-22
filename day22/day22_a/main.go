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

func generateRandom(seed uint64, numNumbers int) uint64 {

	for range numNumbers {
		seed = seed ^ (seed*64)%16777216
		seed = ((seed / 32) ^ seed) % 16777216
		seed = ((seed * 2048) ^ seed) % 16777216

	}

	return seed
}

func solution(puzzle *util.Puzzle, numNumbers int, logger *log.Logger) uint64 {
	ret := uint64(0)

	for _, seed := range puzzle.Codes {
		next := generateRandom(seed, numNumbers)
		println(seed, "->", next)

		ret += next
	}

	return ret
}

func main() {
	logger := SetupLogger()

	logger.Println("Parse input")
	puzzle, err := util.ReadInputFile(os.Args[1])

	// logger.Println("Input was", input)

	if err != nil {
		logger.Fatalln("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, 2000, logger)

	logger.Println("Solution is", sol)
}
