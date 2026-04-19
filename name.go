package masker

import "strings"

// NameMasker is a masker for names.
// It masks the 2nd and 3rd characters, preserving the first and last.
// Space-separated words (e.g. full names) are each masked independently.
type NameMasker struct{}

// Marshal masks a name by replacing the middle characters with the mask char.
//
// Rules:
//   - length 1: returns 2 mask chars
//   - length 2–3: masks the 2nd character
//   - length 4+: masks the 2nd and 3rd characters
//   - spaces: each word is masked independently
//
// Example:
//
//	NameMasker{}.Marshal("*", "John")     // returns "J**n"
//	NameMasker{}.Marshal("*", "John Doe") // returns "J**n D**e"
func (m *NameMasker) Marshal(s string, i string) string {
	l := len([]rune(i))

	if l == 0 {
		return ""
	}

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
