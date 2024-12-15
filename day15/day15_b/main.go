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

func horizontalMove(puzzle *util.Puzzle, pos, d image.Point) bool {
	endPos := pos
	for {
		endPos = endPos.Add(d)

		if puzzle.Board[endPos.Y][endPos.X] != '[' && puzzle.Board[endPos.Y][endPos.X] != ']' {
			// VV: There's either an empty spot or a wall
			break
		}
	}

	//VV: The first tile right after the current Bot position
	first := pos.Add(d)

	if puzzle.Board[endPos.Y][endPos.X] == '#' {
		// VV: when hitting a wall, do nothing
		return false
	}

	// VV: Move everything 1 point horizontally
	row := puzzle.Board[first.Y]
	for x := endPos.X; x != first.X-d.X; x -= d.X {
		row[x-d.X], row[x] = row[x], row[x-d.X]
	}

	return true
}

/*
Call this function using the position and direction of a bot or box.

A bot can move if it's heading to an empty tile OR if the box is pushing can move.

A box can move if each of its halves can move on 1 empty tile each. If the box is about to collide with
any other boxes then this box can move if the other boxes can also move.

If a box cannot move, then this function will execute exactly 1 for that box. If it can move, it will execute
twice. The first verticalMove() will run the tests (i.e. dryRun mode). The second call will first
ask the boxes lying on the path of this box to move and once they're done then this box will move on the
2 tiles that are now empty "in front" of this box.
*/
func verticalMove(puzzle *util.Puzzle, pos, d image.Point, dryRun bool) bool {
	//VV: The next tile right after the current Bot position
	next := pos.Add(d)
	isBox := puzzle.Board[pos.Y][pos.X] == '['

	if puzzle.Board[pos.Y][pos.X] == ']' {
		log.Panicf("Cannot point to the right half of a box")
	}

	nextRow := puzzle.Board[next.Y]

	// VV: By convention, when moving a box, @pos will always point to the left half of the box
	freeLeft := nextRow[next.X] == '.'
	freeRight := nextRow[next.X+1] == '.'

	blockLeft := nextRow[next.X] == '#'
	blockight := nextRow[next.X+1] == '#'

	// VV: Stop if the next tile is a '#'. If the thing that's moving is a box then also
	// check the right half of the box
	if (blockLeft) || (isBox && blockight) {
		return false
	}

	if !isBox {
		if freeLeft {
			return true
		} else {
			// VV: We can only get here if there's a box in front of us
			boxPos := next
			// VV: make sure we're pointing to the left of the box
			if puzzle.Board[next.Y][next.X] == ']' {
				boxPos.X--
			}

			// VV: The bot can move IF the box it's pointing at can also move (first check with a dry-run)
			if verticalMove(puzzle, boxPos, d, true) {
				verticalMove(puzzle, boxPos, d, false)
				return true
			}
			return false
		}
	} else {
		// VV: A box can move if both of its halves can move

		if freeLeft && freeRight {
			goto commit
		}

		// VV: We get here, if this box is blocked by 1 or 2 boxes

		if !freeLeft && !freeRight {
			// VV: There's exaclty 1 box which is aligned with the current one
			if nextRow[pos.X] == '[' {
				if verticalMove(puzzle, next, d, dryRun) {
					goto commit
				}
				return false
			}

			// VV: There are 2 boxes blocking this one
			left := next
			left.X--

			right := next
			right.X++

			if nextRow[left.X] != '[' || nextRow[right.X] != '[' {
				panic("unreachable")
			}

			// VV: If both boxes can move, then so can this one
			if verticalMove(puzzle, left, d, dryRun) && verticalMove(puzzle, right, d, dryRun) {
				goto commit
			}
		} else if !freeLeft && freeRight {
			// VV: there's 1 box which is NOT aligned with this one i.e. it's 1 to the right
			left := next
			left.X--

			if nextRow[left.X] != '[' {
				panic("unreachable")
			}

			if verticalMove(puzzle, left, d, dryRun) {
				goto commit
			}
		} else if freeLeft && !freeRight {
			// VV: there's 1 box which is NOT aligned with this one i.e. it's 1 to the left
			right := next
			right.X++

			if nextRow[right.X] != '[' {
				panic("unreachable")
			}

			if verticalMove(puzzle, right, d, dryRun) {
				goto commit
			}
		} else {
			panic("unreachable")
		}
	}

	return false

commit:
	row := puzzle.Board[pos.Y]

	performMove := func() {
		x := pos.X
		row[x], nextRow[x] = nextRow[x], row[x]
		row[x+1], nextRow[x+1] = nextRow[x+1], row[x+1]
	}

	if !dryRun {
		performMove()
	}
	return true
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {

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

		if (d.Y == 0 && horizontalMove(puzzle, puzzle.Bot, d)) || (d.Y != 0 && verticalMove(puzzle, puzzle.Bot, d, false)) {
			puzzle.Bot = puzzle.Bot.Add(d)
		}
	}

	render(puzzle, '@')

	gps := 0

	for y, row := range puzzle.Board {
		for x, c := range row {
			if c == '[' {
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
