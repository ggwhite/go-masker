package masker

import "net/url"

// URLMasker 是 URL 的 masker。
type URLMasker struct{}

// Mask 遮罩 URL 中的密碼段（若存在）。
// 密碼段固定輸出 xxxxx，由標準函式庫 url.Redacted() 寫死，不受 WithMaskChar 影響。
// Example:
//
//	(&URLMasker{}).Mask("http://john:password@localhost:3000") // returns "http://john:xxxxx@localhost:3000"
func (m *URLMasker) Mask(value string) string {
	u, err := url.Parse(value)
	if err != nil {
		return value
	}
	return u.Redacted()
}
