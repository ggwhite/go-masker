package masker

import "testing"

// TestWithMaskChar 驗證 WithMaskChar 生效且不污染 DefaultMaskerMarshaler（AC-4）。
func TestWithMaskChar(t *testing.T) {
	custom := NewMaskerMarshaler(WithMaskChar('#'))

	got, err := custom.Marshal(TypeMobile, "0987654321")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "0987###321" {
		t.Errorf("custom mask char got %v, want 0987###321", got)
	}

	// 預設 marshaler 仍用 '*'，不受 custom 影響
	def, err := NewMaskerMarshaler().Marshal(TypeMobile, "0987654321")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if def != "0987***321" {
		t.Errorf("default mask char got %v, want 0987***321", def)
	}

	// package-level 便利函式走 DefaultMaskerMarshaler，仍輸出 '*'
	if pkg := Mobile("0987654321"); pkg != "0987***321" {
		t.Errorf("package-level Mobile got %v, want 0987***321", pkg)
	}
}

// TestWithMaskChar_URLAndNoneException 驗證 AC-4b 的例外。
func TestWithMaskChar_URLAndNoneException(t *testing.T) {
	custom := NewMaskerMarshaler(WithMaskChar('#'))

	url, err := custom.Marshal(TypeURL, "http://u:p@host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url != "http://u:xxxxx@host" {
		t.Errorf("URL got %v, want http://u:xxxxx@host", url)
	}

	none, err := custom.Marshal(TypeNone, "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if none != "secret" {
		t.Errorf("None got %v, want secret", none)
	}
}

// TestMarshalAndMustMarshal 驗證 Marshal 回 error、MustMarshal 在未知型別時 panic（AC-11）。
func TestMarshalAndMustMarshal(t *testing.T) {
	m := NewMaskerMarshaler()

	got, err := m.Marshal(TypeMobile, "0987654321")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if must := m.MustMarshal(TypeMobile, "0987654321"); must != got {
		t.Errorf("MustMarshal got %v, want %v", must, got)
	}

	if _, err := m.Marshal("bad", "x"); err == nil {
		t.Errorf("expected error for unknown type")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for unknown type in MustMarshal")
		}
	}()
	m.MustMarshal("bad", "x")
}

// TestStruct 驗證 Struct() 在 v3 仍可用且遮罩結果正確（AC-12）。
func TestStruct(t *testing.T) {
	type Foo struct {
		Name      string `mask:"name"`
		Email     string `mask:"email"`
		Password  string `mask:"password"`
		ID        string `mask:"id"`
		Mobile    string `mask:"mobile"`
		Telephone string `mask:"tel"`
		Credit    string `mask:"credit"`
		URL       string `mask:"url"`
		Plain     string
	}

	m := NewMaskerMarshaler()
	out, err := m.Struct(&Foo{
		Name:      "John Doe",
		Email:     "john@gmail.com",
		Password:  "password",
		ID:        "1234567890",
		Mobile:    "1234567890",
		Telephone: "0227993078",
		Credit:    "4111111111111111",
		URL:       "http://john:password@localhost:3000",
		Plain:     "keep",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.(*Foo)
	want := Foo{
		Name:      "J**n D**e",
		Email:     "joh****@gmail.com",
		Password:  "**************",
		ID:        "123456****",
		Mobile:    "1234***890",
		Telephone: "(02)2799-****",
		Credit:    "411111******1111",
		URL:       "http://john:xxxxx@localhost:3000",
		Plain:     "keep",
	}
	if *got != want {
		t.Errorf("Struct() got %+v, want %+v", *got, want)
	}
}

// TestRegisterGetListUnregister 覆蓋 marshaler 的 registry 操作。
func TestRegisterGetListUnregister(t *testing.T) {
	m := NewMaskerMarshaler()

	if _, err := m.Get(TypeMobile); err != nil {
		t.Errorf("Get(TypeMobile) unexpected error: %v", err)
	}
	if len(m.List()) == 0 {
		t.Errorf("List() should not be empty")
	}

	m.Unregister(TypeMobile)
	if _, err := m.Get(TypeMobile); err == nil {
		t.Errorf("expected error after Unregister")
	}

	m.Register(TypeMobile, &MobileMasker{mask: "*"})
	if _, err := m.Get(TypeMobile); err != nil {
		t.Errorf("Get after Register unexpected error: %v", err)
	}
}
