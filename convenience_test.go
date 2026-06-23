package masker

import "testing"

// TestConvenienceFunctions 驗證 package-level 便利函式輸出（AC-9）。
func TestConvenienceFunctions(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"Mobile", Mobile("0987654321"), "0987***321"},
		{"Email", Email("ggw.chang@gmail.com"), "ggw****@gmail.com"},
		{"Password", Password("secret"), "**************"},
		{"Name", Name("John Doe"), "J**n D**e"},
		{"Address", Address("台北市大安區敦化南路五段7788號378樓"), "台北市大安區******"},
		{"ID", ID("A123456789"), "A12345****"},
		{"Credit", Credit("4111111111111111"), "411111******1111"},
		{"Tel", Tel("0227993078"), "(02)2799-****"},
		{"URL", URL("http://john:password@localhost:3000"), "http://john:xxxxx@localhost:3000"},
		{"None", None("secret"), "secret"},
		{"All", All("secret"), "******"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s() = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

// TestMaskDynamic 驗證動態 Mask 函式（AC-10）：正常型別、未知型別回原值、generic tag。
func TestMaskDynamic(t *testing.T) {
	if got := Mask(TypeMobile, "0987654321"); got != "0987***321" {
		t.Errorf("Mask(TypeMobile) = %v, want 0987***321", got)
	}
	if got := Mask("nonexistent", "value"); got != "value" {
		t.Errorf("Mask(nonexistent) = %v, want value (原值)", got)
	}
	if got := Mask("first-3", "hello"); got != "***lo" {
		t.Errorf("Mask(first-3) = %v, want ***lo", got)
	}
}
