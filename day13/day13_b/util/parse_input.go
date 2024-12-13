package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X, Y int64
}

type Gamba struct {
	ButtonA, ButtonB, Prize Point
}

type Puzzle struct {
	Gambas []Gamba
}

func parsePoint(text string) (Point, error) {
	parts := strings.Split(text, ", ")

	if parts[0][0] != 'X' {
		return Point{}, fmt.Errorf("the X offset does not start with X, value was %s", text)
	}

	if parts[1][0] != 'Y' {
		return Point{}, fmt.Errorf("the Y offset does not start with Y, value was %s", text)
	}

	x, errx := strconv.Atoi(parts[0][2:])

	if errx != nil {
		return Point{}, errx
	}

	y, erry := strconv.Atoi(parts[1][2:])

	if erry != nil {
		return Point{}, erry
	}

	return Point{X: int64(x), Y: int64(y)}, nil
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	gambas := []Gamba{}
	gamba := Gamba{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}
		var err error = nil
		if strings.HasPrefix(line, "Button A: ") {
			gamba.ButtonA, err = parsePoint(line[10:])
			if err != nil {
				return &Puzzle{Gambas: gambas}, err
			}
		} else if strings.HasPrefix(line, "Button B: ") {
			gamba.ButtonB, err = parsePoint(line[10:])
			if err != nil {
				return &Puzzle{Gambas: gambas}, err
			}
		} else if strings.HasPrefix(line, "Prize: ") {
			gamba.Prize, err = parsePoint(line[7:])
			gamba.Prize.X += 10000000000000
			gamba.Prize.Y += 10000000000000

			if err != nil {
				return &Puzzle{Gambas: gambas}, err
			}
			gambas = append(gambas, gamba)
		} else {
			panic(line)
		}
	}

	return &Puzzle{
		Gambas: gambas,
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
