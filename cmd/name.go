package cmd

// Name mask name
func Name(input string) string {

	length := len([]rune(input))

	if length == 2 || length == 3 {
		return overlay(input, "*", 1, 2)
	}

	if length > 3 {
		return overlay(input, "**", 1, 3)
	}

	return input
}
