package main

import (
	"container/heap"
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

func constructBoard(puzzle *util.Puzzle, width, height int) [][]int {
	board := [][]int{}

	for range height {
		board = append(board, make([]int, width))
	}

	for i, b := range puzzle.Bytes {
		board[b.Y][b.X] = i + 1
	}

	return board
}

func findPath(start, end image.Point, board [][]int, width, height int) int {
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
			board[start.Y][start.X] = MAGIC_NUMBER
			for p := cur; p != start; p = reversePath[p] {
				board[p.Y][p.X] = MAGIC_NUMBER
			}

			board[start.Y][start.X] = MAGIC_NUMBER

			prettyPrint(board)

			return gScore[cur]
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

	panic("oh no")
}

func solution(puzzle *util.Puzzle, width, height, cutoff int, logger *log.Logger) int {
	if len(puzzle.Bytes) > cutoff {
		puzzle.Bytes = puzzle.Bytes[0:cutoff]
	}

	board := constructBoard(puzzle, width, height)

	// prettyPrint(board)

	return findPath(image.Pt(0, 0), image.Pt(width-1, height-1), board, width, height)
}

func main() {
	logger := SetupLogger()

	logger.Println("Parse input")
	puzzle, err := util.ReadInputFile(os.Args[1])

	// logger.Println("Input was", input)

	if err != nil {
		logger.Fatalln("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, 71, 71, 1024, logger)

	logger.Println("Solution is", sol)
}
