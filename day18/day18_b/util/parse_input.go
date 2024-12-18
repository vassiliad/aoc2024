package util

import (
	"bufio"
	"image"
	"os"
	"strconv"
	"strings"
)

type Puzzle struct {
	Bytes []image.Point
}

func parsePoint(text string) (image.Point, error) {
	parts := strings.Split(text, ",")

	x, errx := strconv.Atoi(parts[0])

	if errx != nil {
		return image.Point{}, errx
	}

	y, erry := strconv.Atoi(parts[1])

	if erry != nil {
		return image.Point{}, erry
	}

	return image.Pt(x, y), nil
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	bytes := []image.Point{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}
		pt, err := parsePoint(line)

		if err != nil {
			return &Puzzle{
				Bytes: bytes,
			}, err
		}

		bytes = append(bytes, pt)
	}

	return &Puzzle{
		Bytes: bytes,
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
