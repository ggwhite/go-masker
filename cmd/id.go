package cmd

// ID mask ID
func ID(input string) string {
	return overlay(input, "****", 6, 10)
}
