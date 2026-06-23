package zapfield

import (
	"regexp"
	"strings"

	masker "github.com/ggwhite/go-masker/v3"
	"go.uber.org/zap/zapcore"
)

// Rule 定義一條攔截規則：以 Keywords／Patterns 比對 field key，命中時套用 MaskerType
// 指定的遮罩型別。讓不同欄位能用最合適的遮罩（如 phone 用 TypeMobile 保留前後碼），
// 而非一律全遮罩。
//
// MaskerType 為 zero value（空字串）時退化為 TypeAll（全遮罩），維持與舊版相同行為。
type Rule struct {
	// Keywords 對 field key 做 case-insensitive 子字串比對（strings.Contains）。
	Keywords []string
	// Patterns 對 field key 做 regex 比對；需於建立 Rule 時自行編譯，
	// 以保證 WrapCore 不需回傳 error。
	Patterns []*regexp.Regexp
	// MaskerType 指定命中時套用的遮罩型別；zero value 視為 masker.TypeAll。
	MaskerType masker.MaskerType
}

// InterceptRules 定義 WrapCore 攔截時的 field key 比對規則。
//
// 比對僅針對 field 的 key 名稱，分兩種：Keywords 以小寫化後的子字串比對、
// Patterns 以已編譯的 regex 比對；任一命中即視為需遮罩。
//
// 兩種設定來源：
//   - 頂層 Keywords／Patterns：命中即全遮罩（TypeAll），為最早的簡易用法。
//   - Rules：每條規則可指定自己的 MaskerType，命中時套用對應遮罩。
//
// 比對順序為「先頂層、後 Rules」，回傳首個命中者的遮罩型別。
//
// ⚠️ 此規則本質是「靠 key 名稱猜測敏感欄位」，必然不精準：
//   - false positive：key 名碰巧含關鍵字（如 phone_count）會被誤遮。
//   - false negative：敏感欄位的 key 名不在清單內則漏遮。
//
// 因此這只適合當作改不動業務 code 時的被動防線。若業務端已改用 Sensitive[T]
// 或 zapfield.Phone() 等顯式遮罩 helper，本層即為多餘，不應視為主力方案。
type InterceptRules struct {
	// Keywords 對 field key 做 case-insensitive 子字串比對（strings.Contains），命中即全遮罩。
	Keywords []string
	// Patterns 對 field key 做 regex 比對；需於建立 InterceptRules 時自行編譯，
	// 以保證 WrapCore 不需回傳 error。命中即全遮罩。
	Patterns []*regexp.Regexp
	// Rules 為每條可指定 MaskerType 的攔截規則，於頂層 Keywords／Patterns 未命中時依序比對。
	Rules []Rule
}

// maskerFor 回傳 key 應套用的遮罩型別與是否命中。
// 先比對頂層 Keywords/Patterns（命中即 TypeAll，維持既有行為），
// 再依序比對 Rules，回傳首條命中規則的 MaskerType（zero value 視為 TypeAll）。
// 全部未命中時回傳 false，退化為「不攔截、原樣放行」，此為最安全且不 panic 的預設行為。
func (r InterceptRules) maskerFor(key string) (masker.MaskerType, bool) {
	if matchKey(key, r.Keywords, r.Patterns) {
		return masker.TypeAll, true
	}
	for _, rule := range r.Rules {
		if matchKey(key, rule.Keywords, rule.Patterns) {
			t := rule.MaskerType
			if t == "" {
				t = masker.TypeAll
			}
			return t, true
		}
	}
	return "", false
}

// matchKey 回報 key 是否命中任一 keyword（小寫子字串）或 pattern（regex）。
func matchKey(key string, keywords []string, patterns []*regexp.Regexp) bool {
	if len(keywords) > 0 {
		lower := strings.ToLower(key)
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				return true
			}
		}
	}
	for _, p := range patterns {
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
//
// 進階用法：以 Rules 對不同欄位指定遮罩型別，而非一律全遮罩。
//
//	core := zapfield.WrapCore(originalCore, zapfield.InterceptRules{
//	    Rules: []zapfield.Rule{
//	        {Keywords: []string{"phone", "mobile"}, MaskerType: masker.TypeMobile},
//	        {Keywords: []string{"password", "secret"}, MaskerType: masker.TypePassword},
//	    },
//	})
func WrapCore(core zapcore.Core, rules InterceptRules) zapcore.Core {
	rules.Keywords = lowerKeywords(rules.Keywords)
	if len(rules.Rules) > 0 {
		normalizedRules := make([]Rule, len(rules.Rules))
		for i, rule := range rules.Rules {
			rule.Keywords = lowerKeywords(rule.Keywords)
			normalizedRules[i] = rule
		}
		rules.Rules = normalizedRules
	}
	return &maskingCore{Core: core, rules: rules}
}

// lowerKeywords 回傳一份全小寫化的 keywords copy；輸入為空時原樣回傳，
// 避免就地改寫呼叫端傳入的 slice。
func lowerKeywords(in []string) []string {
	if len(in) == 0 {
		return in
	}
	out := make([]string, len(in))
	for i, kw := range in {
		out[i] = strings.ToLower(kw)
	}
	return out
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
// 套用對應遮罩型別，其餘原樣保留。為避免就地改寫呼叫端傳入的 slice，一律建立 copy。
func (c *maskingCore) interceptFields(fields []zapcore.Field) []zapcore.Field {
	out := make([]zapcore.Field, len(fields))
	copy(out, fields)
	for i := range out {
		if out[i].Type != zapcore.StringType {
			continue
		}
		if t, ok := c.rules.maskerFor(out[i].Key); ok {
			out[i].String = masker.Mask(t, out[i].String)
		}
	}
	return out
}
