package masker

import (
	"fmt"
	"strconv"
	"strings"
)

// AllMasker 是把每個字元都換成遮罩字元的 masker。
type AllMasker struct {
	mask string
}

// Mask 將值中的每個字元都換成遮罩字元。
// Example:
//
//	(&AllMasker{mask: "*"}).Mask("secret") // returns "******"
func (m *AllMasker) Mask(value string) string {
	return strLoop(m.mask, len([]rune(value)))
}

// parseGenericMask 處理不對應已註冊 Masker 的動態 tag 樣式：
//   - "first-N"  — 遮罩前 N 個字元
//   - "last-N"   — 遮罩後 N 個字元
//
// 命中時回傳 (遮罩值, true, nil)；非動態樣式回傳 ("", false, nil)；
// 樣式被識別但 N 不合法（負數或非數字）時回傳 ("", false, err)。
func parseGenericMask(maskChar, tag, value string) (string, bool, error) {
	if strings.HasPrefix(tag, "first-") {
		n, err := strconv.Atoi(strings.TrimPrefix(tag, "first-"))
		if err != nil || n < 0 {
			return "", false, fmt.Errorf("invalid mask tag %q: N must be a non-negative integer", tag)
		}
		r := []rune(value)
		if n > len(r) {
			n = len(r)
		}
		return overlay(value, strLoop(maskChar, n), 0, n), true, nil
	}

	if strings.HasPrefix(tag, "last-") {
		n, err := strconv.Atoi(strings.TrimPrefix(tag, "last-"))
		if err != nil || n < 0 {
			return "", false, fmt.Errorf("invalid mask tag %q: N must be a non-negative integer", tag)
		}
		r := []rune(value)
		start := len(r) - n
		if start < 0 {
			start = 0
		}
		return overlay(value, strLoop(maskChar, len(r)-start), start, len(r)), true, nil
	}

	return "", false, nil
}
