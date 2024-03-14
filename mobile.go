package masker

// MobileMasker is a masker for mobile
type MobileMasker struct{}

// Marshal masks mobile
// It mask 3 digits from the 4'th digit
// Example:
//
//	MobileMasker{}.Marshal("*", "0987654321") // returns "0987***321"
func (m *MobileMasker) Marshal(s string, i string) string {
	if len(i) == 0 {
		return ""
	}
	return overlay(i, strLoop(s, 3), 4, 7)
}
