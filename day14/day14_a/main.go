package main

import (
	"os"
	"puzzle/util"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

func solution(puzzle *util.Puzzle, width, height, seconds int, logger *log.Logger) int {
	topLeft, topRight, bottomLeft, bottomRight := 0, 0, 0, 0

	midWidth := width / 2
	midHeight := height / 2
	for _, bot := range puzzle.Bots {
		x := bot.Position.X + bot.Velocity.X*seconds
		y := bot.Position.Y + bot.Velocity.Y*seconds

		finalX := (x%width + width) % width
		finalY := (y%height + height) % height

		if finalX == midWidth || finalY == midHeight {
			continue
		}

		if finalX < midWidth && finalY < midHeight {
			topLeft++
		} else if finalX > midWidth && finalY < midHeight {
			topRight++
		} else if finalX > midWidth && finalY > midHeight {
			bottomRight++
		} else {
			bottomLeft++
		}
	}

	return topLeft * topRight * bottomLeft * bottomRight
}

func main() {
	logger := SetupLogger()

	logger.Println("Parse input")
	puzzle, err := util.ReadInputFile(os.Args[1])

	// logger.Println("Input was", input)

	if err != nil {
		logger.Fatalln("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, 101, 103, 100, logger)

	logger.Println("Solution is", sol)
}
