package util

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Puzzle struct {
	Lines [][]int
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	lines := [][]int{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		numbers := []int{}

		for _, s := range strings.Fields(line) {
			num, err := strconv.Atoi(s)

			numbers = append(numbers, num)

			if err != nil {
				lines = append(lines, numbers)
				return &Puzzle{
					Lines: lines,
				}, err
			}
		}

		lines = append(lines, numbers)
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
