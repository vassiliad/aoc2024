package main

import (
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

func isInOrder(rules []image.Point, update map[int]int) bool {

	for _, r := range rules {
		left, okl := update[r.X]
		right, okr := update[r.Y]

		if (okl && okr) && left > right {
			return false
		}
	}

	return true
}

func reorder(dependencies map[int][]int, update []int) {
	/*VV: This is a weird insertion sort of sorts (pun intended).

	Starting from a dependencies map (see calculateDependencies()) it finds out the right position
	for a page.

	This solution assumes that there are rules for all pages.

	First, we filter the dependencies so that they only references to pages that the current update contains.

	Then, we just look at each page P and make sure it's positioned at an index that leaves enough space
	after P for the pages that should come after P.
	*/
	filteredDependencies := map[int][]int{}

	for k, after := range dependencies {
		if !slices.Contains(update, k) {
			continue
		}

		new_v := []int{}
		for _, p := range after {
			if slices.Contains(update, p) {
				new_v = append(new_v, p)
			}
		}

		filteredDependencies[k] = new_v
	}

	new_update := make([]int, len(update))

	// VV: If a page didn't have a rule for it, we could pop them at the back of new_update
	for _, p := range update {
		idx := len(update) - len(filteredDependencies[p]) - 1
		new_update[idx] = p
	}

	copy(update, new_update)
}

func calculateDependencies(rules []image.Point) map[int][]int {
	/*VV: Create a map whose keys are pages and values are list of pages that should go after the key page*/
	dependencies := map[int][]int{}

	for _, r := range rules {
		if after, ok := dependencies[r.X]; ok {
			dependencies[r.X] = append(after, r.Y)
		} else {
			dependencies[r.X] = []int{r.Y}
		}
	}

	return dependencies
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {

	solution := 0

	dependencies := calculateDependencies(puzzle.Rules)

	for _, update := range puzzle.Updates {
		updateOrder := map[int]int{}
		for i, p := range update {
			updateOrder[p] = i
		}

		if !isInOrder(puzzle.Rules, updateOrder) {
			reorder(dependencies, update)
			solution += update[len(update)/2]
		}
	}

	return solution
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
