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

	i = strings.Replace(i, " ", "", -1)
	i = strings.Replace(i, "(", "", -1)
	i = strings.Replace(i, ")", "", -1)
	i = strings.Replace(i, "-", "", -1)

	l = len([]rune(i))

	if l != 10 && l != 8 {
		return i
	}

	ans := ""

	if l == 10 {
		ans += "("
		ans += i[:2]
		ans += ")"
		i = i[2:]
	}

	ans += i[:4]
	ans += "-"
	ans += "****"

	return ans
}
