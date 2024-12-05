package util

import (
	"bufio"
	"image"
	"os"
	"strconv"
	"strings"
)

type Puzzle struct {
	Rules   []image.Point
	Updates [][]int
}

func ReadScanner(scanner *bufio.Scanner) (*Puzzle, error) {
	rules := []image.Point{}
	updates := [][]int{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			break
		}

		parts := strings.Split(line, "|")

		x, err_l := strconv.Atoi(parts[0])

		if err_l != nil {
			return &Puzzle{
				Rules:   rules,
				Updates: updates,
			}, err_l
		}

		y, err_r := strconv.Atoi(parts[1])

		if err_r != nil {
			return &Puzzle{
				Rules:   rules,
				Updates: updates,
			}, scanner.Err()
		}

		rules = append(rules, image.Pt(x, y))
	}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		pages := []int{}

		for _, n := range strings.Split(line, ",") {
			p, err := strconv.Atoi(n)

			if err != nil {
				updates = append(updates, pages)
				return &Puzzle{
					Rules:   rules,
					Updates: updates,
				}, err
			}

			pages = append(pages, p)
		}

		updates = append(updates, pages)
	}

	return &Puzzle{
		Rules:   rules,
		Updates: updates,
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
