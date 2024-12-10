package main

import (
	"image"
	"os"
	"puzzle/util"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

func canMove(pos image.Point, level rune, puzzle *util.Puzzle) bool {
	return pos.Y >= 0 && pos.Y < len(puzzle.Board) &&
		pos.X >= 0 && pos.X < len(puzzle.Board[pos.Y]) &&
		puzzle.Board[pos.Y][pos.X] == level+1
}

func countTrailheads(start image.Point, book map[image.Point]int, puzzle *util.Puzzle) int {
	if val, seen := book[start]; seen {
		return val
	}

	deltas := []image.Point{
		{X: 1, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: -1}, {X: -1, Y: 0},
	}

	level := puzzle.Board[start.Y][start.X]

	for _, d := range deltas {
		pos := start.Add(d)
		if canMove(pos, level, puzzle) {
			if level == '8' {
				// VV: Record a new peak
				book[start] = book[start] + 1
			} else {
				trails := countTrailheads(pos, book, puzzle)
				book[start] = book[start] + trails
			}
		}
	}

	return book[start]
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	trails := 0

	book := map[image.Point]int{}

	for y, row := range puzzle.Board {
		for x, c := range row {
			if c == '0' {
				start := image.Pt(x, y)
				trails += countTrailheads(start, book, puzzle)
			}
		}
	}

	return trails
}

func main() {
	logger := SetupLogger()

	logger.Println("Parse input")
	puzzle, err := util.ReadInputFile(os.Args[1])

	// logger.Println("Input was", input)

	if err != nil {
		logger.Fatalln("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, logger)

	logger.Println("Solution is", sol)
}
