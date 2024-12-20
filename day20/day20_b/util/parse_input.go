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
	board := [][]rune{}
	start := image.Pt(-1, -1)
	end := image.Pt(-1, -1)

	y := 0

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		row := []rune(line)
		board = append(board, row)

		if x := slices.Index(row, 'S'); x > -1 {
			start = image.Pt(x, y)
		}

		if x := slices.Index(row, 'E'); x > -1 {
			end = image.Pt(x, y)
		}
		y++
	}

	if start.X == -1 {
		return &Puzzle{
			Board: board,
			Start: start,
			End:   end,
		}, fmt.Errorf("missing the start point")
	}

	if end.X == -1 {
		return &Puzzle{
			Board: board,
			Start: start,
			End:   end,
		}, fmt.Errorf("missing the end point")
	}

	return &Puzzle{
		Board: board,
		Start: start,
		End:   end,
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
