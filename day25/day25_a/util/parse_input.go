package util

import (
	"bufio"
	"os"
	"strings"
)

type Puzzle struct {
	Keys  [][]int
	Locks [][]int
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	keys := [][]int{}
	locks := [][]int{}

	firstRow := []rune{}
	current := []int{}
	y := 0

	addLockOrKey := func() {
		if len(firstRow) > 0 {
			if firstRow[0] == '#' {
				locks = append(locks, current)
			} else {
				for idx := range current {
					current[idx] = y - current[idx] - 2
				}
				keys = append(keys, current)
			}
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			if len(firstRow) == 0 {
				// VV: this is just an empty line
				y = 0
				continue
			}

			addLockOrKey()

			y = 0
			firstRow = []rune{}
			current = []int{}
			continue
		}
		y++

		row := []rune(line)
		if len(firstRow) == 0 {
			firstRow = row
			for range firstRow {
				current = append(current, 0)
			}
		} else {
			for idx, val := range row {
				if val == firstRow[idx] {
					current[idx]++
				}
			}
		}
	}

	addLockOrKey()

	return &Puzzle{
		Keys: keys, Locks: locks,
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
