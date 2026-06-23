package masker

import (
	"sync"
	"testing"
)

// 本檔驗證 v3 對外 API 表面齊備且行為對照 v2（對照 masker_test.go 的
// TestNewMaskerMarshaler / TestMaskerMarshaler_List / TestMaskerMarshaler_ConcurrentMarshal）。

// TestNewMaskerMarshaler_DefaultMaskers 對照 AC-17：List() 涵蓋全部靜態 masker type，
// 動態 first-N / last-N 不需註冊即可 Marshal。
func TestNewMaskerMarshaler_DefaultMaskers(t *testing.T) {
	m := NewMaskerMarshaler()

	wantTypes := []MaskerType{
		TypeNone, TypePassword, TypeName, TypeAddress, TypeEmail,
		TypeMobile, TypeTel, TypeID, TypeCredit, TypeURL, TypeAbuse, TypeAll,
	}

	got := make(map[MaskerType]struct{}, len(m.List()))
	for _, ty := range m.List() {
		got[ty] = struct{}{}
	}
	for _, ty := range wantTypes {
		if _, ok := got[ty]; !ok {
			t.Errorf("List() 缺少靜態 masker type: %v", ty)
		}
	}

	// 動態 tag 不需註冊
	if _, ok := got["first-3"]; ok {
		t.Error("first-N 不應出現在 List()（屬動態 tag）")
	}
	if v, err := m.Marshal("first-3", "hello"); err != nil || v != "***lo" {
		t.Errorf("Marshal(first-3) = %v, err %v, want ***lo", v, err)
	}
	if v, err := m.Marshal("last-3", "hello"); err != nil || v != "he***" {
		t.Errorf("Marshal(last-3) = %v, err %v, want he***", v, err)
	}
}

// TestDefaultMaskerMarshaler_Usable 對照 AC-15：全域實例存在、可用、遮罩字元為 '*'。
func TestDefaultMaskerMarshaler_Usable(t *testing.T) {
	if DefaultMaskerMarshaler == nil {
		t.Fatal("DefaultMaskerMarshaler 為 nil")
	}
	got, err := DefaultMaskerMarshaler.Marshal(TypeMobile, "0987654321")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "0987***321" {
		t.Errorf("Default Marshal = %v, want 0987***321", got)
	}
}

// TestMarshalAndGet_UnknownType 對照 AC-15：未知型別的 Marshal / Get 回傳 error。
func TestMarshalAndGet_UnknownType(t *testing.T) {
	m := NewMaskerMarshaler()
	if _, err := m.Marshal("nonexistent", "x"); err == nil {
		t.Error("Marshal(unknown) 應回傳 error")
	}
	if _, err := m.Get("nonexistent"); err == nil {
		t.Error("Get(unknown) 應回傳 error")
	}
}

// TestRegisterOverride 對照 AC-15：Register 覆蓋既有 masker。
func TestRegisterOverride(t *testing.T) {
	m := NewMaskerMarshaler()
	m.Register(TypeAll, &NoneMasker{})
	got, err := m.Marshal(TypeAll, "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret" {
		t.Errorf("覆蓋後 Marshal(TypeAll) = %v, want secret", got)
	}
}

// TestConcurrent_MarshalAndStruct 對照 v2 TestMaskerMarshaler_ConcurrentMarshal（AC-18）：
// 並行 Marshal 與 Struct 在 -race 下無資料競爭。
func TestConcurrent_MarshalAndStruct(t *testing.T) {
	type foo struct {
		Name string `mask:"name"`
		ID   string `mask:"id"`
	}
	m := NewMaskerMarshaler()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, _ = m.Marshal(TypeName, "John Doe")
		}()
		go func() {
			defer wg.Done()
			_, _ = m.Struct(&foo{Name: "John Doe", ID: "A123456789"})
		}()
	}
	wg.Wait()
}
