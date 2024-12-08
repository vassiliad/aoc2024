package util

import (
	"bufio"
	"os"
	"strings"
)

type Puzzle struct {
	Board [][]rune
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	board := [][]rune{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		board = append(board, []rune(line))
	}

	return &Puzzle{
		Board: board,
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
