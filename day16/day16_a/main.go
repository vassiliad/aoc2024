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
	dir rune
}

func neighbours(state *State) map[State]int {
	deltas := map[rune]image.Point{
		'>': image.Pt(+1, +0),
		'<': image.Pt(-1, +0),
		'^': image.Pt(+0, -1),
		'v': image.Pt(+0, +1),
	}

	if state.dir == '<' || state.dir == '>' {
		return map[State]int{
			{pos: state.pos, dir: '^'}:                              1000,
			{pos: state.pos, dir: 'v'}:                              1000,
			{pos: state.pos.Add(deltas[state.dir]), dir: state.dir}: 1,
		}
	} else if state.dir == '^' || state.dir == 'v' {
		return map[State]int{
			{pos: state.pos, dir: '<'}:                              1000,
			{pos: state.pos, dir: '>'}:                              1000,
			{pos: state.pos.Add(deltas[state.dir]), dir: state.dir}: 1,
		}
	}

	panic(fmt.Sprintf("unknown direction for state %+v", state))
}

func findPath(start, end image.Point, board [][]rune) int {
	initial := State{pos: start, dir: '>'}
	gScore := map[State]int{initial: 0}

	reversePath := map[State]State{}
	open := make(util.PriorityQueue, 1)

	open[0] = &util.HeapItem{
		Value:    initial,
		Priority: 0,
	}
	heap.Init(&open)

	for open.Len() > 0 {

		heapItem := heap.Pop(&open)

		item := heapItem.(*util.HeapItem)

		cur := item.Value.(State)

		if cur.pos == end {
			for p := cur; p.pos != start; p = reversePath[p] {
				board[p.pos.Y][p.pos.X] = p.dir
			}

			for _, row := range board {
				println(string(row))
			}

			return gScore[cur]
		}

		curGScore := gScore[cur]

		for next, moveCost := range neighbours(&cur) {
			if board[next.pos.Y][next.pos.X] == '#' {
				continue
			}

			nextCost := curGScore + moveCost

			if knownCost, ok := gScore[next]; ok && knownCost <= nextCost {
				continue
			} else {
				gScore[next] = nextCost
				reversePath[next] = cur

				updated := false

				for idx, q := range open {
					if q.Value.(State) == next {
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

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	return findPath(puzzle.Start, puzzle.End, puzzle.Board)
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
