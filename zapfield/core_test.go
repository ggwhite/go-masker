package zapfield

import (
	"regexp"
	"testing"

	masker "github.com/ggwhite/go-masker/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// newObservedLogger 以 observer core 包 WrapCore，回傳 logger 與可觀測 logs。
func newObservedLogger(rules InterceptRules) (*zap.Logger, *observer.ObservedLogs) {
	obsCore, logs := observer.New(zapcore.InfoLevel)
	return zap.New(WrapCore(obsCore, rules)), logs
}

// fieldByKey 從 entry 取出指定 key 的 field。
func fieldByKey(t *testing.T, entry observer.LoggedEntry, key string) zapcore.Field {
	t.Helper()
	for _, f := range entry.Context {
		if f.Key == key {
			return f
		}
	}
	t.Fatalf("找不到 key %q 的 field", key)
	return zapcore.Field{}
}

func TestWrapCore_KeywordMatchMasksString(t *testing.T) {
	logger, logs := newObservedLogger(InterceptRules{Keywords: []string{"phone"}})
	logger.Info("login", zap.String("phone", "0987654321"))

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "phone").String; got != "**********" {
		t.Errorf("phone 應被全遮罩，got %q", got)
	}
}

func TestWrapCore_PatternMatchMasksString(t *testing.T) {
	rules := InterceptRules{Patterns: []*regexp.Regexp{regexp.MustCompile(`(?i)secret`)}}
	logger, logs := newObservedLogger(rules)
	logger.Info("auth", zap.String("api_secret", "abcdef"))

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "api_secret").String; got != "******" {
		t.Errorf("api_secret 應被全遮罩，got %q", got)
	}
}

func TestWrapCore_NonStringFieldUntouched(t *testing.T) {
	logger, logs := newObservedLogger(InterceptRules{Keywords: []string{"phone"}})
	logger.Info("stat", zap.Int("phone_count", 3))

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "phone_count").Integer; got != 3 {
		t.Errorf("非 string field 不應被改動，got %d", got)
	}
}

func TestWrapCore_UnmatchedStringUntouched(t *testing.T) {
	logger, logs := newObservedLogger(InterceptRules{Keywords: []string{"phone"}})
	logger.Info("profile", zap.String("nickname", "johnny"))

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "nickname").String; got != "johnny" {
		t.Errorf("未命中 field 應原樣保留，got %q", got)
	}
}

func TestWrapCore_EmptyRulesPassThrough(t *testing.T) {
	logger, logs := newObservedLogger(InterceptRules{})
	logger.Info("pass", zap.String("phone", "0987654321"), zap.String("token", "abc"))

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "phone").String; got != "0987654321" {
		t.Errorf("空 rules 應原樣放行 phone，got %q", got)
	}
	if got := fieldByKey(t, entry, "token").String; got != "abc" {
		t.Errorf("空 rules 應原樣放行 token，got %q", got)
	}
}

func TestWrapCore_KeywordCaseInsensitive(t *testing.T) {
	logger, logs := newObservedLogger(InterceptRules{Keywords: []string{"phone"}})
	logger.Info("ci",
		zap.String("Phone", "0987654321"),
		zap.String("userPhone", "0911222333"),
	)

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "Phone").String; got != "**********" {
		t.Errorf("Phone 應命中（case-insensitive），got %q", got)
	}
	if got := fieldByKey(t, entry, "userPhone").String; got != "**********" {
		t.Errorf("userPhone 應命中（子字串），got %q", got)
	}
}

func TestWrapCore_WithContextFieldIntercepted(t *testing.T) {
	logger, logs := newObservedLogger(InterceptRules{Keywords: []string{"password"}})
	logger.With(zap.String("password", "hunter2")).Info("ctx")

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "password").String; got != "*******" {
		t.Errorf("With 帶入的 context field 應被攔截，got %q", got)
	}
}

func TestWrapCore_RuleAppliesSpecifiedMaskerType(t *testing.T) {
	rules := InterceptRules{
		Rules: []Rule{
			{Keywords: []string{"phone", "mobile"}, MaskerType: masker.TypeMobile},
			{Keywords: []string{"password", "secret"}, MaskerType: masker.TypePassword},
		},
	}
	logger, logs := newObservedLogger(rules)
	logger.Info("login",
		zap.String("phone", "0987654321"),
		zap.String("password", "hunter2"),
	)

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "phone").String; got != "0987***321" {
		t.Errorf("phone 應套用 TypeMobile，got %q", got)
	}
	if got := fieldByKey(t, entry, "password").String; got != "**************" {
		t.Errorf("password 應套用 TypePassword（固定 14 字元），got %q", got)
	}
}

func TestWrapCore_RuleZeroMaskerTypeFallsBackToAll(t *testing.T) {
	rules := InterceptRules{
		Rules: []Rule{
			{Keywords: []string{"token"}},
		},
	}
	logger, logs := newObservedLogger(rules)
	logger.Info("auth", zap.String("token", "abcdef"))

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "token").String; got != "******" {
		t.Errorf("MaskerType zero value 應退化為全遮罩，got %q", got)
	}
}

func TestWrapCore_RuleKeywordCaseInsensitive(t *testing.T) {
	rules := InterceptRules{
		Rules: []Rule{
			{Keywords: []string{"phone"}, MaskerType: masker.TypeMobile},
		},
	}
	logger, logs := newObservedLogger(rules)
	logger.Info("ci", zap.String("UserPhone", "0987654321"))

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "UserPhone").String; got != "0987***321" {
		t.Errorf("Rule 的 Keywords 應 case-insensitive 命中，got %q", got)
	}
}

func TestWrapCore_RulePatternAppliesMaskerType(t *testing.T) {
	rules := InterceptRules{
		Rules: []Rule{
			{Patterns: []*regexp.Regexp{regexp.MustCompile(`(?i)mobile`)}, MaskerType: masker.TypeMobile},
		},
	}
	logger, logs := newObservedLogger(rules)
	logger.Info("auth", zap.String("user_mobile", "0987654321"))

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "user_mobile").String; got != "0987***321" {
		t.Errorf("Rule 的 Patterns 命中應套用 TypeMobile，got %q", got)
	}
}

func TestWrapCore_TopLevelTakesPrecedenceOverRules(t *testing.T) {
	rules := InterceptRules{
		Keywords: []string{"phone"},
		Rules: []Rule{
			{Keywords: []string{"phone"}, MaskerType: masker.TypeMobile},
		},
	}
	logger, logs := newObservedLogger(rules)
	logger.Info("login", zap.String("phone", "0987654321"))

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "phone").String; got != "**********" {
		t.Errorf("頂層 Keywords 應優先於 Rules（全遮罩），got %q", got)
	}
}

func TestWrapCore_RuleWithContextFieldIntercepted(t *testing.T) {
	rules := InterceptRules{
		Rules: []Rule{
			{Keywords: []string{"phone"}, MaskerType: masker.TypeMobile},
		},
	}
	logger, logs := newObservedLogger(rules)
	logger.With(zap.String("phone", "0987654321")).Info("ctx")

	entry := logs.All()[0]
	if got := fieldByKey(t, entry, "phone").String; got != "0987***321" {
		t.Errorf("With 帶入的 context field 應依 Rule 套用 TypeMobile，got %q", got)
	}
}
