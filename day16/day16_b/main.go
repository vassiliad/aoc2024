package main

import (
	"container/heap"
	"fmt"
	"image"
	"os"
	"puzzle/util"
	"slices"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

type BaseState struct {
	pos image.Point
	dir rune
}

type State struct {
	path []image.Point
	BaseState
}

func neighbours(state *State) map[BaseState]int {
	deltas := map[rune]image.Point{
		'>': image.Pt(+1, +0),
		'<': image.Pt(-1, +0),
		'^': image.Pt(+0, -1),
		'v': image.Pt(+0, +1),
	}

	if state.dir == '<' || state.dir == '>' {
		return map[BaseState]int{
			{pos: state.pos, dir: '^'}:                              1000,
			{pos: state.pos, dir: 'v'}:                              1000,
			{pos: state.pos.Add(deltas[state.dir]), dir: state.dir}: 1,
		}
	} else if state.dir == '^' || state.dir == 'v' {
		return map[BaseState]int{
			{pos: state.pos, dir: '<'}:                              1000,
			{pos: state.pos, dir: '>'}:                              1000,
			{pos: state.pos.Add(deltas[state.dir]), dir: state.dir}: 1,
		}
	}

	panic(fmt.Sprintf("unknown direction for state %+v", state))
}

func findAllShortestPaths(start, end image.Point, board [][]rune) int {
	initial := State{BaseState: BaseState{pos: start, dir: '>'}, path: []image.Point{start}}
	gScore := map[BaseState]int{initial.BaseState: 0}

	niceLocation := map[image.Point]int{}
	reversePath := map[image.Point][]image.Point{}

	open := make(util.PriorityQueue, 1)

	open[0] = &util.HeapItem{
		Value:    initial,
		Priority: 0,
	}
	heap.Init(&open)

	bestPathCost := -1
	for open.Len() > 0 {

		heapItem := heap.Pop(&open)

		item := heapItem.(*util.HeapItem)

		cur := item.Value.(State)

		if cur.pos == end && bestPathCost == -1 || bestPathCost == gScore[cur.BaseState] {
			bestPathCost = gScore[cur.BaseState]

			for _, p := range cur.path {
				niceLocation[p] = 1
			}
			continue
		}

		if cur.pos == end {
			continue
		}

		curGScore := gScore[cur.BaseState]

		for baseNext, moveCost := range neighbours(&cur) {
			if board[baseNext.pos.Y][baseNext.pos.X] == '#' {
				continue
			}

			nextCost := curGScore + moveCost

			if knownCost, ok := gScore[baseNext]; ok && knownCost < nextCost {
				continue
			} else {
				gScore[baseNext] = nextCost
				path := make([]image.Point, len(cur.path)+1)
				copy(path, cur.path)
				path = append(path, baseNext.pos)

				next := State{
					BaseState: baseNext, path: path,
				}

				if path, known := reversePath[baseNext.pos]; known {
					if !slices.Contains(path, cur.pos) {
						path = append(path, cur.pos)
						reversePath[baseNext.pos] = path
					}
				} else {
					reversePath[baseNext.pos] = []image.Point{cur.pos}
				}

				t := &util.HeapItem{
					Value:    next,
					Priority: nextCost,
				}

				heap.Push(&open, t)
			}
		}
	}

	for pos, _ := range niceLocation {
		board[pos.Y][pos.X] = 'O'
	}

	for _, row := range board {
		println(string(row))
	}

	return len(niceLocation) - 1
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	return findAllShortestPaths(puzzle.Start, puzzle.End, puzzle.Board)
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
