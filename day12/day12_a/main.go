package main

import (
	"container/list"
	"os"
	"puzzle/util"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

type Point struct {
	x, y, direction int
}

func (p *Point) Add(other Point, direction int) Point {
	return Point{x: p.x + other.x, y: p.y + other.y, direction: direction}
}

func fenceTheLotOfThem(point Point, board [][]rune, membership [][]int, groupId int) (int, int) {
	flower := board[point.y][point.x]

	deltas := []Point{
		{x: 1, y: 0}, {x: 0, y: 1}, {x: -1, y: 0}, {x: 0, y: -1},
	}

	pending := list.New()
	pending.PushFront(point)

	/*
		VV: Maintain a record of all cardinal-adjacent plots with different flowers
		Need to also keep track the direction of the vector connecting the 2 adjacent plots to account for "corners".
		An adjacent plot P requires 1 fence for each of the plots in the current region that are touching it.
	*/
	border := map[Point]int{}
	area := 0

	for pending.Front() != nil {
		front := pending.Front()
		pending.Remove(front)

		pos := front.Value.(Point)

		// VV: out of bounds check
		if !(pos.y >= 0 && pos.y < len(board) && pos.x >= 0 && pos.x < len(board[pos.y])) ||
			// VV: Check that this is an adjacent plot with different flowers
			board[pos.y][pos.x] != flower {
			border[pos] = 1
			continue
		}

		if membership[pos.y][pos.x] != 0 {
			continue
		}

		area++
		membership[pos.y][pos.x] = groupId

		for direction, d := range deltas {
			neighbour := pos.Add(d, direction)
			pending.PushFront(neighbour)
		}
	}

	return area, len(border)
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	cost := 0
	membership := [][]int{}

	nextGroupID := 1
	for _, row := range puzzle.Board {
		membership = append(membership, make([]int, len(row)))
	}

	for y, row := range puzzle.Board {
		for x := range row {
			if membership[y][x] != 0 {
				continue
			}

			pos := Point{x: x, y: y}
			area, perimeter := fenceTheLotOfThem(pos, puzzle.Board, membership, nextGroupID)
			nextGroupID++

			t := area * perimeter
			logger.Printf("A region of %s plants with price %d * %d = %d\n", string(puzzle.Board[y][x]), area, perimeter, t)

			cost += t
		}

	}

	alNum := []rune{}
	for a := 'A'; a <= 'Z'; a++ {
		alNum = append(alNum, a)
	}

	for a := 'a'; a <= 'z'; a++ {
		alNum = append(alNum, a)
	}

	for a := '0'; a <= '0'; a++ {
		alNum = append(alNum, a)
	}

	alNum = append(alNum, []rune("!\"£$%^&*()-_=`¬#[]{};:~,<.>/?\\|+")...)
	for _, row := range membership {
		for _, c := range row {
			print(string(alNum[(c-1)%len(alNum)]))
		}
		println()
	}

	return cost
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
