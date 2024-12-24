package main

import (
	"fmt"
	"log/slog"
	"os"
	"puzzle/util"
	"strconv"
	"strings"
)

/*
VV: While there's at least 1 Z-bit without a value, keep resolving operations which involve operands with known values
*/

func solution(puzzle *util.Puzzle) uint64 {
	ret := uint64(0)
	resolvedBitsOfZ := 0

	for resolvedBitsOfZ < puzzle.ZBits {
		for idx := 0; idx < len(puzzle.Operations); {
			op := puzzle.Operations[idx]

			if _, resolved := puzzle.WireValues[op.Res]; resolved {
				puzzle.Operations[idx] = puzzle.Operations[len(puzzle.Operations)-1]
				puzzle.Operations = puzzle.Operations[:len(puzzle.Operations)-1]

				// VV: Don't increase idx
				continue
			}

			v1, seen1 := puzzle.WireValues[op.A]
			v2, seen2 := puzzle.WireValues[op.B]

			if !seen1 || !seen2 {
				// VV: leave this op in for a future evaluation

				idx++
				continue
			}

			var res int
			switch op.Op {
			case "AND":
				res = v1 & v2
			case "OR":
				res = v1 | v2
			case "XOR":
				res = v1 ^ v2
			default:
				panic(fmt.Sprintf("unknown op in %+v", op))
			}

			puzzle.WireValues[op.Res] = res
			puzzle.Operations[idx] = puzzle.Operations[len(puzzle.Operations)-1]
			puzzle.Operations = puzzle.Operations[:len(puzzle.Operations)-1]

			if zbits := strings.TrimPrefix(op.Res, "z"); zbits != op.Res {
				fmt.Printf("%s = %d %s %d -> %d\n", op.Res, v1, op.Op, v2, res)

				resolvedBitsOfZ++

				bits, err := strconv.Atoi(zbits)

				if err != nil {
					panic(fmt.Errorf("%w for %+v", err, op))
				}
				if res == 1 {
					ret |= 1 << bits
				}
			}
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
