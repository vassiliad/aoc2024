package main

import (
	"os"
	"puzzle/util"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

type Block struct {
	size, idx int
}

func (b *Block) print() {
	for range b.size {
		if b.idx > -1 {
			print(b.idx)

		} else {
			print(".")
		}
	}
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int {
	blocks := []Block{{size: puzzle.BlockSizes[0], idx: 0}}

	for idx, free := range puzzle.Free {
		free_block := Block{size: free, idx: -1}
		file_block := Block{size: puzzle.BlockSizes[idx+1], idx: idx + 1}

		blocks = append(blocks, free_block, file_block)
	}

	// VV: This is basically best fit. There's also no joining of consecutive empty blocks.
	// If there was, in the example we'd see 8888 move between 333 and 5555 here:
	//    0099.111777244.333....5555.6666.....8888..

	// VV: This is actually a lie, it also keeps track of blocks skipped (i.e. those which are empty)
	blocksMoved := 1

	// VV: We need only try to move a block once!
	for idx := len(blocks) - 1; idx > 0; {

	next_loop:
		idx := len(blocks) - blocksMoved
		if idx == 0 {
			break
		}
		blocksMoved++

		if blocks[idx].idx == -1 {
			goto next_loop
		}

		for moveToIdx := 1; moveToIdx < idx; moveToIdx++ {
			if blocks[moveToIdx].idx != -1 {
				continue
			}

			if blocks[moveToIdx].size == blocks[idx].size {
				blocks[moveToIdx].idx = blocks[idx].idx
				blocks[idx].idx = -1

				goto next_loop
			} else if blocks[moveToIdx].size > blocks[idx].size {
				remainder := blocks[moveToIdx].size - blocks[idx].size

				blocks[moveToIdx].idx = blocks[idx].idx
				blocks[moveToIdx].size = blocks[idx].size

				blocks[idx].idx = -1

				newBlocks := []Block{}

				newBlocks = append(newBlocks, blocks[:moveToIdx+1]...)
				newBlocks = append(newBlocks, Block{size: remainder, idx: -1})
				newBlocks = append(newBlocks, blocks[moveToIdx+1:]...)

				blocks = newBlocks
				goto next_loop
			}
		}
	}

	sol := 0
	idx := 0
	for _, b := range blocks {
		if b.idx > 0 {
			for i := 0; i < b.size; i++ {
				sol += (idx + i) * b.idx
			}
		}
		idx += b.size
	}

	// for _, b := range blocks {
	// 	b.print()
	// }

	// println()

	return sol
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
