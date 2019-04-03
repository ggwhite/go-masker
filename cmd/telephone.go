package cmd

import "strings"

// Telephone mask telephone
func Telephone(input string) string {
	input = strings.Replace(input, " ", "", -1)
	input = strings.Replace(input, "(", "", -1)
	input = strings.Replace(input, ")", "", -1)
	input = strings.Replace(input, "-", "", -1)

	length := len(input)

	if length != 10 && length != 8 {
		return input
	}

	result := ""

	if length == 10 {
		result += "("
		result += input[:2]
		result += ")"
		input = input[2:]
	}

	result += input[:4]
	result += "-"
	result += "****"

	return result
}
