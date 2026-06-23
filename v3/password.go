package masker

// PasswordMasker 是密碼的 masker。
type PasswordMasker struct {
	mask string
}

// Mask 遮罩密碼，固定回傳 14 個遮罩字元，與原密碼長度無關。
// Example:
//
//	(&PasswordMasker{mask: "*"}).Mask("password") // returns "**************"
func (m *PasswordMasker) Mask(value string) string {
	return strLoop(m.mask, 14)
}
