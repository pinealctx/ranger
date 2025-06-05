package jda

import "math"

const (
	calMinBuff float64 = 1e-5 // 0.00001
	usdMinBuff float64 = 0.1  // we use 0.01 as the usd unit, so 0.015 is a good buffer
	ratioDiff          = 0.002
)

func equalsF64(a, b float64) bool {
	// if a, b are different direction, one > 0, another < 0, return false
	if (a > 0 && b < 0) || (a < 0 && b > 0) {
		return false
	}
	return math.Abs(a-b) < calMinBuff
}

func equalsUsd(a, b float64) bool {
	if (a > 0 && b < 0) || (a < 0 && b > 0) {
		return false
	}
	return math.Abs(a-b) < usdMinBuff
}

func ratioEqual(a, b float64) bool {
	if (b == 0 && a != 0) || (a == 0 && b != 0) {
		return false
	}
	return a/b > (1-ratioDiff) && a/b < (1+ratioDiff)
}
