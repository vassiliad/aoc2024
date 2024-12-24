package main

import (
	"fmt"
	"log/slog"
	"os"
	"puzzle/util"
	"slices"
	"strconv"
	"strings"
)

func plot(operations []util.Operation) {
	nameToOp := map[string]string{}

	for _, op := range operations {
		nameToOp[op.Res] = fmt.Sprintf("%s = %s %s %s", op.Res, op.A, op.Op, op.B)
	}

	fmt.Println("digraph G {")
	for _, op := range operations {

		fmt.Printf("\"%s\" -> \"%s\"\n", nameToOp[op.A], nameToOp[op.Res])
		fmt.Printf("\"%s\" -> \"%s\"\n", nameToOp[op.B], nameToOp[op.Res])
	}
	fmt.Println("}")
}

/*
VV:

	Need to troubleshoot a broken adder circuit

	The Z and Carry in bits are:

	Z_i = (X_i XOR Y_i) XOR C_i-1
	C_i = (X_i AND Y_i) OR ((X_i XOR Y_i) AND C_i-1)

	The 1st bit is Z_0 = X_0 XOR Y_0
	The nth bit is Z_N = C_(N-1)

	I stopped implementing checks as soon as I had enough to identify the 8 swapped gates in my input.
	Going any further felt like a chore to me.
*/
func solution(puzzle *util.Puzzle) string {

	// plot(puzzle.Operations)

	// VV: Make a record of useful information and potentially detect problems
	andGatesOfInputs := map[int]string{}
	xorGatesOfInputs := map[int]string{}

	producedBy := map[string]util.Operation{}
	consumedBy := map[string][]util.Operation{}
	foundAndGatesOfInputs := 0
	foundXorGatesOfInputs := 0

	swappedGates := map[string]int{}

	for _, op := range puzzle.Operations {
		producedBy[op.Res] = op

		for _, name := range [2]string{op.A, op.B} {
			if _, ok := consumedBy[name]; !ok {
				consumedBy[name] = []util.Operation{op}
			} else {
				consumedBy[name] = append(consumedBy[name], op)
			}
		}

		// VV: We expect Xi and Yi to be used in exactly 2 gates a XOR and an AND
		if (strings.HasPrefix(op.A, "x") && strings.HasPrefix(op.B, "y")) || (strings.HasPrefix(op.B, "x") && strings.HasPrefix(op.A, "y")) {
			bitsA, _ := strconv.Atoi(op.A[1:])
			bitsB, _ := strconv.Atoi(op.B[1:])
			if bitsA != bitsB {
				panic(fmt.Sprintf("Using inputs with different bit-index %+v\n", op))
			}

			if op.Op == "OR" {
				panic(fmt.Sprintf("Using inputs with an OR %+v\n", op))
			}

			// VV: This is a gate that consumes Xi and Yi with either a XOR or an AND
			if op.Op == "AND" {
				andGatesOfInputs[bitsA] = op.Res
				foundAndGatesOfInputs++
			} else {
				xorGatesOfInputs[bitsA] = op.Res
				foundXorGatesOfInputs++
			}
		}
	}

	fmt.Printf("ConsumedBy:\n")
	for name, ops := range consumedBy {
		fmt.Printf("  %s: %+v\n", name, ops)
	}

	for bit := range puzzle.ZBits - 1 {
		if _, ok := andGatesOfInputs[bit]; !ok {
			panic(fmt.Sprintf("Missing andGatesOfInput for %d\n", bit))
		}
		if _, ok := xorGatesOfInputs[bit]; !ok {
			panic(fmt.Sprintf("Missing xorGatesOfInputs for %d\n", bit))
		}
	}

	// VV: First validate Z_0 and Z_N

	{
		// VV: C_0 = X_0 XOR Y_0
		op := producedBy["z00"]
		if op.Op != "XOR" || op.A != "x00" || op.B != "y00" {
			fmt.Printf("z00 is not produced by a xor but by a %+v\n", op)
			swappedGates["z00"] = 1
		}
	}

	{
		/*VV:
		Z_N = C_N-1
		C_i = (X_i AND Y_i) OR ((X_i XOR Y_i) AND C_i-1)

		Therefore
		Z_N = (X_(N-1) AND Y_(N-1)) OR ((X_(N-1) XOR Y_(N-1)) AND C_(N-2))
		*/

		name := fmt.Sprintf("z%02d", puzzle.ZBits-1)
		op, ok := producedBy[name]

		if !ok {
			panic("missing last z")
		}

		if op.Op != "OR" {
			fmt.Printf("last Z bit %s is not produced by an or but by a %+v (problem %s)\n", name, op, op.Res)
			swappedGates[op.Res] = 1
		} else {
			left := producedBy[op.A]
			right := producedBy[op.B]

			if left.Op != "AND" || right.Op != "AND" {
				panic("good luck")
			}
		}
	}

	for bit := 1; bit < puzzle.ZBits-1; bit++ {
		/*
			VV: Validate these now

			Z_i = (X_i XOR Y_i) XOR C_i-1
			C_i = (X_i AND Y_i) OR ((X_i XOR Y_i) AND C_i-1)
		*/
		name := fmt.Sprintf("z%02d", bit)

		op, ok := producedBy[name]

		if !ok {
			panic(fmt.Sprintf("missing %s", name))
		}

		if op.Op != "XOR" {
			fmt.Printf("Z bit %s is not produced by a xor but by a %+v (problem %s)\n", name, op, op.Res)
			swappedGates[op.Res] = 1
			continue
		}

		// VV: exactly one of these must be a XOR involving X and Y (the left one)
		myXorIn := xorGatesOfInputs[bit]
		left := producedBy[op.A]
		right := producedBy[op.B]

		if right.Res == myXorIn {
			left, right = right, left
		}

		leftOk := true
		// VV: When left.Res == "" this actually means that left is a swapped Z xor
		// those are handled above
		if left.Res != myXorIn && left.Res != "" {
			fmt.Printf("Z bit %s was supposed to consume the xor %s but it consumed %+v and %+v (problem %s)\n", name, myXorIn, left, right, myXorIn)
			swappedGates[myXorIn] = 1
			leftOk = false
		}

		rightOk := true
		// VV: The right one must be an OR
		if right.Op != "OR" {
			// VV: C_0 is just an and between x00 and y00
			if !(bit == 1 && right.Res == andGatesOfInputs[0]) {
				fmt.Printf("Z bit %s was supposed to consume an or but it consumed %+v and %+v (problem %s)\n", name, left, right, right.Res)
				swappedGates[right.Res] = 1
				rightOk = false
			}
		}

		if !leftOk && rightOk {
			fmt.Printf("Z bit %s was supposed to consume the xor %s but it consumed %+v and %+v, also the right is correct (problem %s)\n", name, myXorIn, left, right, myXorIn)
			swappedGates[left.Res] = 1
		}

		if !rightOk || right.Op != "OR" {
			continue
		}

		// VV: now check C_i
		left = producedBy[right.A]
		right = producedBy[right.B]

		// VV: This check was good enough for my input, not going to bother with more checks
		if left.Op != "AND" {
			fmt.Printf("LEFT C%d  %s was supposed to be an and but it was %+v (problem %s)\n", bit, left.Res, left, left.Res)
			swappedGates[left.Res] = 1
		}

		// VV: Didn't have this in my input
		if right.Op != "AND" {
			fmt.Printf("LEFT C%d  %s was supposed to be an and but it was %+v (problem %s)\n", bit, left.Res, left, left.Res)
			swappedGates[right.Res] = 1
		}
	}

	swappedWireNames := []string{}
	for name := range swappedGates {
		swappedWireNames = append(swappedWireNames, name)
	}
	slices.Sort(swappedWireNames)
	return strings.Join(swappedWireNames, ",")
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
