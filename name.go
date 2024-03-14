package masker

import "strings"

type NameMasker struct{}

// Marshal masks name
// It mask the second letter and the third letter
// Example:
//
//	NameMasker{}.Marshal("*", "name") // returns "n**e"
//	NameMasker{}.Marshal("*", "ABCD") // returns "A**D"
func (m *NameMasker) Marshal(s string, i string) string {
	l := len([]rune(i))

	if l == 0 {
		return ""
	}

	// if has space
	if strs := strings.Split(i, " "); len(strs) > 1 {
		tmp := make([]string, len(strs))
		for idx, str := range strs {
			tmp[idx] = m.Marshal(s, str)
		}
		return strings.Join(tmp, " ")
	}

	if l == 2 || l == 3 {
		return overlay(i, strLoop(s, 2), 1, 2)
	}

	if l > 3 {
		return overlay(i, strLoop(s, 2), 1, 3)
	}

	return strLoop(s, 2)
}
