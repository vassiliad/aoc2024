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

func solution(input *util.Puzzle, logger *log.Logger) int {
	solution := 0

	re := regexp.MustCompile(`mul\((?P<left>\d+),(?P<right>\d+)\)`)

	for _, line := range input.Lines {
		matches := re.FindAll([]byte(line), -1)

		for _, m := range matches {
			op := string(m[4 : len(m)-1])
			parts := strings.Split(op, ",")

			left, _ := strconv.Atoi(parts[0])
			right, _ := strconv.Atoi(parts[1])

			solution += left * right
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
