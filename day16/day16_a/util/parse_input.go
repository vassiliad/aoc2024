package util

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"slices"
	"strings"
)

type Puzzle struct {
	Board      [][]rune
	Start, End image.Point
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	start := image.Pt(-1, -1)
	end := image.Pt(-1, -1)
	board := [][]rune{}

	y := 0

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		row := []rune(line)
		if x := slices.Index(row, 'S'); x > -1 {
			start = image.Pt(x, y)
		} else if x := slices.Index(row, 'E'); x > -1 {
			end = image.Pt(x, y)
		}
		board = append(board, row)
		y++
	}

	if start.X < 0 || start.Y < 0 {
		return &Puzzle{
			Board: board, Start: start, End: end,
		}, fmt.Errorf("invalid start position %+v", start)
	}

	if end.X < 0 || end.Y < 0 {
		return &Puzzle{
			Board: board, Start: start, End: end,
		}, fmt.Errorf("invalid end position %+v", end)
	}

	return &Puzzle{
		Board: board, Start: start, End: end,
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
