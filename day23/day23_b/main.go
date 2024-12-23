package main

import (
	"fmt"
	"log/slog"
	"os"
	"puzzle/util"
	"slices"
	"strings"
)

func constructGraph(connections []util.Connection) map[string][]string {
	ret := map[string][]string{}

	for _, edge := range connections {
		nodesToA := ret[edge.A]
		if !slices.Contains(nodesToA, edge.B) {
			nodesToA = append(nodesToA, edge.B)
			ret[edge.A] = nodesToA
		}

		nodesToB := ret[edge.B]
		if !slices.Contains(nodesToB, edge.A) {
			nodesToB = append(nodesToB, edge.A)
			ret[edge.B] = nodesToB
		}
	}

	return ret
}

// VV: this is no longer necessary. For details, see wall-of-text in main()
// func findNumberOfCommonEdges(graph map[string][]string, excludeNode, node string) []string {
// 	commonNodes := []string{}

// 	for nodeName, nodes := range graph {
// 		if nodeName == excludeNode {
// 			continue
// 		}

// 		if len(commonNodes) == 0 {
// 			commonNodes = append(commonNodes, nodes...)
// 		} else {
// 			for idx := 0; idx < len(commonNodes); {
// 				if !slices.Contains(nodes, commonNodes[idx]) {
// 					commonNodes[idx] = commonNodes[len(commonNodes)-1]
// 					commonNodes = commonNodes[:len(commonNodes)-1]
// 				} else {
// 					idx++
// 				}
// 			}
// 		}
// 	}

// 	return commonNodes
// }

func solution(puzzle *util.Puzzle) string {
	graph := constructGraph(puzzle.Connections)
	nodeNames := []string{}

	for node := range graph {
		nodeNames = append(nodeNames, node)
	}

	largestGroup := []string{}

	for len(nodeNames) > 0 {
		group := []string{}

		updated := true
		for updated {
			// VV: Keep adding nodes to the group till we cannot add any more
			updated = false

			for idx := 0; idx < len(nodeNames); {
				candidateNode := nodeNames[idx]

				addedToGroup := false

				if len(group) == 0 {
					addedToGroup = true
					group = append(group, candidateNode)
				} else {
					// VV: Check if this node is connected with all the nodes in this group
					// and record if it's not connected with 1 of the nodes
					missing := []string{}

					for _, node := range group {
						knownNodes := graph[candidateNode]
						if !slices.Contains(knownNodes, node) {
							missing = append(missing, node)
						}
					}

					/* VV: This is actually incomplete but it worked for my input.
					The missing part is when a candidate node doesn't form a fully connected graph
					because it doesn't have common nodes with X of the nodes in the current group.
					In this scenario we should check whether adding @candidate and removing X results
					can result in a fullyConnected graph that can grow larger.
					I was going to calculate the common nodes between the candidate node and the group
					and then compare those against the common nodes between the @missing nodes (i.e. X)
					if the common nodes of @candidate were more than those of X I'd remove X and add
					@candidate.

					I'd start with implementing this for len(X) == 1 and move from there till my code
					provided a valid solution for my input. However, I figured it wouldn't hurt to try
					the answer I got before implementing this and I was surprised to see that it was
					good enough as is.
					*/
					if len(missing) == 0 {
						addedToGroup = true
						group = append(group, candidateNode)
					}
				}

				updated = updated || addedToGroup
				// VV: when adding a node to the group just replace it with the last known node
				if addedToGroup {
					nodeNames[idx] = nodeNames[len(nodeNames)-1]
					nodeNames = nodeNames[:len(nodeNames)-1]
				} else {
					idx++
				}
			}
		}

		if len(group) > len(largestGroup) {
			largestGroup = group
		}
	}

	slices.Sort(largestGroup)
	return strings.Join(largestGroup, ",")
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
