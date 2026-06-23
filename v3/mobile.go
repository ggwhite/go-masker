package masker

// MobileMasker 是手機號碼的 masker。
type MobileMasker struct {
	mask string
}

// Mask 遮罩手機號碼，從第 4 位起遮罩 3 碼。
// Example:
//
//	(&MobileMasker{mask: "*"}).Mask("0987654321") // returns "0987***321"
func (m *MobileMasker) Mask(value string) string {
	if len(value) == 0 {
		return ""
	}
	return overlay(value, strLoop(m.mask, 3), 4, 7)
}
