package util

import "math"

func NumberOfDigits(number int) int {
	digits := 0

	for number > 0 {
		digits++
		number /= 10
	}

	return digits
}

func SplitNumber(number, digits int) (int, int) {
	power := int(math.Pow10(digits))

	right := number % power
	left := number / power

	return left, right
}
