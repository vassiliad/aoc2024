package main

import (
	"container/list"
	"fmt"
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

func tileHasGroup(point Point, membership [][]int, groupId int) bool {
	if point.y < 0 || point.y >= len(membership) || point.x < 0 || point.x >= len(membership[point.y]) {
		return false
	}

	return membership[point.y][point.x] == groupId
}

func handleVerticalEdge(point, pos Point, hull map[Point]int, seenVertical map[Point]int, membership [][]int, groupId int) int {
	last := pos
	for cur := last; hull[cur] != 0; {
		last = cur

		cur = cur.Add(Point{y: -1}, 0)
	}

	topMost := last

	if _, ok := seenVertical[topMost]; ok {
		return 0
	}

	seenVertical[topMost] = 1

	regionLeft := true
	regionRight := true
	vertEdges := 0
	for cur := last; hull[cur] != 0; {
		last = cur

		left := last.Add(Point{x: -1}, 0)
		right := last.Add(Point{x: +1}, 0)

		thisRegionLeft := tileHasGroup(left, membership, groupId)
		if regionLeft != thisRegionLeft && !thisRegionLeft {
			vertEdges++
		}
		regionLeft = thisRegionLeft

		thisRegionRight := tileHasGroup(right, membership, groupId)
		if regionRight != thisRegionRight && !thisRegionRight {
			vertEdges++
		}
		regionRight = thisRegionRight

		cur = cur.Add(Point{y: +1}, 0)
	}

	seenVertical[topMost] = 1

	return vertEdges

}

func fenceTheLotOfThem(point Point, board [][]rune, membership [][]int, groupId int) (int, int) {
	/*
		VV:	Logic:

		1. build the hull of a region (i.e. the points of a region which are neighbouring other regions)
		2. Find all the horizontal and vertical edges of the hull
		3. Walk each edge and take notes whenever a neighbour tile on either side of the edge switches from
			"part of this region" to "not part of this region".  Everytime this happens you basically
			just discovered a new edge
		4. Since we're moving back and forth the hull we have to avoid counting the same edge twice.
			My hack is to always find the leftmost/topmost point of an edge and then if I've already handled
			that "corner" I skip it. The trick is to maintain 2 records of visited nodes, one for vertical
			and another for horizontal edges.
			4.1 Of course, if you only look for vertical edges at the leftmost corner of an edge you'll
				miss out on all the vertical edges on the rightmost part of an edge so I check for vertical edges
				on the rightmost part of a horizontal edge too


		Here's how a walk looks like:

		Let's imagine this Hull where `E` and `i` denotes points on the edge, and innerpoints respectively
			   i i
			EEEEEEE
			  i

		The more I look at the above, the less I like the term "edge". It's actually points sandwitched between 1 or more edges.
		So above we go from "not in region" to "in the region" twice for the top side of the "edge" and
		once for the bottom side of the edge. Therefore we have 3 horizontal edges!

		I hear you say, but that's wrong VV, there're actually 5 edges.

		Well, we need an initial condition to account for the leftmost part of the edge.

		Let's pretend that there're 2 phantom inner points p like so:
			p   i i
			 EEEEEEE
			p  i
		So now, we get the 5 edges we were hoping to see.
	*/
	flower := board[point.y][point.x]

	deltas := []Point{
		{x: 1, y: 0},
		{x: 0, y: 1},
		{x: -1, y: 0},
		{x: 0, y: -1},
	}

	pending := list.New()
	pending.PushFront(point)

	region := map[Point]int{}

	area := 0

	for pending.Front() != nil {
		front := pending.Front()
		pending.Remove(front)

		pos := front.Value.(Point)

		// VV: out of bounds check
		if !(pos.y >= 0 && pos.y < len(board) && pos.x >= 0 && pos.x < len(board[pos.y])) ||
			// VV: Check that this is an adjacent plot with different flowers
			board[pos.y][pos.x] != flower {
			continue
		}

		if membership[pos.y][pos.x] != 0 {
			continue
		}

		area++
		membership[pos.y][pos.x] = groupId
		region[Point{x: pos.x, y: pos.y}] = 1

		for direction, d := range deltas {
			neighbour := pos.Add(d, direction)
			pending.PushFront(neighbour)
		}
	}
	seenHorizontal := map[Point]int{}
	seenVertical := map[Point]int{}

	edges := 0

	hull := map[Point]int{}

	for pos, _ := range region {

		for _, d := range deltas {
			n := pos.Add(d, 0)

			if !tileHasGroup(n, membership, groupId) {
				hull[pos] = 1
				break
			}
		}
	}

	for pos, _ := range hull {
		last := pos

		horizEdges := 0

		for cur := last; hull[cur] != 0; {
			last = cur

			cur = cur.Add(Point{x: -1}, 0)
		}

		leftMost := last

		regionUp := true
		regionDown := true

		for cur := last; hull[cur] != 0; {
			last = cur

			up := last.Add(Point{y: -1}, 0)
			down := last.Add(Point{y: +1}, 0)

			thisRegionUp := tileHasGroup(up, membership, groupId)
			if regionUp != thisRegionUp && !thisRegionUp {
				horizEdges++
			}
			regionUp = thisRegionUp

			thisRegionDown := tileHasGroup(down, membership, groupId)
			if regionDown != thisRegionDown && !thisRegionDown {
				horizEdges++
			}
			regionDown = thisRegionDown

			cur = cur.Add(Point{x: +1}, 0)
		}

		rightMost := last

		if _, ok := seenHorizontal[leftMost]; !ok {
			edges += horizEdges
			seenHorizontal[leftMost] = 1
		}

		edges += handleVerticalEdge(point, leftMost, hull, seenVertical, membership, groupId)
		if leftMost != rightMost {
			edges += handleVerticalEdge(point, rightMost, hull, seenVertical, membership, groupId)
		}

	}

	return area, edges
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

			fmt.Printf("A region of %s plants with price %d * %d = %d\n", string(puzzle.Board[y][x]), area, perimeter, t)
			cost += t
		}

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
