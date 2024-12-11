package main

import (
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

func setStone(stones []int, idx, stone int) []int {
	if idx < len(stones) {
		stones[idx] = stone
	} else if idx == len(stones) {
		stones = append(stones, stone)
	} else {
		panic(fmt.Sprintf("Adding stone %d with index %d to array %+v with len %d", stone, idx, stones, len(stones)))
	}

	return stones
}

func solution(puzzle *util.Puzzle, blinkTimes int, logger *log.Logger) int {
	newStones := []int{}

	for range blinkTimes {
		newIdx := 0

		for _, stone := range puzzle.Stones {
			if stone == 0 {
				newStones = setStone(newStones, newIdx, 1)
			} else if digits := util.NumberOfDigits(stone); digits%2 == 0 {
				left, right := util.SplitNumber(stone, digits/2)
				newStones = setStone(newStones, newIdx, left)
				newIdx++
				newStones = setStone(newStones, newIdx, right)
			} else {
				newStones = setStone(newStones, newIdx, stone*2024)
			}

			newIdx++
		}

		puzzle.Stones, newStones = newStones, puzzle.Stones
	}

	return len(puzzle.Stones)
}

func main() {
	logger := SetupLogger()

	logger.Println("Parse input")
	puzzle, err := util.ReadInputFile(os.Args[1])

	// logger.Println("Input was", input)

	if err != nil {
		logger.Fatalln("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, 25, logger)

	logger.Println("Solution is", sol)
}
