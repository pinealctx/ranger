package jda

import "math"

const (
	calMinBuff     float64 = 1e-5 // 0.00001
	usdMinBuff     float64 = 0.1  // we use 0.01 as the usd unit, so 0.015 is a good buffer
	ratioDiff              = 0.002
	rateMinBuff    float64 = 1e-8 // 0.00000001, used for rate comparison
	twoDecimalBUff         = 0.01
)

func equalsF64(a, b float64) bool {
	return _equalsF64WithBuff(a, b, calMinBuff)
}

func equalsPrice(a, b float64) bool {
	// for price comparison, we use a larger buffer
	return _equalsF64WithBuff(a, b, rateMinBuff)
}

func equalsUsd(a, b float64) bool {
	return _equalsF64WithBuff(a, b, usdMinBuff)
}

func equals2Decimal(a, b float64) bool {
	return _equalsF64WithBuff(a, b, twoDecimalBUff)
}

func ratioEqual(a, b float64) bool {
	if (b == 0 && a != 0) || (a == 0 && b != 0) {
		return false
	}
	return a/b > (1-ratioDiff) && a/b < (1+ratioDiff)
}

func _equalsF64WithBuff(a, b, buff float64) bool {
	// if a, b are different direction, one > 0, another < 0, return false
	if (a > 0 && b < 0) || (a < 0 && b > 0) {
		return false
	}
	x := math.Abs(a - b)
	return x < buff
}
