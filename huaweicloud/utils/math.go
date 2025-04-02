package utils

import "math"

// A sufficiently small error value.
const epsilon = 1e-9

// Power is a method for calculating powers of exponents.
// The result is base^exponent (^ is the exponential operator).
func Power(base int, exponent int) int {
	if exponent == 0 {
		return 1
	}
	return base * Power(base, exponent-1)
}

// Round is a method that used to round a floating point number to a specified number of decimal places.
// Argument inputs:
// + num: the number to be rounded.
// + precision: the number of decimal places to be retained.
func Round(num float64, precision int) float64 {
	// Calculate the precision of 10 as the base
	base := float64(Power(10, precision))

	absValue := math.Abs(num)
	// Multiply by the base first, round it off, then divide by the base
	rounded := math.Round(absValue*base) / base

	if num < 0 {
		return -rounded
	}
	return rounded
}

// EqualFloat is a method that used to compare two floating point numbers for equality (accounting for precision).
func EqualFloat(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}
