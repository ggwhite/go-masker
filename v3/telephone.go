package masker

import "strings"

// TelephoneMasker 是市話的 masker。
type TelephoneMasker struct {
	mask string
}

// Mask 遮罩市話，移除 "("、")"、" "、"-" 後遮罩最後 4 碼，格式化為 "(XX)XXXX-****"。
// Example:
//
//	(&TelephoneMasker{mask: "*"}).Mask("0227993078") // returns "(02)2799-****"
func (m *TelephoneMasker) Mask(value string) string {
	l := len([]rune(value))
	if l == 0 {
		return ""
	}

	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "(", "")
	value = strings.ReplaceAll(value, ")", "")
	value = strings.ReplaceAll(value, "-", "")

	r := []rune(value)
	l = len(r)

	if l != 10 && l != 8 {
		return value
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
	ans += strLoop(m.mask, 4)

	return ans
}
