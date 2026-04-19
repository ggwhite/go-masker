package masker

// CreditMasker is a masker for credit card numbers.
// It masks characters at positions 7–12 (1-indexed), preserving the first 6 and last 4 digits.
type CreditMasker struct{}

// Marshal masks a credit card number by replacing the middle 6 digits with the mask char.
//
// Example:
//
//	CreditMasker{}.Marshal("*", "4111111111111111") // returns "411111******1111"
func (m *CreditMasker) Marshal(s string, i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return overlay(i, strLoop(s, 6), 6, 12)
}
