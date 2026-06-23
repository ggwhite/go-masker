package zapfield

import (
	"regexp"
	"strings"

	masker "github.com/ggwhite/go-masker/v3"
	"go.uber.org/zap/zapcore"
)

// InterceptRules 定義 WrapCore 攔截時的 field key 比對規則。
//
// 比對僅針對 field 的 key 名稱，分兩種：Keywords 以小寫化後的子字串比對、
// Patterns 以已編譯的 regex 比對；任一命中即視為需遮罩。
//
// ⚠️ 此規則本質是「靠 key 名稱猜測敏感欄位」，必然不精準：
//   - false positive：key 名碰巧含關鍵字（如 phone_count）會被誤遮。
//   - false negative：敏感欄位的 key 名不在清單內則漏遮。
//
// 因此這只適合當作改不動業務 code 時的被動防線。若業務端已改用 Sensitive[T]
// 或 zapfield.Phone() 等顯式遮罩 helper，本層即為多餘，不應視為主力方案。
type InterceptRules struct {
	// Keywords 對 field key 做 case-insensitive 子字串比對（strings.Contains）。
	Keywords []string
	// Patterns 對 field key 做 regex 比對；需於建立 InterceptRules 時自行編譯，
	// 以保證 WrapCore 不需回傳 error。
	Patterns []*regexp.Regexp
}

// match 回報 key 是否命中任一 keyword 或 pattern。
// Keywords 與 Patterns 皆空（或 nil）時一律回傳 false，退化為「不攔截、原樣放行」，
// 此為最安全且不 panic 的預設行為。
func (r InterceptRules) match(key string) bool {
	if len(r.Keywords) > 0 {
		lower := strings.ToLower(key)
		for _, kw := range r.Keywords {
			if strings.Contains(lower, kw) {
				return true
			}
		}
	}
	for _, p := range r.Patterns {
		if p != nil && p.MatchString(key) {
			return true
		}
	}
	return false
}

// maskingCore 包裝既有 zapcore.Core，在寫出前依 rules 攔截並遮罩 string field。
type maskingCore struct {
	zapcore.Core
	rules InterceptRules
}

// WrapCore 以 rules 包裝既有 zapcore.Core，於 log 寫出前對「key 命中規則」的
// string field 整串遮罩（每字元換成遮罩字元），作為敏感資料外洩的最後一道防線。
//
// 典型用法：對改不動的既有 logger 被動加上一層防護。
//
//	core := zapfield.WrapCore(originalCore, zapfield.InterceptRules{
//	    Keywords: []string{"phone", "password", "token", "secret"},
//	})
//	logger := zap.New(core)
//
// ⚠️ 攔截依 field key 名稱猜測，必然有誤遮（如 phone_count）與漏遮（key 不在清單）；
// 詳見 InterceptRules 說明。若業務 code 已採顯式遮罩（Sensitive[T] / zapfield.Phone 等），
// 本層多餘，不應視為主力方案。
//
// 僅處理頂層 string 型別 field；數字、布林、結構等非 string field 一律原樣保留，
// 巢狀 object / array 內部欄位不在此層攔截範圍。
func WrapCore(core zapcore.Core, rules InterceptRules) zapcore.Core {
	normalized := make([]string, len(rules.Keywords))
	for i, kw := range rules.Keywords {
		normalized[i] = strings.ToLower(kw)
	}
	rules.Keywords = normalized
	return &maskingCore{Core: core, rules: rules}
}

// With 先攔截累積的 context field，再委派內層 With，並以同樣 rules 重新包裝，
// 確保 logger.With(...) 帶入的 context field 同樣會被攔截。
func (c *maskingCore) With(fields []zapcore.Field) zapcore.Core {
	return &maskingCore{
		Core:  c.Core.With(c.interceptFields(fields)),
		rules: c.rules,
	}
}

// Check 依 Enabled 慣例把自己加入 CheckedEntry，為標準 wrapper 寫法。
func (c *maskingCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

// Write 先攔截 field，再委派內層 Write。
func (c *maskingCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	return c.Core.Write(ent, c.interceptFields(fields))
}

// interceptFields 回傳一份新 slice：對 key 命中 rules 且型別為 string 的 field
// 整串遮罩，其餘原樣保留。為避免就地改寫呼叫端傳入的 slice，一律建立 copy。
func (c *maskingCore) interceptFields(fields []zapcore.Field) []zapcore.Field {
	out := make([]zapcore.Field, len(fields))
	copy(out, fields)
	for i := range out {
		if out[i].Type == zapcore.StringType && c.rules.match(out[i].Key) {
			out[i].String = masker.All(out[i].String)
		}
	}
	return out
}
