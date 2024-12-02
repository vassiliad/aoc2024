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

func solution(input *util.Puzzle, logger *log.Logger) int {
	safe := 0

	for _, line := range input.Lines {

		sign := 99

		for i := 1; i < len(line); i++ {
			delta := line[i-1] - line[i]

			if delta == 0 || util.Abs(delta) > 3 {
				goto next_line
			}

			this_sign := (delta) / util.Abs(delta)

			if sign != 99 && this_sign != sign {
				goto next_line
			}

			sign = this_sign
		}

		safe++
	next_line:
	}

	return safe
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
