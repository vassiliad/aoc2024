package main

import (
	"fmt"
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

func renderBoard(board map[image.Point]int, spot image.Point, width, height int) {
	for y := range height {
		for x := range width {
			pos := image.Pt(x, y)

			if spot == pos {
				print("*")
			} else if board[pos] > 0 {
				print(fmt.Sprintf("%x", board[pos]))
			} else {
				print(".")
			}
		}

		println()
	}
}

func workerSimulate(board [][]int, width, height, seconds int, input chan *util.Bot, output chan int) {
	for bot := range input {
		x := bot.Position.X + bot.Velocity.X*seconds
		y := bot.Position.Y + bot.Velocity.Y*seconds

		x = (x%width + width) % width
		y = (y%height + height) % height

		board[y][x]++

		output <- 0
	}
}

/*
Find the position of bots at time @seconds. A tree will have a square in it somewhere,
whenever one of the bots end up inside one of those squares that indicates that the bots
just created the Tree Christmas egg formation.

In my input, a square whose size is 3x3 was good enough as a test (i.e. no false positives).
*/
func simulate(puzzle *util.Puzzle, width, height, seconds int) bool {

	board := make([][]int, height)

	for y := range height {
		board[y] = make([]int, width)
	}
	const blockSize = 3

	simulate := func(idx int) (int, int) {
		bot := puzzle.Bots[idx]
		x := bot.Position.X + bot.Velocity.X*seconds
		y := bot.Position.Y + bot.Velocity.Y*seconds

		x = (x%width + width) % width
		y = (y%height + height) % height

		board[y][x] = board[y][x] + 1

		return x, y
	}

	idx := 0
	for ; idx < len(puzzle.Bots); idx++ {
		x, y := simulate(idx)

		for dy := 0; dy < blockSize; dy++ {
			for dx := 0; dx < blockSize; dx++ {
				if y+dy < height && y+dy >= 0 && x+dx >= 0 && x+dx < width {
					if board[y+dy][x+dx] == 0 {
						goto next
					}
				} else {
					goto next
				}

			}
		}

		goto christmasEgg

	next:
	}

	return false

christmasEgg:

	// VV: We may have had an early exit, fill in the spots of the remaining bots
	for ; idx < len(puzzle.Bots); idx++ {
		simulate(idx)
	}

	println("Christmas egg!")
	for y := range height {
		for x := range width {
			if board[y][x] != 0 {
				print(fmt.Sprintf("%x", board[y][x]))
			} else {
				print(".")
			}
		}
		println()
	}

	return true
}

func solution(puzzle *util.Puzzle, width, height int, logger *log.Logger) int {

	for seconds := 1; ; seconds++ {
		if simulate(puzzle, width, height, seconds) {
			return seconds
		}
	}

	panic("oh no")
}

func main() {
	logger := SetupLogger()

	logger.Println("Parse input")
	puzzle, err := util.ReadInputFile(os.Args[1])

	// logger.Println("Input was", input)

	if err != nil {
		logger.Fatalln("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, 101, 103, logger)

	logger.Println("Solution is", sol)
}
