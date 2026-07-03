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

	if strings.HasPrefix(tag, "tel-") {
		return parseTelMask(maskChar, tag, value)
	}

	return "", false, nil
}

// cleanTelValue 清理電話號碼中的格式字元，移除 "(", ")", " ", "-" 及開頭的 "+"。
func cleanTelValue(value string) string {
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "(", "")
	value = strings.ReplaceAll(value, ")", "")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.TrimPrefix(value, "+")
	return value
}

// parseTelMask 處理 tel- 前綴的動態 tag，依區碼／號碼／國際碼位數遮罩市話號碼。
func parseTelMask(maskChar, tag, value string) (string, bool, error) {
	rest := strings.TrimPrefix(tag, "tel-")
	tokens := strings.Split(rest, "-")

	var intlLen, regionLen, numberLen int
	sep := "-"

	switch len(tokens) {
	case 2:
		// [regionLen, numberLen]
		var err error
		regionLen, err = strconv.Atoi(tokens[0])
		if err != nil || regionLen <= 0 {
			return "", false, fmt.Errorf("invalid mask tag %q: regionLen must be a positive integer", tag)
		}
		numberLen, err = strconv.Atoi(tokens[1])
		if err != nil || numberLen < 4 {
			return "", false, fmt.Errorf("invalid mask tag %q: numberLen must be an integer >= 4", tag)
		}
	case 3:
		if tokens[2] == "dash" || tokens[2] == "space" {
			// [regionLen, numberLen, sep]
			var err error
			regionLen, err = strconv.Atoi(tokens[0])
			if err != nil || regionLen <= 0 {
				return "", false, fmt.Errorf("invalid mask tag %q: regionLen must be a positive integer", tag)
			}
			numberLen, err = strconv.Atoi(tokens[1])
			if err != nil || numberLen < 4 {
				return "", false, fmt.Errorf("invalid mask tag %q: numberLen must be an integer >= 4", tag)
			}
			if tokens[2] == "space" {
				sep = " "
			}
		} else if n, err := strconv.Atoi(tokens[2]); err == nil && n > 0 {
			// [intlLen, regionLen, numberLen]
			intlLen, _ = strconv.Atoi(tokens[0])
			if intlLen <= 0 {
				return "", false, fmt.Errorf("invalid mask tag %q: intlLen must be a positive integer", tag)
			}
			regionLen, _ = strconv.Atoi(tokens[1])
			if regionLen <= 0 {
				return "", false, fmt.Errorf("invalid mask tag %q: regionLen must be a positive integer", tag)
			}
			numberLen = n
			if numberLen < 4 {
				return "", false, fmt.Errorf("invalid mask tag %q: numberLen must be an integer >= 4", tag)
			}
		} else {
			return "", false, fmt.Errorf("invalid mask tag %q: third token must be a positive integer or dash/space", tag)
		}
	case 4:
		// [intlLen, regionLen, numberLen, sep]
		var err error
		intlLen, err = strconv.Atoi(tokens[0])
		if err != nil || intlLen <= 0 {
			return "", false, fmt.Errorf("invalid mask tag %q: intlLen must be a positive integer", tag)
		}
		regionLen, err = strconv.Atoi(tokens[1])
		if err != nil || regionLen <= 0 {
			return "", false, fmt.Errorf("invalid mask tag %q: regionLen must be a positive integer", tag)
		}
		numberLen, err = strconv.Atoi(tokens[2])
		if err != nil || numberLen < 4 {
			return "", false, fmt.Errorf("invalid mask tag %q: numberLen must be an integer >= 4", tag)
		}
		if tokens[3] != "dash" && tokens[3] != "space" {
			return "", false, fmt.Errorf("invalid mask tag %q: separator must be dash or space", tag)
		}
		if tokens[3] == "space" {
			sep = " "
		}
	default:
		return "", false, fmt.Errorf("invalid mask tag %q: expected 2-4 tokens after tel-", tag)
	}

	cleaned := cleanTelValue(value)
	l := len([]rune(cleaned))
	if l == 0 {
		return "", true, nil
	}

	total := intlLen + regionLen + numberLen
	if l != total {
		return cleaned, true, nil
	}

	r := []rune(cleaned)
	var b strings.Builder
	pos := 0

	if intlLen > 0 {
		b.WriteString("+")
		b.WriteString(string(r[pos : pos+intlLen]))
		b.WriteString(sep)
		pos += intlLen
	}

	b.WriteString(string(r[pos : pos+regionLen]))
	b.WriteString(sep)
	pos += regionLen

	prefix := numberLen - 4
	if prefix > 0 {
		b.WriteString(string(r[pos : pos+prefix]))
		b.WriteString(sep)
	}

	b.WriteString(strLoop(maskChar, 4))

	return b.String(), true, nil
}
