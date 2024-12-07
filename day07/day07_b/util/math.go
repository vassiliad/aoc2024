package util

func NumberOfDigits(number int) int {
	digits := 0

	for number > 0 {
		digits++
		number /= 10
	}

	return digits
}
