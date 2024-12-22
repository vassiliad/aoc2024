package main

import (
	"puzzle/util"
	"testing"
)

func TestSmall(t *testing.T) {
	small := `
1
10
100
2024
`

	// println(small)

	puzzle, err := util.ReadString(small)

	logger := SetupLogger()
	logger.Printf("Input is %+v\n", puzzle)

	if err != nil {
		t.Fatal("Ran into problems while reading input. Problem", err)
	}

	sol := solution(puzzle, 2000, logger)
	const correctAnswer = 37327623

	if sol != correctAnswer {
		t.Fatal("Expected answer to be", correctAnswer, "but it was", sol)
	}
}

func TestIndividual(t *testing.T) {
	expected := map[uint64]uint64{
		1:    8685429,
		10:   4700978,
		100:  15273692,
		2024: 8667524,
	}

	for seed, correctAnswer := range expected {
		sol := generateRandom(seed, 2000)

		if sol != correctAnswer {
			t.Fatal("Expected answer to be", correctAnswer, "but it was", sol)
		}
	}

}

func TestSingleStep(t *testing.T) {
	expected := map[uint64]uint64{
		123:      15887950,
		15887950: 16495136,
		16495136: 527345,
		527345:   704524,
		704524:   1553684,
		1553684:  12683156,
		12683156: 11100544,
		11100544: 12249484,
		12249484: 7753432,
		7753432:  5908254,
	}

	wrong := 0
	for seed, correctAnswer := range expected {
		sol := generateRandom(seed, 1)

		if sol != correctAnswer {
			wrong++
			t.Log("Expected answer to be", correctAnswer, "but it was", sol)
		}
	}

	if wrong > 0 {
		t.Fatal("mistakes:", wrong)
	}

}
