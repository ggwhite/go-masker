package masker

// IDMasker is a masker for ID
type IDMasker struct{}

// Marshal masks ID
// It mask last 4 digits of ID number
// Example:
//
//	IDMasker{}.Marshal("*", "1234567890") // returns "1234****890"
func (m *IDMasker) Marshal(s string, i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return overlay(i, strLoop(s, 4), 6, 10)
}
