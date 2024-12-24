package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Operation struct {
	Res, A, Op, B string
}

type Puzzle struct {
	WireValues map[string]int
	Operations []Operation
	ZBits      int
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	wireValues := map[string]int{}
	operations := []Operation{}
	zbits := 0

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if strings.Contains(line, ":") {
			inputWireValue := strings.Split(line, ": ")
			value, err := strconv.Atoi(inputWireValue[1])

			if err != nil {
				return &Puzzle{
					Operations: operations, ZBits: zbits, WireValues: wireValues,
				}, err
			}
			wireName := inputWireValue[0]
			wireValues[wireName] = value

			if bitsStr := strings.TrimPrefix(wireName, "z"); bitsStr != wireName {
				bits, err := strconv.Atoi(bitsStr)
				if err != nil {
					return &Puzzle{
						Operations: operations, ZBits: zbits, WireValues: wireValues,
					}, fmt.Errorf("invalid Z wire in line %s due to %w", line, err)
				}
				zbits = max(zbits, bits+1)
			}
		} else {
			var a string
			var b string
			var op string
			var res string

			n, err := fmt.Sscanf(line, "%s %s %s -> %s", &a, &op, &b, &res)
			if err != nil {
				return &Puzzle{
					Operations: operations, ZBits: zbits, WireValues: wireValues,
				}, err
			}

			if n != 4 {
				return &Puzzle{
					Operations: operations, ZBits: zbits, WireValues: wireValues,
				}, fmt.Errorf("parsed %d/4 elements from line %s", n, line)
			}

			// VV: Make my life easier by sorting the names of operands, this also
			// converts Y_i $OP X_i into X_I $OP Y_I
			if a > b {
				a, b = b, a
			}

			operations = append(operations, Operation{A: a, B: b, Op: op, Res: res})
			if bitsStr := strings.TrimPrefix(res, "z"); bitsStr != res {
				bits, _ := strconv.Atoi(bitsStr)
				if err != nil {
					return &Puzzle{
						Operations: operations, ZBits: zbits, WireValues: wireValues,
					}, fmt.Errorf("invalid Z wire in line %s due to %w", line, err)
				}
				zbits = max(zbits, bits+1)
			}

		}
	}

	return &Puzzle{
		Operations: operations, ZBits: zbits, WireValues: wireValues,
	}, scanner.Err()
}

func ReadString(text string) (*Puzzle, error) {
	scanner := bufio.NewScanner(strings.NewReader(text))

	return ReadScanner(scanner)
}

func ReadInputFile(path string) (*Puzzle, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	return ReadScanner(scanner)
}
