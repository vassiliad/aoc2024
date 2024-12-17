package main

import (
	"container/heap"
	"fmt"
	"os"
	"puzzle/util"
	"strings"

	"log"
)

func SetupLogger() *log.Logger {
	logger := log.New(os.Stdout, "", 0)

	logger.SetFlags(log.LstdFlags)
	return logger
}

const (
	ADV = 0
	BXL = 1
	BST = 2
	JNZ = 3
	BXC = 4
	OUT = 5
	BDV = 6
	CDV = 7
)

var OP_LABELS = map[int64]string{
	0: "ADV", 1: "BXL", 2: "BST", 3: "JNZ", 4: "BXC", 5: "OUT", 6: "BDV", 7: "CDV",
}

func division(numerator, comboValue int64) int64 {
	// denominator := int64(math.Pow(2, float64(comboValue)))
	// ret := numerator / denominator
	// fmt.Printf("%d / %d = %d", numerator, denominator, ret)
	// return ret

	return numerator >> comboValue
}

func prettyPrint(puzzle *util.Puzzle) {
	for i, value := range puzzle.Registers {
		fmt.Printf("Register %s: %d\n", string(i+'A'), value)
	}

	fmt.Println()
	for i := 0; i < len(puzzle.Program); i += 2 {
		opcode := puzzle.Program[i]
		operand := puzzle.Program[i+1]

		fmt.Printf("%2d: %s %d\n", i, OP_LABELS[opcode], operand)
	}
}

func simulate(puzzle *util.Puzzle) {
	opcode := puzzle.Program[puzzle.PC]
	operand := puzzle.Program[puzzle.PC+1]

	// VV: These 2 instructions use literal operands, everything else a combo operand
	if opcode != BXL && opcode != JNZ {
		if operand > 3 && operand < 7 {
			operand = puzzle.Registers[operand-4]
		} else if operand == 7 {
			panic(fmt.Sprintf("opcode %d with operand %d at %d", opcode, operand, puzzle.PC))
		}

	}

	puzzle.PC += 2

	switch opcode {
	case ADV:
		puzzle.Registers[0] = division(puzzle.Registers[0], operand)
	case BXL:
		ret := puzzle.Registers[1] ^ operand
		puzzle.Registers[1] = ret
	case BST:
		ret := operand & 7
		puzzle.Registers[1] = ret
	case JNZ:

		if puzzle.Registers[0] != 0 {
			puzzle.PC = operand
		}
	case BXC:
		ret := puzzle.Registers[1] ^ puzzle.Registers[2]
		puzzle.Registers[1] = ret
	case OUT:
		puzzle.Output = append(puzzle.Output, operand&7)
	case BDV:
		puzzle.Registers[1] = division(puzzle.Registers[0], operand)
	case CDV:
		puzzle.Registers[2] = division(puzzle.Registers[0], operand)
	default:
		panic(fmt.Sprintf("Unknown opcode %d and input %d at %d", opcode, operand, puzzle.PC-2))
	}
}

type State struct {
	alpha int64 // VV: this is the WIP value of Alpha - we're building it startin from its leftmost bits
	found int   // VV: keeps track of how many output's we've found so far
}

/*
VV: the puzzle program looks at the rightmost 6 bits of register A and based on those, it:

1. prints a number
2. decides whether it should loop one more time or not

Multiple 6-bit combinations print the same number.

The solution is to try 6-bit combinations to discover the rightmost 6-bits of A.

The 6-bit combination is wrong when:
 1. the program does not produce 2 outputs
 2. the 2 outputs it produces are not the ones we expect to find (we start from the rightmost OpCode/Operands
    and move our way towards the front - these correspond to the left most bits of Alpha)
 3. the 2 outputs it produces are the ones we expect BUT the program exits early

We're looking for the minimum value of A, so I'm using a priority queue to keep track of the incomplete Alpha values
we're trying out. This works because we're building Alpha from the leftmost bits.

Finally, avoid retrying the same bits over over by using a memoization cache
*/
func reverseEngineerAlpha(puzzle *util.Puzzle) int64 {
	pending := make(util.PriorityQueue, 1)

	pending[0] = &util.HeapItem{
		Value:    State{},
		Priority: 0,
	}
	heap.Init(&pending)
	seen := map[State]int{}

	for pending.Len() > 0 {

		heapItem := heap.Pop(&pending)
		item := heapItem.(*util.HeapItem)
		cur := item.Value.(State)

		for attempt := int64(0); attempt < 0b1_000_000; attempt++ {
			// VV: We discover the last 6 bits every time, the higher bits are only there
			// for the puzzle code to decide whether it'll produce more outputs in the future or not
			alpha := cur.alpha<<6 | attempt

			puzzle.PC = 0
			puzzle.Output = puzzle.Output[:0]

			// VV: We don't care about the other registers, the puzzle code generates them using the value of A
			puzzle.Registers[0] = alpha

			for len(puzzle.Output) != 2 && puzzle.PC < int64(len(puzzle.Program)) {
				simulate(puzzle)
			}

			// VV: Puzzle code exited early and didn't produce the 2 outputs
			if len(puzzle.Output) != 2 {
				continue
			}

			if puzzle.Output[1] == puzzle.Program[len(puzzle.Program)-cur.found-1] && puzzle.Output[0] == puzzle.Program[len(puzzle.Program)-cur.found-1-1] {
				next := cur
				next.found += 2

				if next.found == len(puzzle.Program) {
					// VV: all done
					return alpha
				} else {
					// VV: The 2 outputs are those we expected. Need to check the rest
					next.alpha = alpha

					if _, ok := seen[next]; ok {
						continue
					}
					seen[next] = 1

					t := &util.HeapItem{
						Value:    next,
						Priority: int(next.alpha),
					}

					heap.Push(&pending, t)
				}
			}

		}
	}

	panic("good luck")
}

func solution(puzzle *util.Puzzle, logger *log.Logger) int64 {
	prettyPrint(puzzle)

	answer := int64(0)
	answer = reverseEngineerAlpha(puzzle)

	puzzle.PC = 0
	puzzle.Registers = [3]int64{answer, 0, 0}
	puzzle.Output = puzzle.Output[:0]
	for puzzle.PC < int64(len(puzzle.Program)) {
		simulate(puzzle)
	}

	ret := []string{}
	for _, out := range puzzle.Program {
		ret = append(ret, fmt.Sprint(out))
	}

	correct := strings.Join(ret, ",")
	println("Correct: ", correct)

	ret = ret[:0]
	for _, out := range puzzle.Output {
		ret = append(ret, fmt.Sprint(out))
	}

	computed := strings.Join(ret, ",")
	println("Computed:", computed)

	if correct != computed {
		panic("good luck")
	}

	return answer
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
