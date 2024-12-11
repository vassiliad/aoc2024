package util

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Puzzle struct {
	Stones []int
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	stones := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		for _, s := range strings.Fields(line) {
			num, err := strconv.Atoi(s)

			if err != nil {
				return &Puzzle{
					Stones: stones,
				}, err
			}

			stones = append(stones, num)
		}
	}

	return &Puzzle{
		Stones: stones,
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
