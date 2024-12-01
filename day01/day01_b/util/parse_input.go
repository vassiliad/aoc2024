package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Puzzle struct {
	Left, Right []int
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	left := []int{}
	right := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		for i, s := range strings.Fields(line) {
			num, err := strconv.Atoi(s)

			if err != nil {
				return &Puzzle{
					Left:  left,
					Right: right,
				}, err
			}

			if i == 0 {
				left = append(left, num)
			} else if i == 1 {
				right = append(right, num)
			} else {
				return &Puzzle{
					Left: left, Right: right,
				}, fmt.Errorf("too many numbers in line: %s", line)
			}
		}
	}

	return &Puzzle{
		Left: left, Right: right,
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
