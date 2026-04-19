package masker

import (
	"fmt"
	"strconv"
	"strings"
)

// AllMasker is a masker that replaces every character with the mask char.
type AllMasker struct{}

// Marshal replaces every character in the value with the mask char.
//
// Example:
//
//	AllMasker{}.Marshal("*", "secret") // returns "******"
func (m *AllMasker) Marshal(s, i string) string {
	return strLoop(s, len([]rune(i)))
}

// parseGenericMask handles dynamic tag patterns that do not map to a registered Masker:
//   - "all"      — mask every character
//   - "first-N"  — mask the first N characters
//   - "last-N"   — mask the last N characters
//
// Returns (maskedValue, true, nil) on match, ("", false, nil) when the tag is not a dynamic pattern,
// or ("", false, err) when the pattern is recognized but N is invalid.
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
