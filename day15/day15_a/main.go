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

func render(puzzle *util.Puzzle, move rune) {
	println(string(move))

	for y, row := range puzzle.Board {
		for x, c := range row {
			if puzzle.Bot.X == x && puzzle.Bot.Y == y {
				print(string(move))
			} else {
				print(string(c))
			}
		}

		println()
	}
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {

	pos := puzzle.Bot

	deltas := map[rune]image.Point{
		'>': image.Pt(+1, +0),
		'<': image.Pt(-1, +0),
		'v': image.Pt(+0, +1),
		'^': image.Pt(+0, -1),
	}

	for _, move := range puzzle.Moves {
		d, ok := deltas[move]

		if !ok {
			log.Panicf("Invalid move %s", string(move))
		}

		/*
			VV: Find the tile ! which will be occupied by either the bot or a box that the bot just moved into the tile

			Some examples:

			- ">.OO" --> ".!OO"
			- ">OO." --> ".OO!"
			- ">#xx" --> "!#xx"
			- ">OO#" --> "!OO#"
		*/

		endPos := pos
		for {
			endPos = endPos.Add(d)

			if puzzle.Board[endPos.Y][endPos.X] != 'O' {
				// VV: There's either an empty spot or a wall
				break
			}
		}

		//VV: The first tile right after the current Bot position
		first := pos.Add(d)

		puzzle.Bot = pos

		if puzzle.Board[endPos.Y][endPos.X] == '#' {
			// VV: when hitting a wall, do nothing
			goto prettyRender
		}

		// VV: We know that we can move N blocks from pos to endPos for distance=1
		// We also know that N-1 blocks are consecutive boxes moving boxes ">123." is the same as swapping 1 with . like so: ">.231"
		puzzle.Board[first.Y][first.X], puzzle.Board[endPos.Y][endPos.X] = puzzle.Board[endPos.Y][endPos.X], puzzle.Board[first.Y][first.X]
		pos = first
		puzzle.Bot = first
	prettyRender:
		// render(puzzle, move)
	}

	render(puzzle, '@')

	gps := 0

	for y, row := range puzzle.Board {
		for x, c := range row {
			if c == 'O' {
				gps += 100*y + x
			}
		}
	}

	return gps
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
