package masker

import "strings"

// EmailMasker is a masker for email
type EmailMasker struct{}

// Marshal masks email
// It keep domain and the first 3 letters
// Example:
//
//	AddressMasker{}.Marshal("*", "ggw.chang@gmail.com") // returns "ggw****@gmail.com"
func (m *EmailMasker) Marshal(s string, i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}

	tmp := strings.Split(i, "@")
	if len(tmp) == 1 {
		return overlay(i, strLoop(s, 4), 3, 7)
	}

	addr := tmp[0]
	domain := tmp[1]

	addr = overlay(addr, strLoop(s, 4), 3, 7)

	return addr + "@" + domain
}
