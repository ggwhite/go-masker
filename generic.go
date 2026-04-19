package masker

import (
	"fmt"
	"strconv"
	"strings"
)

// AllMasker masks every character in the string
type AllMasker struct{}

func (m *AllMasker) Marshal(s, i string) string {
	return strLoop(s, len([]rune(i)))
}

// parseGenericMask parses dynamic tag patterns: "all", "first-N", "last-N"
// Returns the masked string and true if the tag matched, otherwise empty string and false.
func parseGenericMask(maskChar, tag, value string) (string, bool, error) {
	if tag == "all" {
		return strLoop(maskChar, len([]rune(value))), true, nil
	}

	if strings.HasPrefix(tag, "first-") {
		n, err := strconv.Atoi(strings.TrimPrefix(tag, "first-"))
		if err != nil || n < 0 {
			return "", false, fmt.Errorf("invalid mask tag %q: N must be a non-negative integer", tag)
		}
		r := []rune(value)
		if n > len(r) {
			n = len(r)
		}
		return overlay(value, strLoop(maskChar, n), 0, n), true, nil
	}

	if strings.HasPrefix(tag, "last-") {
		n, err := strconv.Atoi(strings.TrimPrefix(tag, "last-"))
		if err != nil || n < 0 {
			return "", false, fmt.Errorf("invalid mask tag %q: N must be a non-negative integer", tag)
		}
		r := []rune(value)
		start := len(r) - n
		if start < 0 {
			start = 0
		}
		return overlay(value, strLoop(maskChar, len(r)-start), start, len(r)), true, nil
	}

	return "", false, nil
}
