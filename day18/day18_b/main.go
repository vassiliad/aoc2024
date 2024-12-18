package main

import (
	"container/heap"
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

type State struct {
	pos image.Point
}

var MAGIC_NUMBER = 9_999_999_999_999

func prettyPrint(board [][]int) {
	for _, row := range board {
		for _, t := range row {
			if t == 0 {
				print(".")
			} else if t == MAGIC_NUMBER {
				print("O")
			} else {
				print("#")
			}
		}
		println()
	}
}

func constructBoard(puzzle *util.Puzzle, width, height, cutoff int) [][]int {
	board := [][]int{}

	for range height {
		row := make([]int, width)
		board = append(board, row)
	}

	for _, b := range puzzle.Bytes[:cutoff] {
		board[b.Y][b.X] = 1
	}

	return board
}

/*
Returns True if there's a path from start to end
*/
func findPath(start, end image.Point, board [][]int, width, height int) bool {
	gScore := map[image.Point]int{start: 0}

	reversePath := map[image.Point]image.Point{}
	open := make(util.PriorityQueue, 1)

	open[0] = &util.HeapItem{
		Value:    start,
		Priority: 0,
	}
	heap.Init(&open)

	for open.Len() > 0 {
		heapItem := heap.Pop(&open)

		item := heapItem.(*util.HeapItem)

		cur := item.Value.(image.Point)

		if cur == end {
			return true
		}

		curGScore := gScore[cur]
		nextCost := curGScore + 1
		deltas := []image.Point{
			image.Pt(+1, +0),
			image.Pt(-1, +0),
			image.Pt(+0, -1),
			image.Pt(+0, +1),
		}

		for _, d := range deltas {
			next := cur.Add(d)

			if next.Y < 0 || next.X < 0 || next.Y >= height || next.X >= width || board[next.Y][next.X] > 0 {
				continue
			}

			if knownCost, ok := gScore[next]; ok && knownCost <= nextCost {
				continue
			} else {
				gScore[next] = nextCost
				reversePath[next] = cur

				updated := false

				for idx, q := range open {
					if q.Value.(image.Point) == next {
						updated = true

						q.Priority = nextCost
						heap.Fix(&open, idx)
						break
					}
				}

				if !updated {
					t := &util.HeapItem{
						Value:    next,
						Priority: nextCost,
					}

					heap.Push(&open, t)
				}
			}
		}
	}

	return false
}

func solution(puzzle *util.Puzzle, width, height int, logger *log.Logger) string {
	left := 1
	right := len(puzzle.Bytes) - 1

	// VV: A binary search to find the exact number of "bytes" that end up blocking the path to the exit
	// If there's a path, add more bytes. If there's not a path, remove a few bytes.
	// Stop when there're no more left options to explore (i.e. [these lead to paths]<this point>[these lead to no paths]])

	for {
		middle := (left + right) / 2
		// VV: can make this a bit faster by generating the board once and then using the integers
		// to decide whether a byte is blocking a tile or not
		board := constructBoard(puzzle, width, height, middle)

		if left == right-1 {
			pos := puzzle.Bytes[middle]
			return fmt.Sprintf("%d,%d", pos.X, pos.Y)
		}
		if !findPath(image.Pt(0, 0), image.Pt(width-1, height-1), board, width, height) {
			right = middle
		} else {
			left = middle
		}

	}
}

func main() {
	logger := SetupLogger()

	logger.Println("Parse input")
	puzzle, err := util.ReadInputFile(os.Args[1])

	// logger.Println("Input was", input)

	if err != nil {
		logger.Fatalln("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, 71, 71, logger)

	logger.Println("Solution is", sol)
}
