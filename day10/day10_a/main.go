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

func countTrailheads(start image.Point, book map[image.Point]map[image.Point]int, puzzle *util.Puzzle) int {
	if val, seen := book[start]; seen {
		return len(val)
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
				peaks, ok := book[start]
				if !ok {
					peaks = map[image.Point]int{}
				}
				peaks[pos] = 1
				book[start] = peaks
			} else {
				countTrailheads(pos, book, puzzle)

				if book[pos] == nil || len(book[pos]) == 0 {
					continue
				}

				// VV: Start can go to all the peaks that @pos leads to
				peaks, ok := book[start]

				if !ok {
					peaks = map[image.Point]int{}
				}

				for peak := range book[pos] {
					peaks[peak] = 1
				}

				book[start] = peaks
			}
		}
	}

	return len(book[start])
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	trailheads := 0
	// VV: Maintain a record of all the peaks that all points lead to
	book := map[image.Point]map[image.Point]int{}

	for y, row := range puzzle.Board {
		for x, c := range row {
			if c == '0' {
				start := image.Pt(x, y)
				// VV: Report the sum of unique peaks that positions with level '0' lead to
				trailheads += countTrailheads(start, book, puzzle)
			}
		}
	}

	return trailheads
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
