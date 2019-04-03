package cmd

import "math"

// Address mask address
func Address(input string) string {
	length := len([]rune(input))
	if length <= 6 {
		return "******"
	}
	return overlay(input, "******", 6, math.MaxInt64)
}
