package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Puzzle struct {
	Codes []uint64
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	codes := []uint64{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		code, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return &Puzzle{
				Codes: codes,
			}, fmt.Errorf("line %s raised error %w", line, err)
		}
		codes = append(codes, uint64(code))
	}

	return &Puzzle{
		Codes: codes,
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
