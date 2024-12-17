package util

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Puzzle struct {
	Program   []int64
	Registers [3]int64
	PC        int64
	Output    []int64
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	registers := [3]int64{0, 0, 0}
	program := []int64{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "Register ") {
			val, err := strconv.Atoi(line[len("Register A: "):])
			if err != nil {
				return &Puzzle{
					Registers: registers, Program: program,
				}, err
			}
			idx := int([]rune(line)[len("Register ")] - 'A')

			registers[idx] = int64(val)
		} else {
			for _, s := range strings.Split(line[len("Program: "):], ",") {
				opcode, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return &Puzzle{
						Registers: registers, Program: program,
					}, err
				}

				program = append(program, opcode)
			}
		}
	}

	return &Puzzle{
		Registers: registers, Program: program,
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
