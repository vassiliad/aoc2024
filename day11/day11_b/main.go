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

type Memo struct {
	stone, blinkTimes int
}

func simulate(stone, blinkTimes int, book map[Memo]int) int {
	if blinkTimes == 0 {
		return 1
	}

	if value, seen := book[Memo{stone: stone, blinkTimes: blinkTimes}]; seen {
		return value
	}

	blinkTimes--

	var value int

	if stone == 0 {
		value = simulate(1, blinkTimes, book)
	} else if digits := util.NumberOfDigits(stone); digits%2 == 0 {
		left, right := util.SplitNumber(stone, digits/2)
		value = simulate(left, blinkTimes, book) + simulate(right, blinkTimes, book)
	} else {
		value = simulate(stone*2024, blinkTimes, book)
	}

	book[Memo{stone: stone, blinkTimes: blinkTimes + 1}] = value

	return value
}

func solution(puzzle *util.Puzzle, blinkTimes int, logger *log.Logger) int {
	totalStones := 0

	book := map[Memo]int{}

	for idx, stone := range puzzle.Stones {
		thisStone := simulate(stone, blinkTimes, book)

		logger.Printf("Stone[%d]=%d had %d offsprings\n", idx, stone, thisStone)
		totalStones += thisStone
	}

	return totalStones
}

func main() {
	logger := SetupLogger()

	logger.Println("Parse input")
	puzzle, err := util.ReadInputFile(os.Args[1])

	// logger.Println("Input was", input)

	if err != nil {
		logger.Fatalln("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, 75, logger)

	logger.Println("Solution is", sol)
}
