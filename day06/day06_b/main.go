package main

import (
	"image"
	"os"
	"puzzle/util"
	"time"

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
func nextDirection(puzzle *util.Puzzle, direction Direction, position, obstacle image.Point) Direction {
	for {
		// VV: Try moving forward
		next := position.Add(deltas[direction])

		// VV: if you're about to fall of the map, stop spinning
		if next.Y < 0 || next.Y >= len(puzzle.Board) || next.X < 0 || next.X >= len(puzzle.Board[next.Y]) {
			return -1
		}

		// VV: Turn right whenever you're about to bump into an obstacle
		if (puzzle.Board[next.Y][next.X] != '#') && (next != obstacle) {
			return direction
		}
		direction = (direction + 1) % 4
	}
}

func simulate(guard image.Point, puzzle *util.Puzzle) map[image.Point]Direction {
	visited := map[image.Point]Direction{}

	fake := image.Pt(-1, -1)
	nextDir := Up
	for {
		nextDir = nextDirection(puzzle, nextDir, guard, fake)

		if nextDir == -1 {
			break
		}

		guard = guard.Add(deltas[nextDir])
		visited[guard] = nextDir
	}

	return visited
}

func obstacleCausesLoop(obstacle, guard image.Point, puzzle *util.Puzzle, c chan int) {
	visited := map[image.Point]Direction{}

	nextDir := Up
	for {
		if oldDir, seen := visited[guard]; seen && (oldDir == nextDir) {
			// VV: been on this exact same spot before heading the same direction, therefore this is a loop
			c <- 1
			return
		}

		visited[guard] = nextDir

		nextDir = nextDirection(puzzle, nextDir, guard, obstacle)
		if nextDir == -1 {
			c <- 0
			return
		}

		guard = guard.Add(deltas[nextDir])

	}

}

func worker(guard image.Point, puzzle *util.Puzzle, input chan image.Point, output chan int) {
	for obstacle := range input {
		obstacleCausesLoop(obstacle, guard, puzzle, output)
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

	// VV: We've already mapped every single step the guard makes. Instead of placing blocks in random spots
	// focus on the tiles that the guard is going to walk on
	path := simulate(guard, puzzle)

	logger.Println("Candidates are", len(path))

	// VV: Not smart enough to optimize your solution? throw a bunch of threads in the mix
	now := time.Now()

	start := now.UnixNano()

	const numWorkers = 128

	workerOutput := make(chan int, len(path))
	workerInput := make(chan image.Point, len(path))

	for range numWorkers {
		go worker(guard, puzzle, workerInput, workerOutput)
	}

	for obstacle := range path {
		workerInput <- obstacle
	}

	close(workerInput)

	loops := 0
	for range path {
		loops += <-workerOutput
	}

	logger.Println("Duration", time.Now().UnixNano()-start)

	return loops
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
