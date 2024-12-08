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

func calculateAntinodes(a, b image.Point, antinodes map[image.Point]int, width, height int) {
	/*VV: Doesn't look we need fancy math here. 2 Points are always on a line
	and we can just use manhattan distance to find where "a point A that's twice
	as far to 0 as 1 is" :

	..........
	..........
	..........
	..........
	........0.
	.....1....
	..A.......
	..........
	..........
	..........

	Above, 3 is the antinode of 0 based on 1:
		The delta between 0 and 1 is (dx, dy) = (-3, 1)
		The delta between 0 and A is (-6, 2) = (2*dx, 2*dy)

	Now, we repeat this logic with the only difference that we're constantly increasing the
	distance multiplier (in above example `2`) till we go past the borders of the Board.
	*/

	delta := image.Pt(a.X-b.X, a.Y-b.Y)

	antinodes[a] = 1
	antinodes[b] = 1

	// VV: There are 2 antinode candidates - one "behind" each of the nodes in the pair
	for multiplier := 2; ; multiplier++ { // those behind 1
		anti := a.Sub(delta.Mul(multiplier))
		if anti.Y >= 0 && anti.Y < height && anti.X >= 0 && anti.X < width {
			antinodes[anti] = 1
		} else {
			break
		}
	}

	for multiplier := 2; ; multiplier++ { // those behind 0
		anti := b.Add(delta.Mul(multiplier))
		if anti.Y >= 0 && anti.Y < height && anti.X >= 0 && anti.X < width {
			antinodes[anti] = 1
		} else {
			break
		}
	}
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	antinodes := map[image.Point]int{}

	antennas := map[rune][]image.Point{}

	for y, row := range puzzle.Board {
		for x, c := range row {
			if c != '.' {
				antennas[c] = append(antennas[c], image.Pt(x, y))
			}
		}
	}

	height := len(puzzle.Board)
	width := len(puzzle.Board[0])

	for _, sameAnts := range antennas {
		for i := range sameAnts {
			for j := i + 1; j < len(sameAnts); j++ {
				calculateAntinodes(sameAnts[i], sameAnts[j], antinodes, width, height)
			}
		}
	}

	println("--")
	for y, row := range puzzle.Board {
		for x, _ := range row {
			if antinodes[image.Pt(x, y)] != 0 {
				print("#")
			} else {
				print(".")
			}
		}

		println()
	}

	return len(antinodes)
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
