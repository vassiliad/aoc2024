package main

import (
	"fmt"
	"log/slog"
	"os"
	"puzzle/util"
	"slices"
	"strings"
)

func solution(puzzle *util.Puzzle) int {
	ret := 0

	seen := map[string]int{}
	/*
		VV: For every 3 edges:

		edge0: A -> B (also B -> A)
		edge1: B -> C (also C -> B)
		edge2: C -> A (also C -> A)

		Find the "3rd" node in edge1 (i.e. C). Then find an edge2 which connects the 3rd node to
		either A or B. If any of the 3 nodes starts with a T then the 3 nodes form one of the groups
		we're looking for.
	*/
	for i := 0; i < len(puzzle.Connections)-2; i++ {
		edge0 := puzzle.Connections[i]

		for j := i + 1; j < len(puzzle.Connections)-1; j++ {
			edge1 := puzzle.Connections[j]

			// VV: this will contain the name of the 3rd node in the group
			thirdNode := ""

			if edge0.A == edge1.A {
				thirdNode = edge1.B
			} else if edge0.A == edge1.B {
				thirdNode = edge1.A
			} else if edge0.B == edge1.A {
				thirdNode = edge1.B
			} else if edge0.B == edge1.B {
				thirdNode = edge1.A
			} else {
				// VV: This second edge does not connect to any of the nodes in edge 0
				continue
			}

			anyT := strings.HasPrefix(edge0.A, "t") || strings.HasPrefix(edge0.B, "t") || strings.HasPrefix(thirdNode, "t")

			// VV: None of the 3 nodes start with a "t", skip looking for the edge with the 3rd node
			if !anyT {
				continue
			}

			for k := j + 1; k < len(puzzle.Connections); k++ {
				edge2 := puzzle.Connections[k]

				oneOfFirstTwoNodes := ""
				if edge2.A == thirdNode {
					oneOfFirstTwoNodes = edge2.B
				} else if edge2.B == thirdNode {
					oneOfFirstTwoNodes = edge2.A
				}

				nodes := []string{edge0.A, edge0.B, thirdNode}

				if oneOfFirstTwoNodes == edge0.B || oneOfFirstTwoNodes == edge0.A {
				} else {
					continue
				}

				slices.Sort(nodes)
				key := strings.Join(nodes, ",")

				if seen[key] == 0 {
					ret++
					seen[key] = 1

					break
				}
			}
		}
	}

	return ret
}

func main() {
	slog.Info("Parse input")
	puzzle, err := util.ReadInputFile(os.Args[1])

	// slog.Info("Input was", "input", puzzle)

	if err != nil {
		slog.Error("Ran into problems while reading input. Problem", "error", err)
	}

	sol := solution(puzzle)

	slog.Info(fmt.Sprint("Solution is ", sol))
}
