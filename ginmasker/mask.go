package ginmasker

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"

	masker "github.com/ggwhite/go-masker/v3"
)

// redactionRepeat 為非 string 值的固定 redaction 長度，避免洩漏原值與型別細節。
const redactionRepeat = 4

// maskJSONBody 解析 raw JSON、遞迴遮罩命中 keyword 的欄位，回傳遮罩後的 JSON 字串。
// 解析失敗時回 ("", false)，呼叫端據此跳過遮罩、改記原始標記。
func maskJSONBody(raw []byte, keywords []string, maskChar string) (string, bool) {
	var node any
	if err := json.Unmarshal(raw, &node); err != nil {
		return "", false
	}
	masked := walkAndMask(node, keywords, maskChar)
	out, err := json.Marshal(masked)
	if err != nil {
		return "", false
	}
	return string(out), true
}

// walkAndMask 遞迴走訪 JSON 結構：object 的 key 命中 keyword 時整個 value 遮罩，
// 否則對 value 繼續遞迴；array 對每個元素遞迴；scalar 原樣回傳。
func walkAndMask(node any, keywords []string, maskChar string) any {
	switch v := node.(type) {
	case map[string]any:
		for key, val := range v {
			if keyMatches(key, keywords) {
				v[key] = maskValue(val, maskChar)
				continue
			}
			v[key] = walkAndMask(val, keywords, maskChar)
		}
		return v
	case []any:
		for i, val := range v {
			v[i] = walkAndMask(val, keywords, maskChar)
		}
		return v
	default:
		return v
	}
}

// keyMatches 回報 key 是否命中任一 keyword（雙方皆小寫化後做子字串比對），
// 對齊 zapfield 的 matchKey 語意。
func keyMatches(key string, keywords []string) bool {
	if len(keywords) == 0 {
		return false
	}
	lower := strings.ToLower(key)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// maskValue 全遮一個命中欄位的值：
// string 採全遮（見 maskString）；
// 非 string（number／bool／object／array）一律回傳固定長度 redaction 字串，
// 不洩漏原值與型別細節。
func maskValue(v any, maskChar string) any {
	if s, ok := v.(string); ok {
		return maskString(s, maskChar)
	}
	return strings.Repeat(maskChar, redactionRepeat)
}

// maskString 對字串做全遮：maskChar 為預設 "*" 時委派 masker.All 確保與 v3 core 逐字一致，
// 自訂 maskChar 時則依 rune 數重複該字元。
func maskString(s, maskChar string) string {
	if maskChar == defaultMaskChar {
		return masker.All(s)
	}
	return strings.Repeat(maskChar, len([]rune(s)))
}

// maskedQuery 回傳遮罩後的 query 字串：對 keys 命中的參數（精確比對）全遮其值，
// 其餘原樣保留。操作於 values 的副本上，不影響原 request。
func maskedQuery(values url.Values, keys []string, maskChar string) string {
	masked := make(url.Values, len(values))
	for k, vs := range values {
		masked[k] = append([]string(nil), vs...)
	}
	for _, key := range keys {
		vs, ok := masked[key]
		if !ok {
			continue
		}
		for i, v := range vs {
			vs[i] = maskString(v, maskChar)
		}
	}
	return masked.Encode()
}

// maskedHeaders 回傳遮罩後的 header 字串：對 keys 命中的 header（精確比對、case-insensitive）
// 全遮其值，其餘原樣保留。操作於 header 的副本上，不修改 c.Request.Header 本體。
// 輸出依 header 名稱排序，格式為 "Key=v1,v2"，以分號分隔，確保結果穩定可測。
func maskedHeaders(h http.Header, keys []string, maskChar string) string {
	cloned := h.Clone()
	if cloned == nil {
		cloned = http.Header{}
	}
	for _, key := range keys {
		ck := http.CanonicalHeaderKey(key)
		vs, ok := cloned[ck]
		if !ok {
			continue
		}
		for i, v := range vs {
			vs[i] = maskString(v, maskChar)
		}
	}
	names := make([]string, 0, len(cloned))
	for name := range cloned {
		names = append(names, name)
	}
	sort.Strings(names)
	var b strings.Builder
	for i, name := range names {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(name)
		b.WriteString("=")
		b.WriteString(strings.Join(cloned[name], ","))
	}
	return b.String()
}
