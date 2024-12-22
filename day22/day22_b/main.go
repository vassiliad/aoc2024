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

Book is basically a hashmap whose keys are "{diff[0]},{diff[1]},{diff[2]},{diff[3]}" for the 4
digits of the sliding windows over the prices that buyers are willing to spend.

The range of "diff"  is [-9, 9] which means we can use a base 19 number to keep track of the diffs.
We just need to make sure that each diff is positive i.e. add 9
*/
func generateDifferences(seed uint64, numNumbers int, book *[19 * 19 * 19 * 19]uint64) {
	last := int(seed % 10)

	differences := [4]int{}

	localBook := [19 * 19 * 19 * 19]int{}

	for idx := range numNumbers {
		seed = seed ^ (seed<<6)%16777216
		seed = ((seed >> 5) ^ seed) % 16777216
		seed = ((seed << 11) ^ seed) % 16777216

		cur := int(seed % 10)
		diff := cur - last
		differences[idx%4] = diff

		if idx > 2 {
			start := (idx - 3) % 4

			key := uint64(0)
			power := uint64(1)
			// VV: A diff ranges from [-9, 9] so add 9 to it and make it [0, 18]
			// then build the key as if it were a base 19 number
			for i := start; i < start+4; i++ {
				digit := uint64(9 + int64(differences[i%4]))
				key += digit * power
				power *= 19
			}

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

	book := [19 * 19 * 19 * 19]uint64{}

	for _, seed := range puzzle.Codes {
		generateDifferences(seed, numNumbers, &book)
	}

	maxBananas := uint64(0)
	bestDiff := 0
	for differences, bananas := range book {
		if bananas > maxBananas {
			maxBananas = bananas
			bestDiff = differences
		}
	}

	// VV: Convert the base19 number back to differences
	differences := [4]int{}
	for i := 0; i < 4; i, bestDiff = i+1, bestDiff/19 {
		differences[i] = bestDiff%19 - 9
	}
	println("Buy on", fmt.Sprintf("%d,%d,%d,%d", differences[0], differences[1], differences[2], differences[3]), "for", maxBananas)

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
