package util

func GCD(a, b int) int {
	for {
		if a < b {
			a, b = b, a
		}

		r := a % b
		a, b = b, r
		if b == 0 {
			return a
		}
	}
}

func LCM(a, b int) int {

	lcm := a * b / GCD(a, b)

	if lcm < 0 {
		return -lcm
	}

	return lcm
}
