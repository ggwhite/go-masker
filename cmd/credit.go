package cmd

// Credit mask credit card number
func Credit(input string) string {
	return overlay(input, "******", 6, 12)
}
