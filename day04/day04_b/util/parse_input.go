package util

import (
	"bufio"
	"os"
	"strings"
)

type Puzzle struct {
	Lines []string
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		lines = append(lines, line)
	}

	return &Puzzle{
		Lines: lines,
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
