package masker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

func TestSensitiveReveal(t *testing.T) {
	s := NewPhone("0987654321")
	if got := s.Reveal(); got != "0987654321" {
		t.Fatalf("Reveal() = %q, want %q", got, "0987654321")
	}
}

func TestSensitiveString(t *testing.T) {
	s := NewPhone("0987654321")
	got := fmt.Sprint(s)
	if got != "0987***321" {
		t.Fatalf("Sprint = %q, want %q", got, "0987***321")
	}
	if strings.Contains(got, "654") {
		t.Fatalf("Sprint = %q leaked raw substring", got)
	}
}

func TestSensitiveGoString(t *testing.T) {
	s := NewEmail("ggw.chang@gmail.com")
	got := fmt.Sprintf("%#v", s)
	if got != Email("ggw.chang@gmail.com") {
		t.Fatalf("%%#v = %q, want %q", got, Email("ggw.chang@gmail.com"))
	}
	if strings.Contains(got, "chang") {
		t.Fatalf("%%#v = %q leaked raw account", got)
	}
}

type smsRequest struct {
	Phone   Sensitive[string] `json:"phone"`
	Message string            `json:"message"`
}

func TestSensitiveMarshalJSON(t *testing.T) {
	req := smsRequest{Phone: NewPhone("0987654321"), Message: "驗證碼 1234"}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	want := `{"phone":"0987***321","message":"驗證碼 1234"}`
	if string(b) != want {
		t.Fatalf("Marshal = %s, want %s", b, want)
	}
}

func TestSensitiveMarshalText(t *testing.T) {
	b, err := NewID("A123456789").MarshalText()
	if err != nil {
		t.Fatalf("MarshalText error: %v", err)
	}
	if string(b) != "A12345****" {
		t.Fatalf("MarshalText = %q, want %q", b, "A12345****")
	}
}

func TestSensitiveLogValue(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	logger.Info("login", slog.Any("phone", NewPhone("0987654321")))
	out := buf.String()
	if !strings.Contains(out, "0987***321") {
		t.Fatalf("slog output %q missing masked value", out)
	}
	if strings.Contains(out, "0987654321") {
		t.Fatalf("slog output %q leaked raw value", out)
	}
}

func TestNewSensitiveCustomType(t *testing.T) {
	mask := func(n int) string { return strings.Repeat("#", len(fmt.Sprint(n))) }
	s := NewSensitive(12345, mask)
	if s.String() != "#####" {
		t.Fatalf("String() = %q, want %q", s.String(), "#####")
	}
	if s.Reveal() != 12345 {
		t.Fatalf("Reveal() = %d, want %d", s.Reveal(), 12345)
	}
}

func TestSensitiveEqual(t *testing.T) {
	if !NewPhone("0987654321").Equal(NewPhone("0987654321")) {
		t.Fatal("same raw should be Equal")
	}
	if NewPhone("0987654321").Equal(NewPhone("0911222333")) {
		t.Fatal("different raw should not be Equal")
	}
	// 不同 raw 但 masked 相同（Password 固定 14 顆星）仍應回 false。
	a := NewPassword("aaaaaaaa")
	b := NewPassword("bbbbbbbb")
	if a.String() != b.String() {
		t.Fatalf("precondition failed: masked differ %q %q", a.String(), b.String())
	}
	if a.Equal(b) {
		t.Fatal("same masked but different raw should not be Equal")
	}
}

func TestSensitiveUnmarshalJSON(t *testing.T) {
	req := smsRequest{Phone: NewPhone("")}
	if err := json.Unmarshal([]byte(`{"phone":"0987654321"}`), &req); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if req.Phone.Reveal() != "0987654321" {
		t.Fatalf("Reveal() = %q, want %q", req.Phone.Reveal(), "0987654321")
	}
	if req.Phone.String() != "0987***321" {
		t.Fatalf("String() = %q, want %q", req.Phone.String(), "0987***321")
	}
}

func TestSensitiveNilMaskConstructor(t *testing.T) {
	s := NewSensitive("secret", (func(string) string)(nil))
	if s.String() != "" {
		t.Fatalf("String() = %q, want empty", s.String())
	}
	if s.Reveal() != "secret" {
		t.Fatalf("Reveal() = %q, want %q", s.Reveal(), "secret")
	}
}

func TestSensitiveUnmarshalNilMask(t *testing.T) {
	var s Sensitive[string]
	if err := json.Unmarshal([]byte(`"secret"`), &s); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if s.String() != "" {
		t.Fatalf("String() = %q, want empty (no leak)", s.String())
	}
	if s.Reveal() != "secret" {
		t.Fatalf("Reveal() = %q, want %q", s.Reveal(), "secret")
	}
}

func ExampleSensitive_jSON() {
	req := smsRequest{Phone: NewPhone("0987654321"), Message: "hi"}
	b, _ := json.Marshal(req)
	fmt.Println(string(b))
	// Output: {"phone":"0987***321","message":"hi"}
}
