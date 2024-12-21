package main

import (
	"fmt"
	"image"
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `
029A
980A
179A
456A
379A
`

	// println(small)

	puzzle, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", puzzle)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, logger)
	const correctAnswer = 126384

	if sol != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", sol)
	}
}

func TestTiny(t *testing.T) {
	small := `
179A
`
	puzzle, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", puzzle)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	buttonsFinal := [][]rune{
		{'7', '8', '9'},
		{'4', '5', '6'},
		{'1', '2', '3'},
		{'x', '0', 'A'},
	}

	buttons := [][]rune{
		{'x', '^', 'A'},
		{'<', 'v', '>'},
	}

	firstRobot := calculateCode(buttons, []rune("<v<A>>^A<vA<A>>^AAvAA<^A>A<v<A>>^AAvA^A<vA>^AA<A>A<v<A>A>^AAAvA<^A>A"), image.Pt(2, 0))
	println("First robot:", string(firstRobot))

	secondRobot := calculateCode(buttons, firstRobot, image.Pt(2, 0))
	println("First robot:", string(secondRobot))

	finalCode := calculateCode(buttonsFinal, secondRobot, image.Pt(2, 3))
	println("Final code:", string(finalCode))

	sol := solution(puzzle, logger)

	const correctAnswer = 68 * 179

	if sol != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", sol)
	}
}

func calculateCode(buttons [][]rune, code []rune, start image.Point) []rune {
	ret := []rune{}

	deltas := map[rune]image.Point{
		'>': image.Pt(1, 0),
		'<': image.Pt(-1, 0),
		'v': image.Pt(0, 1),
		'^': image.Pt(0, -1),
	}

	for _, c := range code {
		if c == 'A' {
			ret = append(ret, buttons[start.Y][start.X])
		} else {
			if buttons[start.Y][start.X] == 'x' {
				panic("oh no")
			}
			if delta, ok := deltas[c]; ok {
				start = start.Add(delta)
			} else {
				panic(fmt.Sprintf("unknown button '%s'", string(c)))
			}

			if buttons[start.Y][start.X] == 'x' {
				panic("oh no 2")
			}
		}
	}

	return ret
}
