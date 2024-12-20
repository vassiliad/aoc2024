package main

import (
	"puzzle/util"
	"slices"
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
		50: 32, 52: 31, 54: 29, 56: 39, 58: 25, 60: 23, 62: 20, 64: 19, 66: 12, 68: 14,
		70: 12, 72: 22, 74: 4, 76: 3,
	}

	wrong := 0

	allExpected := []int{}

	for e, _ := range expected {
		allExpected = append(allExpected, e)
	}

	slices.Sort(allExpected)

	for _, saved := range allExpected {

		numPaths := groups[saved]
		correctSaved, ok := expected[saved]

		if !ok {
			t.Fatal("Got an unexpected entry for saving", saved, "picoseconds")
		}

		if correctSaved != numPaths {
			t.Log("Expected numPaths for saved", saved, "to be", correctSaved, "but it was", numPaths)
			wrong++
		}
	}

	if wrong > 0 {
		t.Fatal(wrong, "mistakes")
	}
}
