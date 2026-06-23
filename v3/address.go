package masker

// AddressMasker 是地址的 masker。
type AddressMasker struct {
	mask string
}

// Mask 遮罩地址，遮罩最後 6 個字元。
// Example:
//
//	(&AddressMasker{mask: "*"}).Mask("台北市內湖區內湖路一段737巷1號1樓") // returns "台北市內湖區******"
func (m *AddressMasker) Mask(value string) string {
	l := len([]rune(value))
	if l == 0 {
		return ""
	}
	if l <= 6 {
		return strLoop(m.mask, 6)
	}
	return overlay(value, strLoop(m.mask, 6), 6, l)
}
