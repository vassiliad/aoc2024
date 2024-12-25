package main

import (
	"fmt"
	"log/slog"
	"os"
	"puzzle/util"
)

func solution(puzzle *util.Puzzle) int {
	ret := 0

	for _, lock := range puzzle.Locks {
		for _, key := range puzzle.Keys {
			for i := range lock {
				if lock[i]+key[i] > len(lock) {
					goto tryNextKey
				}
			}

			ret++

		tryNextKey:
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
