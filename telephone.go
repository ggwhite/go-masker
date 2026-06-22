package masker

import "strings"

// TelephoneMasker is a masker for telephone
type TelephoneMasker struct{}

// Marshal masks telephone
// It remove "(", ")", " ", "-" chart, and mask last 4 digits of telephone number, format to "(??)????-????"
// Example:
//
//	TelephoneMasker{}.Marshal("*", "0227993078") // returns "(02)2799-****"
func (m *TelephoneMasker) Marshal(s string, i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}

	i = strings.ReplaceAll(i, " ", "")
	i = strings.ReplaceAll(i, "(", "")
	i = strings.ReplaceAll(i, ")", "")
	i = strings.ReplaceAll(i, "-", "")

	r := []rune(i)
	l = len(r)

	if l != 10 && l != 8 {
		return i
	}

	ans := ""

	if l == 10 {
		ans += "("
		ans += string(r[:2])
		ans += ")"
		r = r[2:]
	}

	ans += string(r[:4])
	ans += "-"
	ans += strLoop(s, 4)

	return ans
}
