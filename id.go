package masker

// IDMasker 是身分證、護照等證號的 masker，遮罩第 7–10 位（1-indexed）。
type IDMasker struct {
	mask string
}

// Mask 遮罩證號，將第 7–10 位字元換成遮罩字元。
// Example:
//
//	(&IDMasker{mask: "*"}).Mask("A123456789") // returns "A12345****"
func (m *IDMasker) Mask(value string) string {
	l := len([]rune(value))
	if l == 0 {
		return ""
	}
	return overlay(value, strLoop(m.mask, 4), 6, 10)
}
