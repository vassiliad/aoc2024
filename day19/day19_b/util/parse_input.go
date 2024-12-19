package util

import (
	"bufio"
	"os"
	"strings"
)

type Puzzle struct {
	Designs []string
	Towels  []string
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	designs := []string{}
	towels := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if len(towels) == 0 {
			towels = strings.Split(line, ", ")
		} else {
			designs = append(designs, line)
		}
	}

	return &Puzzle{
		Designs: designs, Towels: towels,
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
