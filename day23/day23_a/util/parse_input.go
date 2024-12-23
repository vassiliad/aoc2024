package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Connection struct {
	A, B string
}

type Puzzle struct {
	Connections []Connection
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	connections := []Connection{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		nodes := strings.Split(line, "-")

		if len(nodes) != 2 {
			return &Puzzle{
				Connections: connections,
			}, fmt.Errorf("line \"%s\" does not contain 2 nodes", line)
		}

		connections = append(connections, Connection{A: nodes[0], B: nodes[1]})
	}

	return &Puzzle{
		Connections: connections,
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
