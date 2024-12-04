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

func count_stencil(stencil [][]rune, start image.Point, puzzle *util.Puzzle) int {

	if start.Y < 0 || start.Y >= len(puzzle.Lines)-2 || start.X < 0 || start.X >= len(puzzle.Lines[start.Y])-2 {
		return 0
	}

	for y, row := range stencil {
		for x := range row {
			if row[x] != '.' && rune(puzzle.Lines[y+start.Y][x+start.X]) != row[x] {
				return 0
			}
		}
	}

	return 1
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	stencils := [][][]rune{}

	stencils = append(stencils,
		[][]rune{
			[]rune("M.M"),
			[]rune(".A."),
			[]rune("S.S"),
		},
	)

	stencils = append(stencils,
		[][]rune{
			[]rune("S.S"),
			[]rune(".A."),
			[]rune("M.M"),
		},
	)

	stencils = append(stencils,
		[][]rune{
			[]rune("S.M"),
			[]rune(".A."),
			[]rune("S.M"),
		},
	)

	stencils = append(stencils,
		[][]rune{
			[]rune("M.S"),
			[]rune(".A."),
			[]rune("M.S"),
		},
	)

	solution := 0

	for y := range puzzle.Lines {
		for x := range puzzle.Lines[0] {
			for _, s := range stencils {
				solution += count_stencil(s, image.Pt(x, y), puzzle)
			}
		}
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
