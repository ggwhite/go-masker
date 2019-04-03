package cmd

import "strings"

// Email mask email
func Email(input string) string {
	splitInput := strings.Split(input, "@")
	address := splitInput[0]
	domain := splitInput[1]

	address = overlay(address, "****", 3, 7)

	return address + "@" + domain
}
