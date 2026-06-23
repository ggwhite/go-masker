package slogfield

import (
	"log/slog"
	"testing"

	masker "github.com/ggwhite/go-masker/v3"
)

// helper 委派的 core 函式對照表，逐一驗證輸出與 core 逐字一致且型別為 String。
func TestFieldHelpers_DelegateToCore(t *testing.T) {
	const key = "k"
	tests := []struct {
		name  string
		fn    func(string, string) slog.Attr
		core  func(string) string
		value string
	}{
		{"Phone", Phone, masker.Mobile, "0987654321"},
		{"Email", Email, masker.Email, "ggw.chang@gmail.com"},
		{"Password", Password, masker.Password, "password"},
		{"Name", Name, masker.Name, "John Doe"},
		{"Address", Address, masker.Address, "台北市內湖區內湖路一段737巷1號1樓"},
		{"ID", ID, masker.ID, "A123456789"},
		{"Credit", Credit, masker.Credit, "4111111111111111"},
		{"Tel", Tel, masker.Tel, "0227993078"},
		{"URL", URL, masker.URL, "http://john:password@localhost:3000"},
		{"Abuse", Abuse, masker.Abuse, "hello"},
		{"None", None, masker.None, "secret"},
		{"All", All, masker.All, "secret"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(key, tt.value)
			want := tt.core(tt.value)
			if got.Key != key {
				t.Errorf("%s() key = %q, want %q", tt.name, got.Key, key)
			}
			if got.Value.Kind() != slog.KindString {
				t.Errorf("%s() value kind = %v, want KindString（不得用 slog.Any／reflection）", tt.name, got.Value.Kind())
			}
			if got.Value.String() != want {
				t.Errorf("%s() value = %q, want %q", tt.name, got.Value.String(), want)
			}
		})
	}
}

// 空字串輸入的退化行為應與 core 一致，不得 panic。
func TestFieldHelpers_EmptyValue(t *testing.T) {
	got := Phone("k", "")
	if got.Value.String() != masker.Mobile("") {
		t.Errorf("Phone(empty) = %q, want %q", got.Value.String(), masker.Mobile(""))
	}
}

// Masked 以動態 MaskerType 遮罩，輸出與 masker.Mask 逐字一致。
func TestMasked_DelegatesToMask(t *testing.T) {
	got := Masked("k", masker.TypeMobile, "0987654321")
	want := masker.Mask(masker.TypeMobile, "0987654321")
	if got.Key != "k" {
		t.Errorf("Masked() key = %q, want %q", got.Key, "k")
	}
	if got.Value.Kind() != slog.KindString {
		t.Errorf("Masked() value kind = %v, want KindString", got.Value.Kind())
	}
	if got.Value.String() != want {
		t.Errorf("Masked() value = %q, want %q", got.Value.String(), want)
	}
}

// 未知 MaskerType 時 masker.Mask 回傳原值（寬鬆入口），slogfield 不另外處理。
func TestMasked_UnknownTypeReturnsRaw(t *testing.T) {
	got := Masked("k", masker.MaskerType("nonexistent"), "value")
	if got.Value.String() != "value" {
		t.Errorf("Masked(unknown) = %q, want %q（未知型別回原值）", got.Value.String(), "value")
	}
}
