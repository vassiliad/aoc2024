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
	Bot   image.Point
	Board [][]rune
	Moves []rune
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	bot := image.Pt(-1, -1)
	board := [][]rune{}
	moves := []rune{}
	y := 0

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			row := []rune(line)

			x := slices.Index(row, '@')

			if x != -1 {
				bot = image.Pt(x, y)
				row[x] = '.'
			}

			board = append(board, row)
			y++
		} else {
			moves = append(moves, []rune(line)...)
		}
	}

	if bot.X == -1 || bot.Y == -1 {
		return &Puzzle{
			Bot: bot, Board: board, Moves: moves,
		}, fmt.Errorf("no bot")
	}

	return &Puzzle{
		Bot: bot, Board: board, Moves: moves,
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
