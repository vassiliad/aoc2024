package main

import (
	"os"
	"puzzle/util"
	"regexp"
	"strconv"
	"strings"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

func get_product(m string) int {
	op := string(m[4 : len(m)-1])
	parts := strings.Split(op, ",")

	left, _ := strconv.Atoi(parts[0])
	right, _ := strconv.Atoi(parts[1])

	return right * left
}

func solution(input *util.Puzzle, logger *log.Logger) int {
	solution := 0

	rmul := regexp.MustCompile(`^mul\((?P<left>\d+),(?P<right>\d+)\)`)

	do := true
	for _, line := range input.Lines {
		remaining := line

		for len(remaining) > 0 {
			if strings.HasPrefix(remaining, "do()") {
				do = true

				remaining = remaining[len("do()"):]
			} else if strings.HasPrefix(remaining, "don't()") {
				do = false

				remaining = remaining[len("don't()"):]
			} else {
				match := rmul.Find([]byte(remaining))

				if match == nil {
					remaining = remaining[1:]
					continue
				}

				if do {
					solution += get_product(string(match))
				}

				remaining = remaining[len(match):]
			}
		}
	}

	return solution
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
