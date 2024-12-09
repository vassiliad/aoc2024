package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Puzzle struct {
	BlockSizes []int
	Free       []int
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	blockSizes := []int{}
	free := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		for i, c := range line {
			num := int(c - '0')
			if (num < 0) || (num > 9) {
				return &Puzzle{BlockSizes: blockSizes, Free: free}, fmt.Errorf("unexpected character %d", c)
			}

			if i%2 == 0 {
				blockSizes = append(blockSizes, num)
			} else {
				free = append(free, num)
			}
		}
	}

	return &Puzzle{
		BlockSizes: blockSizes,
		Free:       free,
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
