package zapfield

import (
	"testing"

	masker "github.com/ggwhite/go-masker/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// helper 委派的 core 函式對照表，逐一驗證輸出與 core 逐字一致且型別為 String。
func TestFieldHelpers_DelegateToCore(t *testing.T) {
	const key = "k"
	tests := []struct {
		name  string
		fn    func(string, string) zap.Field
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
			want := zap.String(key, tt.core(tt.value))
			if got != want {
				t.Errorf("%s() = %+v, want %+v", tt.name, got, want)
			}
			if got.Type != zapcore.StringType {
				t.Errorf("%s() field type = %v, want StringType（不得用 Any／Reflect）", tt.name, got.Type)
			}
			if got.Key != key {
				t.Errorf("%s() field key = %q, want %q", tt.name, got.Key, key)
			}
		})
	}
}

// 空字串輸入的退化行為應與 core 一致，不得 panic。
func TestFieldHelpers_EmptyValue(t *testing.T) {
	if got := Phone("k", ""); got.String != masker.Mobile("") {
		t.Errorf("Phone(empty) = %q, want %q", got.String, masker.Mobile(""))
	}
}
