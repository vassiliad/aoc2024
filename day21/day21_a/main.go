package main

import (
	"fmt"
	"image"
	"os"
	"puzzle/util"
	"strconv"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

func verifyMovement(buttons [][]rune, code []rune, start image.Point) bool {
	deltas := map[rune]image.Point{
		'>': image.Pt(1, 0),
		'<': image.Pt(-1, 0),
		'v': image.Pt(0, 1),
		'^': image.Pt(0, -1),
	}

	for _, c := range code {
		if c != 'A' {
			if buttons[start.Y][start.X] == 'x' {
				return false
			}
			if delta, ok := deltas[c]; ok {
				start = start.Add(delta)
			} else {
				panic(fmt.Sprintf("unknown button '%s'", string(c)))
			}
			if buttons[start.Y][start.X] == 'x' {
				return false
			}
		}
	}

	return true
}

func computeMovements(buttons [][]rune) map[rune]map[rune][][]rune {
	ret := map[rune]map[rune][][]rune{}

	for y, row := range buttons {
		for x, c := range row {

			if c == 'x' {
				continue
			}

			for otherY, otherRow := range buttons {
				for otherX, otherC := range otherRow {
					if otherC == 'x' {
						continue
					}

					movementsHorizontal := []rune{}
					dx := otherX - x
					{

						if dx != 0 {
							step := 1
							movement := '>'

							if dx < 0 {
								step = -1
								movement = '<'
							}

							for i := x; i != otherX; i += step {
								movementsHorizontal = append(movementsHorizontal, movement)
							}
						}

					}

					movementsVertical := []rune{}
					dy := otherY - y
					{
						if dy != 0 {
							step := 1
							movement := 'v'

							if dy < 0 {
								step = -1
								movement = '^'
							}

							for i := y; i != otherY; i += step {
								movementsVertical = append(movementsVertical, movement)
							}
						}

					}
					movements := [][]rune{}

					if dx != 0 && dy != 0 {
						firstHoriz := []rune{}
						firstHoriz = append(firstHoriz, movementsHorizontal...)
						firstHoriz = append(firstHoriz, movementsVertical...)
						firstHoriz = append(firstHoriz, 'A')

						if verifyMovement(buttons, firstHoriz, image.Pt(x, y)) {
							movements = append(movements, firstHoriz)
						}
					}

					{
						firstVert := []rune{}
						firstVert = append(firstVert, movementsVertical...)
						firstVert = append(firstVert, movementsHorizontal...)
						firstVert = append(firstVert, 'A')

						if verifyMovement(buttons, firstVert, image.Pt(x, y)) {
							movements = append(movements, firstVert)
						}
					}

					if maps, ok := ret[c]; ok {
						maps[otherC] = movements
					} else {
						ret[c] = map[rune][][]rune{
							otherC: movements,
						}
					}
				}
			}
		}
	}

	return ret
}

type Memo struct {
	start, stop rune
	level       int
}

/*
VV: Starting from the symbol @start move to the symbol @stop and count how many buttons the
outtermost robot ended up pressing.

There are 2 optimal ways to move from @start to @stop in 1 level:

a) Move horizontally and then vertically
b) Move vertically and then horizontally

This way the level+1 robot does not need to move its arm, it just presses A.

The question is then how to pick the best out of the 2 movements.

The function uses memoization to figure out the best way. At every level, it tries out all
working movesets (i.e. those which avoid the gap) and then picks the one that ends up making
the least amount of movements in the outtermost level.
*/
func countCost(
	start, stop rune, level, maxLevel int, movementsFinal, movements map[rune]map[rune][][]rune, book map[Memo]uint64) uint64 {

	memo := Memo{start: start, stop: stop, level: level}

	if cost, seen := book[memo]; seen {
		return cost
	}

	myMovements := movementsFinal[start][stop]
	if level != 0 {
		myMovements = movements[start][stop]
	}

	const huge = ^uint64(0)
	bestCost := huge

	for _, moveSet := range myMovements {
		thisCost := uint64(0)

		last := 'A'
		if level < maxLevel {
			for _, m := range moveSet {
				thisCost += countCost(last, m, level+1, maxLevel, movementsFinal, movements, book)
				last = m
			}
		} else {
			thisCost = uint64(len(moveSet))
		}

		if thisCost < bestCost {
			bestCost = thisCost
		}
	}

	if bestCost == huge {
		panic("no movesets here")
	}

	book[memo] = bestCost
	return bestCost
}

func solution(puzzle *util.Puzzle, logger *log.Logger) uint64 {
	const numRobots = 2

	buttonsFinal := [][]rune{
		{'7', '8', '9'},
		{'4', '5', '6'},
		{'1', '2', '3'},
		{'x', '0', 'A'},
	}

	buttons := [][]rune{
		{'x', '^', 'A'},
		{'<', 'v', '>'},
	}

	movementsFinal := computeMovements(buttonsFinal)
	movements := computeMovements(buttons)

	book := map[Memo]uint64{}
	ret := uint64(0)
	for _, code := range puzzle.Codes {
		num, err := strconv.Atoi(string(code[0 : len(code)-1]))
		if err != nil {
			panic("could not decode code")
		}
		print(string(code), "(", num, "):\n")

		last := 'A'
		total := uint64(0)

		for _, c := range code {
			total += countCost(last, c, 0, numRobots, movementsFinal, movements, book)
			last = c
		}

		println(total, "*", num, "=", total*uint64(num))
		ret += total * uint64(num)
	}

	return ret
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
