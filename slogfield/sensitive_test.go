package slogfield

import (
	"log/slog"
	"testing"

	masker "github.com/ggwhite/go-masker/v3"
)

// SensitiveAttr 應只讀快取 masked 值，輸出與 String() 逐字一致且型別為 String，且不等於原值。
func TestSensitiveAttr_UsesMaskedValue(t *testing.T) {
	s := masker.NewPhone("0987654321")
	got := SensitiveAttr("phone", s)
	if got.Key != "phone" {
		t.Errorf("SensitiveAttr() key = %q, want %q", got.Key, "phone")
	}
	if got.Value.Kind() != slog.KindString {
		t.Errorf("SensitiveAttr() value kind = %v, want KindString", got.Value.Kind())
	}
	if got.Value.String() != s.String() {
		t.Errorf("SensitiveAttr() = %q, want %q", got.Value.String(), s.String())
	}
	if got.Value.String() == s.Reveal() {
		t.Errorf("SensitiveAttr() 輸出等於原值，疑似洩漏 raw：%q", got.Value.String())
	}
}

// zero-value Sensitive[T]（mask 為 nil）的 String() 回空字串，adapter 應安全退化為空字串。
func TestSensitiveAttr_ZeroValueSafeDegradation(t *testing.T) {
	var s masker.Sensitive[string]
	got := SensitiveAttr("phone", s)
	if got.Value.String() != "" {
		t.Errorf("SensitiveAttr(zero) = %q, want \"\"（安全退化，不得洩漏 raw）", got.Value.String())
	}
}

// 非 string 型別的 T 也應正確運作（驗證泛型）。
func TestSensitiveAttr_GenericType(t *testing.T) {
	s := masker.NewSensitive(12345, func(int) string { return "*****" })
	got := SensitiveAttr("num", s)
	if got.Value.String() != "*****" {
		t.Errorf("SensitiveAttr(int) = %q, want %q", got.Value.String(), "*****")
	}
}
