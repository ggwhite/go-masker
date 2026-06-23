package masker

import (
	"database/sql/driver"
	"testing"
	"time"
)

// TestSensitiveValueReturnsRaw 確認 Value() 回傳原值（明文），不是遮罩值。
func TestSensitiveValueReturnsRaw(t *testing.T) {
	s := NewPhone("0987654321")
	v, err := s.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if v != "0987654321" {
		t.Fatalf("Value() = %v, want raw %q", v, "0987654321")
	}
}

// TestSensitiveValueIntConversion 確認非 string 的 raw 會被轉成合法 driver.Value（int→int64）。
func TestSensitiveValueIntConversion(t *testing.T) {
	mask := func(int) string { return "***" }
	s := NewSensitive(42, mask)
	v, err := s.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	got, ok := v.(int64)
	if !ok {
		t.Fatalf("Value() type = %T, want int64", v)
	}
	if got != 42 {
		t.Fatalf("Value() = %d, want 42", got)
	}
}

// TestSensitiveScanString 確認從 string／[]byte scan 後重算遮罩值並保留原值。
func TestSensitiveScanString(t *testing.T) {
	cases := []struct {
		name string
		src  any
	}{
		{"string", "0987654321"},
		{"bytes", []byte("0987654321")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := NewPhone("") // 預先綁定 Mobile mask
			if err := s.Scan(c.src); err != nil {
				t.Fatalf("Scan(%v) error: %v", c.src, err)
			}
			if s.Reveal() != "0987654321" {
				t.Fatalf("Reveal() = %q, want %q", s.Reveal(), "0987654321")
			}
			if s.String() != "0987***321" {
				t.Fatalf("String() = %q, want %q", s.String(), "0987***321")
			}
		})
	}
}

// TestSensitiveScanNil 確認 Scan(nil)（DB NULL）設為 zero value 且不 panic。
func TestSensitiveScanNil(t *testing.T) {
	s := NewPhone("0987654321")
	if err := s.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) error: %v", err)
	}
	if s.Reveal() != "" {
		t.Fatalf("Reveal() = %q, want zero value", s.Reveal())
	}
	// mask("") 仍應被重算，不是殘留舊遮罩值。
	if s.String() != Mobile("") {
		t.Fatalf("String() = %q, want %q", s.String(), Mobile(""))
	}
}

// TestSensitiveScanInt 確認整數型別轉換（int64 來源 + 字串數字）。
func TestSensitiveScanInt(t *testing.T) {
	mask := func(n int64) string { return "###" }
	s := NewSensitive(int64(0), mask)
	if err := s.Scan(int64(123456)); err != nil {
		t.Fatalf("Scan(int64) error: %v", err)
	}
	if s.Reveal() != 123456 {
		t.Fatalf("Reveal() = %d, want 123456", s.Reveal())
	}
	if s.String() != "###" {
		t.Fatalf("String() = %q, want %q", s.String(), "###")
	}

	// 驅動回傳字串數字（如 []byte("789")）也應可轉。
	if err := s.Scan([]byte("789")); err != nil {
		t.Fatalf("Scan([]byte number) error: %v", err)
	}
	if s.Reveal() != 789 {
		t.Fatalf("Reveal() = %d, want 789", s.Reveal())
	}
}

// TestSensitiveScanFloat 確認浮點型別轉換（float64 與字串）。
func TestSensitiveScanFloat(t *testing.T) {
	mask := func(float64) string { return "$$$" }
	s := NewSensitive(0.0, mask)
	if err := s.Scan(3.14); err != nil {
		t.Fatalf("Scan(float64) error: %v", err)
	}
	if s.Reveal() != 3.14 {
		t.Fatalf("Reveal() = %v, want 3.14", s.Reveal())
	}
	if err := s.Scan([]byte("2.5")); err != nil {
		t.Fatalf("Scan([]byte float) error: %v", err)
	}
	if s.Reveal() != 2.5 {
		t.Fatalf("Reveal() = %v, want 2.5", s.Reveal())
	}
}

// TestSensitiveScanBool 確認 bool 型別轉換（bool、0/1、字串）。
func TestSensitiveScanBool(t *testing.T) {
	mask := func(bool) string { return "?" }
	s := NewSensitive(false, mask)
	if err := s.Scan(true); err != nil {
		t.Fatalf("Scan(bool) error: %v", err)
	}
	if s.Reveal() != true {
		t.Fatalf("Reveal() = %v, want true", s.Reveal())
	}
	if err := s.Scan(int64(0)); err != nil {
		t.Fatalf("Scan(int64 bool) error: %v", err)
	}
	if s.Reveal() != false {
		t.Fatalf("Reveal() = %v, want false", s.Reveal())
	}
}

// TestSensitiveScanTime 確認 time.Time 這類可直接賦值的型別。
func TestSensitiveScanTime(t *testing.T) {
	mask := func(time.Time) string { return "redacted" }
	s := NewSensitive(time.Time{}, mask)
	now := time.Date(2026, 6, 23, 12, 0, 0, 0, time.UTC)
	if err := s.Scan(now); err != nil {
		t.Fatalf("Scan(time.Time) error: %v", err)
	}
	if !s.Reveal().Equal(now) {
		t.Fatalf("Reveal() = %v, want %v", s.Reveal(), now)
	}
}

// TestSensitiveScanNilMaskNoLeak 確認未綁定 mask 的 zero-value receiver（ORM 場景）
// scan 後安全儲存原值但遮罩值維持空字串，不洩漏原值、不報錯。
func TestSensitiveScanNilMaskNoLeak(t *testing.T) {
	var s Sensitive[string]
	if err := s.Scan("0987654321"); err != nil {
		t.Fatalf("Scan with nil mask error: %v", err)
	}
	if s.Reveal() != "0987654321" {
		t.Fatalf("Reveal() = %q, want %q", s.Reveal(), "0987654321")
	}
	if s.String() != "" {
		t.Fatalf("String() = %q, want empty (no leak without bound mask)", s.String())
	}
}

// TestSensitiveScanTypeMismatch 確認無法轉換的型別回傳 error，不 panic。
func TestSensitiveScanTypeMismatch(t *testing.T) {
	s := NewPhone("")
	if err := s.Scan(struct{ X int }{X: 1}); err == nil {
		t.Fatal("expected error scanning incompatible type, got nil")
	}
}

// TestSensitiveScanOverflow 確認超出目的整數範圍時回傳 error。
func TestSensitiveScanOverflow(t *testing.T) {
	mask := func(int8) string { return "x" }
	s := NewSensitive(int8(0), mask)
	if err := s.Scan(int64(300)); err == nil {
		t.Fatal("expected overflow error, got nil")
	}
}

// TestSensitiveScanValueRoundTrip 確認 scan 進來的值可透過 Value 原樣寫回（明文 round-trip）。
func TestSensitiveScanValueRoundTrip(t *testing.T) {
	s := NewEmail("")
	if err := s.Scan([]byte("ggw.chang@gmail.com")); err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	v, err := s.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if v.(string) != "ggw.chang@gmail.com" {
		t.Fatalf("round-trip Value() = %v, want raw email", v)
	}
	// 顯示路徑仍是遮罩值。
	if s.String() != "ggw****@gmail.com" {
		t.Fatalf("String() = %q, want masked", s.String())
	}
}

// 確認 Value 回傳的型別都是合法 driver.Value。
func TestSensitiveValueIsDriverValue(t *testing.T) {
	s := NewPhone("0987654321")
	v, err := s.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if !driver.IsValue(v) {
		t.Fatalf("Value() = %v (%T) is not a valid driver.Value", v, v)
	}
}
