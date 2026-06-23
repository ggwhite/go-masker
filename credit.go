package masker

// CreditMasker 是信用卡號的 masker，遮罩第 7–12 位（1-indexed），保留前 6 後 4 碼。
type CreditMasker struct {
	mask string
}

// Mask 遮罩信用卡號，將中間 6 碼換成遮罩字元。
// Example:
//
//	(&CreditMasker{mask: "*"}).Mask("4111111111111111") // returns "411111******1111"
func (m *CreditMasker) Mask(value string) string {
	l := len([]rune(value))
	if l == 0 {
		return ""
	}
	return overlay(value, strLoop(m.mask, 6), 6, 12)
}
