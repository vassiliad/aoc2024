package main

import (
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `
###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############
`

	// println(small)

	puzzle, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", puzzle)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	groups := groupPaths(puzzle)
	expected := map[int]int{
		2: 14, 4: 14, 6: 2, 8: 4, 10: 2, 12: 3, 20: 1, 36: 1, 38: 1, 40: 1, 64: 1,
	}

	for saved, numPaths := range groups {
		if saved <= 0 {
			continue
		}

		correctSaved, ok := expected[saved]

		if !ok {
			t.Fatal("Got an unexpected entry for saving", saved, "picoseconds")
		}

		if correctSaved != numPaths {
			t.Fatal("Expected numPaths for saved", saved, "to be", correctSaved, "but it was", numPaths)
		}
	}
}
