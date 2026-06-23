package zapfield

import (
	"testing"

	masker "github.com/ggwhite/go-masker/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Sensitive adapter 應只讀快取 masked 值，輸出與 String() 逐字一致且型別為 String。
func TestSensitive_UsesMaskedValue(t *testing.T) {
	s := masker.NewPhone("0987654321")
	got := Sensitive("phone", s)
	want := zap.String("phone", s.String())
	if got != want {
		t.Errorf("Sensitive() = %+v, want %+v", got, want)
	}
	if got.Type != zapcore.StringType {
		t.Errorf("Sensitive() field type = %v, want StringType", got.Type)
	}
	if got.String == s.Reveal() {
		t.Errorf("Sensitive() 輸出等於原值，疑似洩漏 raw：%q", got.String)
	}
}

// zero-value Sensitive[T]（mask 為 nil）的 String() 回空字串，adapter 應安全退化為空字串。
func TestSensitive_ZeroValueSafeDegradation(t *testing.T) {
	var s masker.Sensitive[string]
	got := Sensitive("phone", s)
	if got.String != "" {
		t.Errorf("Sensitive(zero) = %q, want \"\"（安全退化，不得洩漏 raw）", got.String)
	}
}

// 非 string 型別的 T 也應正確運作（驗證泛型）。
func TestSensitive_GenericType(t *testing.T) {
	s := masker.NewSensitive(12345, func(int) string { return "*****" })
	got := Sensitive("num", s)
	if got.String != "*****" {
		t.Errorf("Sensitive(int) = %q, want %q", got.String, "*****")
	}
}
