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

func count_words(word string, start image.Point, dir image.Point, puzzle *util.Puzzle) int {
	pos := start

	for i := range word {
		if pos.Y < 0 || pos.Y >= len(puzzle.Lines) || pos.X < 0 || pos.X >= len(puzzle.Lines[pos.Y]) {
			return 0
		}

		// fmt.Printf("%d-%q %q %q\n", i, pos, puzzle.Lines[pos.Y][pos.X], word[i])
		if puzzle.Lines[pos.Y][pos.X] != word[i] {
			return 0
		}

		if i == len(word)-1 {
			return 1
		}

		pos = pos.Add(dir)
	}

	return 0
}

func paint(board [][]rune, word string, start, dir image.Point) {
	pos := start

	for i := range word {
		board[pos.Y][pos.X] = rune(word[i])
		pos = pos.Add(dir)
	}
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	const word = "XMAS"
	const rev_word = "SAMX"

	solution := 0
	forward := []image.Point{
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 0, Y: 1},
		{X: 1, Y: -1},
	}

	board := [][]rune{}
	for y := range len(puzzle.Lines) {
		line := make([]rune, len(puzzle.Lines[y]))
		for x := range line {
			line[x] = '.'
		}
		board = append(board, line)
	}

	// VV: For each point on the board look for 2 words but only going forward
	// this way there's no need for a mechanism to avoid counting the same word twice
	for y := range len(puzzle.Lines) {

		for x := range len(puzzle.Lines[0]) {
			pos := image.Pt(x, y)

			for _, delta := range forward {
				if count_words(word, pos, delta, puzzle) > 0 {
					paint(board, word, pos, delta)
					solution++
				}
				if count_words(rev_word, pos, delta, puzzle) > 0 {
					paint(board, rev_word, pos, delta)
					solution++
				}
			}
		}
	}

	for _, line := range board {
		println(string(line))
	}

	return solution
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
