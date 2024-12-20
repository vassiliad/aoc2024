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

func findBestPath(start, end image.Point, board [][]rune, maxCheats int, recordPath bool, cutoff int) (int, []image.Point) {
	height := len(board)
	width := len(board[0])

	initial := State{pos: start, cheatsLeft: maxCheats, cheatStart: image.Pt(-1, -1)}
	gScore := map[State]int{initial: 0}

	open := make(util.PriorityQueue, 1)

	open[0] = &util.HeapItem{
		Value:    initial,
		Priority: 0,
	}
	heap.Init(&open)

	reversePath := map[State]State{}

	for open.Len() > 0 {
		heapItem := heap.Pop(&open)

		item := heapItem.(*util.HeapItem)

		cur := item.Value.(State)
		cost := gScore[cur]

		if cost > cutoff {
			continue
		}

		if cur.pos == end {
			path := make([]image.Point, cost+1)
			path[0] = start

			if recordPath {
				idx := cost
				for ; cur != initial; cur = reversePath[cur] {
					path[idx] = cur.pos
					idx--
				}
			}

			return cost, path
		}

		nextCost := cost + 1

		for _, d := range deltas {
			pos := cur.pos.Add(d)

			if pos.Y < 0 || pos.X < 0 || pos.Y >= height || pos.X >= width {
				continue
			}

			// VV: We can cheat up to N times but we can only cheat if this is the 1st time we cheat OR we just cheated to get here
			next := State{
				cheatsLeft: cur.cheatsLeft, pos: pos, cheatStart: cur.cheatStart,
			}

			if board[pos.Y][pos.X] == '#' {
				if next.cheatsLeft <= 0 {
					continue
				}

				// VV: This is the 1st time we cheat
				if next.cheatStart.X == -1 {
					next.cheatStart = pos
				}

				next.cheatsLeft--
			} else if next.cheatStart.X != -1 {
				// VV: We got here by cheating and this tile is walkable, we can no longer cheat
				next.cheatsLeft = 0
			}

			if knownCost, ok := gScore[next]; ok && knownCost <= nextCost {
				continue
			} else {
				if recordPath {
					reversePath[next] = cur
				}

				gScore[next] = nextCost
				open.Upsert(next, nextCost)
				// open.PushValue(next, nextCost)
			}
		}
	}

	// VV: No path that's shorter than @cutoff
	return -1, []image.Point{}
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

func prettyPrint(puzzle *util.Puzzle, pos, wall, tile image.Point) {
	for y, row := range puzzle.Board {
		for x, c := range row {
			if x == pos.X && y == pos.Y {
				print("@")
			} else if x == wall.X && y == wall.Y {
				print("1")
			} else if x == tile.X && y == tile.Y {
				print("2")
			} else {
				print(string(c))
			}
		}

		println()
	}
}

func groupPaths(puzzle *util.Puzzle) map[int]int {
	// VV: First find the path from start to end without cheating
	shortestPath, tiles := findBestPath(puzzle.Start, puzzle.End, puzzle.Board, 0, true, int(^uint(0)>>1))

	// VV: Find the distances of all points to the end
	distanceToEnd := findAllDistances(puzzle.End, puzzle.Board)

	// VV: Keys are picoseconds saved and values are number of paths
	groups := map[int]int{}

	width := len(puzzle.Board[0])
	height := len(puzzle.Board)

	for idx, pos := range tiles {
		// VV: For each point on the shortest path, find an adjacent wall
		// and then from that wall find an adjacent tile
		for _, wallD := range deltas {
			wall := pos.Add(wallD)
			if puzzle.Board[wall.Y][wall.X] == '#' {
				for _, d := range deltas {
					tile := wall.Add(d)
					if tile == pos {
						continue
					}

					if tile.X < 0 || tile.Y < 0 || tile.X >= width || tile.Y >= height || puzzle.Board[tile.Y][tile.X] == '#' {
						continue
					}

					d, ok := distanceToEnd[tile]

					if !ok {
						// VV: The tile is not connected to the End, cheating didn't help at all
						continue
					}

					// VV: It takes idx steps to reach the tile on the shortestPath
					// 1 step to go into the wall and another to exit the wall onto another tile
					// then it takes d steps to walk to the end
					distance := d + 1 + 1 + idx

					saved := shortestPath - distance

					if saved > 0 {
						groups[saved] = groups[saved] + 1
					}
				}
			}
		}
	}

	return groups
}

func solution(puzzle *util.Puzzle, saveAtLeast int, logger *log.Logger) int {
	const cutoff = 100
	groups := groupPaths(puzzle)

	ret := 0

	for saved, numPaths := range groups {
		if saved >= cutoff {
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
