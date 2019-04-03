package cmd

// Mobile mask mobile
func Mobile(input string) string {
	return overlay(input, "***", 4, 7)
}
