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

type Direction int

const (
	Up    Direction = 0
	Right Direction = 1
	Down  Direction = 2
	Left  Direction = 3
)

var deltas = []image.Point{
	image.Pt(0, -1), image.Pt(1, 0), image.Pt(0, 1), image.Pt(-1, 0),
}

/*Returns the next direction or -1 if the guard is about to exit the board*/
func nextDirection(puzzle *util.Puzzle, direction Direction, position image.Point) Direction {
	for {
		// VV: Try moving forward
		next := position.Add(deltas[direction])

		// VV: if you're about to fall of the map, stop spinning
		if next.Y < 0 || next.Y >= len(puzzle.Board) || next.X < 0 || next.X >= len(puzzle.Board[next.Y]) {
			return -1
		}

		// VV: Turn right whenever you're about to bump into an obstacle
		if puzzle.Board[next.Y][next.X] != '#' {
			return direction
		}
		direction = (direction + 1) % 4
	}
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	guard := image.Pt(-1, -1)

	for y, row := range puzzle.Board {
		for x, c := range row {
			if c == '^' {
				guard = image.Pt(x, y)
				break
			}
		}
	}

	if guard.X == -1 {
		panic("No guard here")
	}

	visited := map[image.Point]int{}

	nextDir := Up
	for {
		puzzle.Board[guard.Y][guard.X] = 'x'
		visited[guard] = 1
		nextDir = nextDirection(puzzle, nextDir, guard)

		if nextDir == -1 {
			break
		}

		guard = guard.Add(deltas[nextDir])
	}

	for _, row := range puzzle.Board {
		println(string(row))
	}

	return len(visited)
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
