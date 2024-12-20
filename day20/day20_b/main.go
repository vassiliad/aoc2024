package main

import (
	"container/heap"
	"image"
	"os"
	"puzzle/util"

	"iter"
	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

type State struct {
	pos        image.Point
	cheatsLeft int
	cheatStart image.Point
}

var deltas = []image.Point{
	image.Pt(+1, +0),
	image.Pt(-1, +0),
	image.Pt(+0, -1),
	image.Pt(+0, +1),
}

func findAllDistances(start image.Point, board [][]rune) map[image.Point]int {
	height := len(board)
	width := len(board[0])

	initial := start
	gScore := map[image.Point]int{initial: 0}

	open := make(util.PriorityQueue, 1)

	open[0] = &util.HeapItem{
		Value:    initial,
		Priority: 0,
	}
	heap.Init(&open)

	for open.Len() > 0 {
		heapItem := heap.Pop(&open)

		item := heapItem.(*util.HeapItem)

		cur := item.Value.(image.Point)

		nextCost := gScore[cur] + 1

		for _, d := range deltas {
			next := cur.Add(d)

			if next.Y < 0 || next.X < 0 || next.Y >= height || next.X >= width || board[next.Y][next.X] == '#' {
				continue
			}

			if knownCost, ok := gScore[next]; ok && knownCost <= nextCost {
				continue
			} else {
				gScore[next] = nextCost
				open.Upsert(next, nextCost)
			}
		}
	}

	return gScore
}

/*
Switch off collision and find out which walkable tiles you can find by phasing through walls
*/
func cheat(start image.Point, board [][]rune, maxCheats int) iter.Seq2[image.Point, int] {
	height := len(board)
	width := len(board[0])

	return func(yield func(image.Point, int) bool) {
		for dy := -maxCheats; dy <= +maxCheats; dy++ {
			for dx := -maxCheats; dx <= +maxCheats; dx++ {

				absDx := dx
				if absDx < 0 {
					absDx = -dx
				}

				absDy := dy
				if absDy < 0 {
					absDy = -dy
				}
				distance := absDx + absDy
				if distance > maxCheats {
					continue
				}

				pos := start
				pos.X += dx
				pos.Y += dy

				if pos.X < 0 || pos.Y < 0 || pos.X >= width || pos.Y >= height || board[pos.Y][pos.X] == '#' {
					continue
				}

				if !yield(pos, distance) {
					return
				}
			}
		}

	}
}

func groupPaths(puzzle *util.Puzzle) map[int]int {
	const maxCheats = 20

	// VV: Find the distances of all points to the end
	distanceToEnd := findAllDistances(puzzle.End, puzzle.Board)

	// VV: Keys are picoseconds saved and values are number of paths
	groups := map[int]int{}

	shortestPath := distanceToEnd[puzzle.Start]

	for y, row := range puzzle.Board {
		for x, c := range row {
			if c == '#' || c == 'E' {
				continue
			}

			pos := image.Pt(x, y)
			distFromStart := shortestPath - distanceToEnd[pos]

			// VV: For each point on the shortest path, switch off collision for maxCheats picoseconds
			cheatTiles := cheat(pos, puzzle.Board, maxCheats)

			for tile, cheatDistance := range cheatTiles {
				if tile == pos {
					continue
				}
				distToEnd, ok := distanceToEnd[tile]

				if !ok {
					// VV: The tile is not connected to the End, cheating didn't help at all
					continue
				}

				// VV: It takes distFromStart steps to reach a walkable tile starting from S, then it takes
				// 1 step to go into a wall and cheatDistance-1 to walk around inside the walls and back onto a walkable tile
				// then it takes distToEnd steps to walk to the end
				distance := distFromStart + cheatDistance + distToEnd

				saved := shortestPath - distance

				if saved >= 50 {
					groups[saved]++
				}

			}
		}
	}

	return groups
}

func solution(puzzle *util.Puzzle, saveAtLeast int, logger *log.Logger) int {
	groups := groupPaths(puzzle)

	ret := 0

	for saved, numPaths := range groups {
		if saved >= saveAtLeast {
			ret += numPaths
		}
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

	sol := solution(puzzle, 100, logger)

	logger.Println("Solution is", sol)
}
