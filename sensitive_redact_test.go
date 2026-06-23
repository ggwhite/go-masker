package masker

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestRedactAllOutputs 確認 redact 模式下所有自動輸出介面都回傳替換文字，不洩漏任何部分資訊。
func TestRedactAllOutputs(t *testing.T) {
	s := NewPhone("0987654321", WithRedact())

	if got := s.String(); got != DefaultRedactText {
		t.Fatalf("String() = %q, want %q", got, DefaultRedactText)
	}
	if got := fmt.Sprintf("%#v", s); got != DefaultRedactText {
		t.Fatalf("GoString() = %q, want %q", got, DefaultRedactText)
	}
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	if got, want := string(jsonBytes), `"[REDACTED]"`; got != want {
		t.Fatalf("MarshalJSON = %s, want %s", got, want)
	}
	text, err := s.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText error: %v", err)
	}
	if got := string(text); got != DefaultRedactText {
		t.Fatalf("MarshalText = %q, want %q", got, DefaultRedactText)
	}
	if got := s.LogValue().String(); got != DefaultRedactText {
		t.Fatalf("LogValue = %q, want %q", got, DefaultRedactText)
	}
}

// TestRedactRevealUnaffected 確認 redact 不影響 Reveal，原值仍可取回。
func TestRedactRevealUnaffected(t *testing.T) {
	const raw = "0987654321"
	if got := NewPhone(raw, WithRedact()).Reveal(); got != raw {
		t.Fatalf("Reveal() = %q, want %q", got, raw)
	}
	if got := NewPhone(raw, WithRedactText("***")).Reveal(); got != raw {
		t.Fatalf("Reveal() = %q, want %q", got, raw)
	}
}

// TestRedactValueUnaffected 確認 redact 不影響 SQL Value，寫入 DB 仍為原值明文。
func TestRedactValueUnaffected(t *testing.T) {
	const raw = "0987654321"
	v, err := NewPhone(raw, WithRedact()).Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if v != raw {
		t.Fatalf("Value() = %v, want %q", v, raw)
	}
}

// TestWithRedactText 確認自訂替換文字生效，且 WithRedactText 隱含啟用 redact。
func TestWithRedactText(t *testing.T) {
	s := NewPhone("0987654321", WithRedactText("***"))
	if got := s.String(); got != "***" {
		t.Fatalf("String() = %q, want %q", got, "***")
	}
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	if got, want := string(jsonBytes), `"***"`; got != want {
		t.Fatalf("MarshalJSON = %s, want %s", got, want)
	}
}

// TestNonRedactStillMasks 確認未開啟 redact 時維持原本部分遮罩行為（向後相容）。
func TestNonRedactStillMasks(t *testing.T) {
	if got := NewPhone("0987654321").String(); got != "0987***321" {
		t.Fatalf("String() = %q, want %q", got, "0987***321")
	}
}

// TestRedactSurvivesUnmarshalJSON 確認 redact 狀態在 UnmarshalJSON 重算遮罩後仍保留。
func TestRedactSurvivesUnmarshalJSON(t *testing.T) {
	s := NewPhone("", WithRedact())
	if err := json.Unmarshal([]byte(`"0912345678"`), &s); err != nil {
		t.Fatalf("UnmarshalJSON error: %v", err)
	}
	if got := s.String(); got != DefaultRedactText {
		t.Fatalf("String() after unmarshal = %q, want %q", got, DefaultRedactText)
	}
	if got := s.Reveal(); got != "0912345678" {
		t.Fatalf("Reveal() after unmarshal = %q, want %q", got, "0912345678")
	}
}

// TestRedactSurvivesScan 確認 redact 狀態在 SQL Scan 重算遮罩後仍保留。
func TestRedactSurvivesScan(t *testing.T) {
	s := NewPhone("", WithRedactText("***"))
	if err := s.Scan("0912345678"); err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if got := s.String(); got != "***" {
		t.Fatalf("String() after scan = %q, want %q", got, "***")
	}
	if got := s.Reveal(); got != "0912345678" {
		t.Fatalf("Reveal() after scan = %q, want %q", got, "0912345678")
	}
}

func ExampleWithRedact() {
	s := NewPhone("0987654321", WithRedact())
	fmt.Println(s)
	fmt.Println(s.Reveal())
	// Output:
	// [REDACTED]
	// 0987654321
}

func ExampleWithRedactText() {
	fmt.Println(NewPhone("0987654321", WithRedactText("***")))
	// Output: ***
}
