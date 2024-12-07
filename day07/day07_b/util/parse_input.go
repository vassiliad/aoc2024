package util

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Equation struct {
	Target  int
	Numbers []int
}

type Puzzle struct {
	Equations []Equation
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	equations := []Equation{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}
		testAndParts := strings.Split(line, ": ")

		target, err := strconv.Atoi(testAndParts[0])

		if err != nil {
			return &Puzzle{
				Equations: equations,
			}, err
		}

		numbers := []int{}
		for _, p := range strings.Split(testAndParts[1], " ") {
			num, err := strconv.Atoi(p)

			if err != nil {
				return &Puzzle{
					Equations: equations,
				}, err
			}

			numbers = append(numbers, num)
		}

		equations = append(equations, Equation{Numbers: numbers, Target: target})
	}

	return &Puzzle{
		Equations: equations,
	}, scanner.Err()
}

func ReadString(text string) (*Puzzle, error) {
	scanner := bufio.NewScanner(strings.NewReader(text))

	return ReadScanner(scanner)
}

func ReadInputFile(path string) (*Puzzle, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	return ReadScanner(scanner)
}
