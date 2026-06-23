package masker

import (
	"fmt"
	"testing"
)

// TestBuiltinConstructorsParity 確認九個建構子的 masked 輸出與對應 v3 package-level 函式逐字一致。
func TestBuiltinConstructorsParity(t *testing.T) {
	cases := []struct {
		name string
		ctor func(string, ...SensitiveOption) Sensitive[string]
		fn   func(string) string
		in   string
	}{
		{"NewPhone", NewPhone, Mobile, "0987654321"},
		{"NewEmail", NewEmail, Email, "ggw.chang@gmail.com"},
		{"NewPassword", NewPassword, Password, "password"},
		{"NewID", NewID, ID, "A123456789"},
		{"NewCredit", NewCredit, Credit, "4111111111111111"},
		{"NewName", NewName, Name, "John Doe"},
		{"NewAddress", NewAddress, Address, "台北市內湖區內湖路一段737巷1號1樓"},
		{"NewTel", NewTel, Tel, "0227993078"},
		{"NewURL", NewURL, URL, "http://john:password@localhost:3000"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.ctor(c.in)
			if got, want := s.String(), c.fn(c.in); got != want {
				t.Fatalf("%s masked = %q, want %q", c.name, got, want)
			}
			if s.Reveal() != c.in {
				t.Fatalf("%s Reveal() = %q, want %q", c.name, s.Reveal(), c.in)
			}
		})
	}
}

// TestPhoneVsTel 確認 NewPhone（手機）與 NewTel（市話）不混用，輸出格式不同。
func TestPhoneVsTel(t *testing.T) {
	if got := NewPhone("0987654321").String(); got != "0987***321" {
		t.Fatalf("NewPhone = %q, want %q", got, "0987***321")
	}
	if got := NewTel("0227993078").String(); got != "(02)2799-****" {
		t.Fatalf("NewTel = %q, want %q", got, "(02)2799-****")
	}
}

func ExampleNewPhone() {
	s := NewPhone("0987654321")
	fmt.Println(s)
	fmt.Println(s.Reveal())
	// Output:
	// 0987***321
	// 0987654321
}

func ExampleNewEmail() {
	fmt.Println(NewEmail("ggw.chang@gmail.com"))
	// Output: ggw****@gmail.com
}
