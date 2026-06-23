package masker

import "strings"

// NameMasker 是姓名的 masker，遮罩第 2、3 個字元並保留首尾。
// 以空白分隔的字（例如全名）各自獨立遮罩。
type NameMasker struct {
	mask string
}

// Mask 遮罩姓名，將中間字元換成遮罩字元。
//
// 規則：
//   - 長度 1：回傳 2 個遮罩字元
//   - 長度 2–3：遮罩第 2 個字元
//   - 長度 4 以上：遮罩第 2、3 個字元
//   - 含空白：每個字各自獨立遮罩
//
// Example:
//
//	(&NameMasker{mask: "*"}).Mask("John")     // returns "J**n"
//	(&NameMasker{mask: "*"}).Mask("John Doe") // returns "J**n D**e"
func (m *NameMasker) Mask(value string) string {
	l := len([]rune(value))

	if l == 0 {
		return ""
	}

	if strs := strings.Split(value, " "); len(strs) > 1 {
		tmp := make([]string, len(strs))
		for idx, str := range strs {
			tmp[idx] = m.Mask(str)
		}
		return strings.Join(tmp, " ")
	}

	if l == 2 || l == 3 {
		return overlay(value, strLoop(m.mask, 2), 1, 2)
	}

	if l > 3 {
		return overlay(value, strLoop(m.mask, 2), 1, 3)
	}

	return strLoop(m.mask, 2)
}
