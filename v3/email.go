package masker

import "strings"

// EmailMasker 是 email 的 masker。
type EmailMasker struct {
	mask string
}

// Mask 遮罩 email，保留網域與前 3 個字元。
// Example:
//
//	(&EmailMasker{mask: "*"}).Mask("ggw.chang@gmail.com") // returns "ggw****@gmail.com"
func (m *EmailMasker) Mask(value string) string {
	l := len([]rune(value))
	if l == 0 {
		return ""
	}

	tmp := strings.SplitN(value, "@", 2)
	if len(tmp) == 1 {
		return overlay(value, strLoop(m.mask, 4), 3, 7)
	}

	addr := tmp[0]
	domain := tmp[1]

	addr = overlay(addr, strLoop(m.mask, 4), 3, len(addr))

	return addr + "@" + domain
}
