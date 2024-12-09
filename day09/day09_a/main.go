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

func checksum(blockSizes, free []int) int {
	total := 0

	free_idx := 0
	left_idx := 0
	right_idx := len(blockSizes) - 1

	for idx := 0; ; {
		if blockSizes[left_idx] == 0 && blockSizes[right_idx] == 0 && right_idx == 0 {
			break
		}

		if blockSizes[left_idx] > 0 {
			this := idx * left_idx
			total += this

			idx++
			blockSizes[left_idx]--
		} else if free[free_idx] > 0 {
			if blockSizes[right_idx] > 0 {
				// VV: There's enough free space, therefore we can consume from the end
				this := idx * right_idx
				total += this

				idx++
				blockSizes[right_idx]--
				free[free_idx]--

				if blockSizes[right_idx] == 0 {
					// VV: Consumed the entire right-most file, move on to the next one
					right_idx = max(0, right_idx-1)
				}
			} else {
				right_idx = max(0, right_idx-1)
			}
		} else {
			// VV: If we got here, then all the "free" blocks in between files are filled in and all the
			// files from the left have been checksumed. We can only continue towards the right
			left_idx = min(len(blockSizes)-1, left_idx+1)
			free_idx = min(len(free)-1, free_idx+1)
		}
	}

	return total
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {

	return checksum(puzzle.BlockSizes, puzzle.Free)
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
