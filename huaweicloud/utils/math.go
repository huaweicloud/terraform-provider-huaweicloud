package utils

// Power is a method for calculating powers of exponents.
// The result is base^exponent (^ is the exponential operator).
func Power(base int, exponent int) int {
	if exponent == 0 {
		return 1
	}
	return base * Power(base, exponent-1)
}
