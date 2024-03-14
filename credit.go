package masker

type CreditMasker struct{}

// Marshal masks credit card number
// It mask 6 digits from the 7'th digit
// Example:
//
//	CreditMasker{}.Marshal("*", "4111111111111111") // returns "**** **** **** 1111"
func (m *CreditMasker) Marshal(s string, i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return overlay(i, strLoop(s, 6), 6, 12)
}
