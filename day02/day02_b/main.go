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

func apply_skip(line []int, skip int) []int {
	if skip == -1 {
		return line
	}

	r := make([]int, len(line)-1)
	j := 0
	for i := 0; i < len(line); i++ {
		if i == skip {
			continue
		}
		r[j] = line[i]
		j++
	}

	return r
}

func is_safe(line []int, skip int) bool {
	sign := 99

	line = apply_skip(line, skip)

	for i := 1; i < len(line); i++ {
		delta := line[i-1] - line[i]

		if delta == 0 || util.Abs(delta) > 3 {
			return false
		}

		this_sign := (delta) / util.Abs(delta)

		if sign != 99 && this_sign != sign {
			return false
		}

		sign = this_sign
	}

	return true
}

func solution(input *util.Puzzle, logger *log.Logger) int {
	safe := 0

	for _, line := range input.Lines {
		for skip := -1; skip < len(line); skip++ {
			if is_safe(line, skip) {
				safe++
				goto next_line
			}
		}
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
