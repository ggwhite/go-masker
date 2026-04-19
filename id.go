package masker

// IDMasker is a masker for ID numbers (e.g. national ID, passport).
// It masks characters at positions 7–10 (1-indexed).
type IDMasker struct{}

// Marshal masks an ID number by replacing characters at positions 7–10 with the mask char.
//
// Example:
//
//	IDMasker{}.Marshal("*", "A123456789") // returns "A12345****"
func (m *IDMasker) Marshal(s string, i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return overlay(i, strLoop(s, 4), 6, 10)
}
