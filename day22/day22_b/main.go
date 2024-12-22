package main

import (
	"fmt"
	"os"
	"puzzle/util"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

/*
VV: Produce the @numNumbers prices and maintain a record of all 4-window differences.

Whenever a difference shows up for the 1st time, record how much this buyer would pay into @book
After calling this method for all buyers, book[key] will record how much you'd get if you asked
the monkey to sell on @key differences
*/
func generateDifferences(seed uint64, numNumbers int, book map[string]uint64) {
	last := int(seed % 10)

	differences := []int{}

	localBook := map[string]uint64{}

	for range numNumbers {
		seed = seed ^ (seed*64)%16777216
		seed = ((seed / 32) ^ seed) % 16777216
		seed = ((seed * 2048) ^ seed) % 16777216

		cur := int(seed % 10)

		if len(differences) == 4 {
			differences = differences[1:4]
		}

		differences = append(differences, cur-last)

		if len(differences) == 4 {
			key := fmt.Sprintf("%d,%d,%d,%d", differences[0], differences[1], differences[2], differences[3])
			// VV: If this is the 1st time that this difference shows up, make a record of it and note down how much
			// bananas you'd get from selling the hiding spot to this buyer
			if localBook[key] == 0 {
				book[key] += uint64(cur)
				localBook[key] += 1
			}
		}

		last = cur
	}
}

func solution(puzzle *util.Puzzle, logger *log.Logger) uint64 {
	const numNumbers = 2000

	book := map[string]uint64{}

	for _, seed := range puzzle.Codes {
		generateDifferences(seed, numNumbers, book)
	}

	maxBananas := uint64(0)
	bestDiff := ""
	for differences, bananas := range book {
		// println("diff", differences, "bananas", bananas)
		if bananas > maxBananas {
			maxBananas = bananas
			bestDiff = differences
		}
	}

	println("Buy on", bestDiff, "for", maxBananas)

	return maxBananas
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
