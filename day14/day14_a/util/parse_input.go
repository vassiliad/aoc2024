package util

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"
)

type Bot struct {
	Position, Velocity image.Point
}

type Puzzle struct {
	Bots []Bot
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
	bots := []Bot{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}
		if !strings.HasPrefix(line, "p=") {
			return &Puzzle{
				Bots: bots,
			}, fmt.Errorf("line %s does not begin with p=", line)
		}

		parts := strings.Split(line[2:], " v=")

		if len(parts) != 2 {
			return &Puzzle{
				Bots: bots,
			}, fmt.Errorf("line %s does not contain exactly one v=", line)
		}

		pos, errPos := parsePoint(parts[0])
		if errPos != nil {
			return &Puzzle{
				Bots: bots,
			}, fmt.Errorf("line %s contains invalid position %s", line, parts[0])
		}
		vel, errVel := parsePoint(parts[1])
		if errVel != nil {
			return &Puzzle{
				Bots: bots,
			}, fmt.Errorf("line %s contains invalid velocity %s", line, parts[1])
		}

		bots = append(bots, Bot{Position: pos, Velocity: vel})
	}

	return &Puzzle{
		Bots: bots,
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
