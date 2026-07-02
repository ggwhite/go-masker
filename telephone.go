package masker

import "strings"

// TelephoneMasker 是市話的 masker。
type TelephoneMasker struct {
	mask          string
	regionCodeLen int
	numberLen     int
}

// Mask 遮罩市話。
// 當未設定 regionCodeLen 和 numberLen 時（零值），使用預設邏輯：
//   - 移除 "("、")"、" "、"-" 等特殊字元
//   - 支援 8 位和 10 位市話
//   - 10位格式："(XX)XXXX-****"，8位格式："XXXX-****"
//   - 長度不匹配時返回原值
//
// 當設定 regionCodeLen 和 numberLen 時，使用可配置邏輯：
//   - 僅保留數字
//   - 格式："XX XXXX****"（有区号）或 "XXXX****"（無区号）
//   - 長度不匹配時全部遮罩
//
// Example:
//
//	(&TelephoneMasker{mask: "*"}).Mask("0227993078")                      // returns "(02)2799-****"
//	(&TelephoneMasker{mask: "*", regionCodeLen: 2, numberLen: 8}).Mask("0227993078") // returns "02 2799****"
//	(&TelephoneMasker{mask: "*", regionCodeLen: 0, numberLen: 8}).Mask("12345678")   // returns "1234****"
func (m *TelephoneMasker) Mask(value string) string {
	if m.regionCodeLen == 0 && m.numberLen == 0 {
		return m.maskDefault(value)
	}
	return m.maskConfigurable(value)
}

func (m *TelephoneMasker) maskDefault(value string) string {
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
	ans += strings.Repeat(m.mask, 4)

	return ans
}

func (m *TelephoneMasker) maskConfigurable(value string) string {
	digits := extractDigits(value)
	expectedLen := m.regionCodeLen + m.numberLen

	if len(digits) != expectedLen {
		return strings.Repeat(m.mask, len(digits))
	}

	if m.numberLen < 4 {
		return strings.Repeat(m.mask, len(digits))
	}

	start := m.regionCodeLen
	end := start + m.numberLen - 4

	var result strings.Builder
	if m.regionCodeLen > 0 {
		result.WriteString(digits[:start])
		result.WriteByte(' ')
	}
	result.WriteString(digits[start:end])
	result.WriteString(strings.Repeat(m.mask, 4))
	return result.String()
}

func extractDigits(s string) string {
	var builder strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
