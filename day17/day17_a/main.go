package main

import (
	"fmt"
	"math"
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
	denominator := int64(math.Pow(2, float64(comboValue)))
	ret := numerator / denominator
	fmt.Printf("%d / %d = %d", numerator, denominator, ret)
	return ret
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

	fmt.Printf("%2d: %s %d (%d) ", puzzle.PC, OP_LABELS[opcode], operand, puzzle.Program[puzzle.PC+1])

	puzzle.PC += 2

	switch opcode {
	case ADV:
		fmt.Printf("-> A:")
		puzzle.Registers[0] = division(puzzle.Registers[0], operand)
	case BXL:
		ret := puzzle.Registers[1] ^ operand
		fmt.Printf("-> B: %d ^ %d = %d", puzzle.Registers[1], operand, ret)
		puzzle.Registers[1] = ret
	case BST:
		ret := operand & 7
		fmt.Printf("-> B: %d %% %d = %d", puzzle.Registers[1], operand, ret)
		puzzle.Registers[1] = ret
	case JNZ:
		fmt.Printf("A=%d ", puzzle.Registers[0])
		if puzzle.Registers[0] != 0 {
			fmt.Printf("Jump to %d", operand)
			puzzle.PC = operand
		}
	case BXC:
		ret := puzzle.Registers[1] ^ puzzle.Registers[2]
		fmt.Printf("-> B: %d ^ %d = %d", puzzle.Registers[1], puzzle.Registers[2], ret)
		puzzle.Registers[1] = ret
	case OUT:
		fmt.Printf("-> Output %d (%d)", operand&7, operand)
		puzzle.Output = append(puzzle.Output, operand&7)
	case BDV:
		fmt.Printf("-> B:")
		puzzle.Registers[1] = division(puzzle.Registers[0], operand)
	case CDV:
		fmt.Printf("-> C:")
		puzzle.Registers[2] = division(puzzle.Registers[0], operand)
	default:
		panic(fmt.Sprintf("Unknown opcode %d and input %d at %d", opcode, operand, puzzle.PC-2))
	}

	fmt.Println()
}

func solution(puzzle *util.Puzzle, logger *log.Logger) string {
	for puzzle.PC < int64(len(puzzle.Program)) {
		simulate(puzzle)
	}

	ret := []string{}
	for _, out := range puzzle.Output {
		ret = append(ret, fmt.Sprint(out))
	}

	return strings.Join(ret, ",")
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
