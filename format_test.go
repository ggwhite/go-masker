package masker

import (
	"strings"
	"testing"
)

type formatFoo struct {
	Name     string     `mask:"name"`
	Email    string     `mask:"email"`
	Password string     `mask:"password"`
	ID       string     `mask:"id"`
	Note     string     // 無 tag，直接複製
	Self     *formatFoo `mask:"struct"`
}

func newFormatFoo() *formatFoo {
	return &formatFoo{
		Name:     "John Doe",
		Email:    "john@gmail.com",
		Password: "password",
		ID:       "A123456789",
		Note:     "plain note",
		Self: &formatFoo{
			Name:     "Jane Doe",
			Email:    "jane@gmail.com",
			Password: "secret",
			ID:       "B223456789",
			Note:     "nested note",
		},
	}
}

// TestFormatNoPanic 驗證 AC-6：nil input 回空字串，且各種退化 input 皆不 panic。
func TestFormatNoPanic(t *testing.T) {
	m := NewMaskerMarshaler()

	if got := m.Format(nil); got != "" {
		t.Errorf("Format(nil) = %q, want empty string", got)
	}

	cases := []struct {
		name string
		in   interface{}
	}{
		{"typed nil pointer", (*formatFoo)(nil)},
		{"empty struct", struct{}{}},
		{"struct with nil pointer field", formatFoo{Name: "John Doe"}},
		{"non struct int", 42},
		{"non struct string", "hello"},
		{"slice", []string{"a", "b"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Format(%v) panicked: %v", c.in, r)
				}
			}()
			_ = m.Format(c.in)
		})
	}
}

// TestFormatParity 驗證 AC-7：Format() 對每個遮罩後 string 欄位的輸出與 Struct() byte-identical。
func TestFormatParity(t *testing.T) {
	m := NewMaskerMarshaler()
	in := newFormatFoo()

	masked, err := m.Struct(in)
	if err != nil {
		t.Fatalf("Struct() error: %v", err)
	}
	mf := masked.(*formatFoo)

	out := m.Format(in)

	for _, want := range []string{mf.Name, mf.Email, mf.Password, mf.ID, mf.Self.Name, mf.Self.Email, mf.Self.Password, mf.Self.ID} {
		if !strings.Contains(out, want) {
			t.Errorf("Format output %q does not contain masked value %q", out, want)
		}
	}
	// 無 tag 欄位直接保留原值
	if !strings.Contains(out, "plain note") {
		t.Errorf("Format output %q does not contain untagged value", out)
	}
}

// TestFormatDeterministic 驗證 AC-8：輸出不含記憶體位址、多次呼叫結果相同、ptr-to-struct 已遞迴展開。
func TestFormatDeterministic(t *testing.T) {
	m := NewMaskerMarshaler()
	in := newFormatFoo()

	out1 := m.Format(in)
	out2 := m.Format(in)
	if out1 != out2 {
		t.Errorf("Format not deterministic:\n  out1=%q\n  out2=%q", out1, out2)
	}
	if strings.Contains(out1, "0x") {
		t.Errorf("Format output contains memory address: %q", out1)
	}
	// ptr-to-struct 以 &{...} 展開，且巢狀遮罩後值已 inline
	if !strings.Contains(out1, "&{") {
		t.Errorf("Format output %q did not inline ptr-to-struct as &{...}", out1)
	}
	// 巢狀的 nil pointer 顯示 <nil>
	if !strings.HasSuffix(strings.TrimSuffix(out1, "}"), "<nil>}") {
		t.Errorf("Format output %q did not render nil pointer as <nil>", out1)
	}
}

// TestFormatMaskedContainers 驗證 slice / mapstruct / 巢狀 struct 欄位皆被遮罩，原始敏感值不外洩。
func TestFormatMaskedContainers(t *testing.T) {
	m := NewMaskerMarshaler()

	type item struct {
		ID string `mask:"id"`
	}
	in := struct {
		IDs   []string     `mask:"id"`
		Items map[int]item `mask:"mapstruct"`
		Ptrs  []*item      `mask:"struct"`
	}{
		IDs:   []string{"A123456789", "B223456789"},
		Items: map[int]item{1: {ID: "C323456789"}},
		Ptrs:  []*item{{ID: "D423456789"}, nil},
	}

	out := m.Format(in)
	leaks := []string{"A123456789", "B223456789", "C323456789", "D423456789"}
	for _, raw := range leaks {
		if strings.Contains(out, raw) {
			t.Errorf("Format output %q leaked raw sensitive value %q", out, raw)
		}
	}
	for _, want := range []string{"A12345****", "B22345****", "C32345****", "D42345****"} {
		if !strings.Contains(out, want) {
			t.Errorf("Format output %q missing masked value %q", out, want)
		}
	}
	if !strings.Contains(out, "<nil>") {
		t.Errorf("Format output %q missing <nil> for nil slice pointer element", out)
	}
}

// gtItem 是 ground-truth 測試用的巢狀 struct。
type gtItem struct {
	ID string `mask:"id"`
}

// gtFull 涵蓋 Format() 各條輸出分支，供逐字 ground-truth 斷言使用。
type gtFull struct {
	Name      string            `mask:"name"` // string 遮罩
	Note      string            // 無 tag → 原值
	hidden    string            `mask:"name"`         // 未匯出 → zero value
	Unknown   string            `mask:"nosuchmasker"` // 未知 tag → error fallback mask char
	Val       gtItem            `mask:"struct"`       // 巢狀 struct value
	Ptr       *gtItem           `mask:"struct"`       // 巢狀 struct ptr
	NilPtr    *gtItem           `mask:"struct"`       // nil pointer → <nil>
	IDs       []string          `mask:"id"`           // []string 遮罩
	NilSlice  []string          `mask:"id"`           // nil slice → []
	Structs   []gtItem          `mask:"struct"`       // []struct
	PtrSlice  []*gtItem         `mask:"struct"`       // []*struct（含 nil 元素）
	Ifaces    []interface{}     `mask:"struct"`       // []interface{}
	MapStruct map[string]gtItem `mask:"mapstruct"`    // map mapstruct
	NilMap    map[string]gtItem `mask:"mapstruct"`    // nil map → map[]
}

// TestFormatGroundTruth 驗證 AC-7：Format() 各輸出分支與 ground-truth 字串逐字一致。
func TestFormatGroundTruth(t *testing.T) {
	m := NewMaskerMarshaler()
	in := gtFull{
		Name:    "John Doe",
		Note:    "plain note",
		hidden:  "should be zeroed",
		Unknown: "leakme",
		Val:     gtItem{ID: "A123456789"},
		Ptr:     &gtItem{ID: "B223456789"},
		NilPtr:  nil,
		IDs:     []string{"C323456789", "D423456789"},
		Structs: []gtItem{{ID: "E523456789"}},
		PtrSlice: []*gtItem{
			{ID: "F623456789"},
			nil,
		},
		Ifaces:    []interface{}{gtItem{ID: "G723456789"}, &gtItem{ID: "H823456789"}},
		MapStruct: map[string]gtItem{"k1": {ID: "I923456789"}},
	}

	const want = "{J**n D**e plain note  ****** {A12345****} &{B22345****} <nil> [C32345**** D42345****] [] [{E52345****}] [&{F62345****} <nil>] [{G72345****} &{H82345****}] map[k1:{I92345****}] map[]}"

	got := m.Format(in)
	if got != want {
		t.Errorf("Format ground-truth mismatch:\n got: %q\nwant: %q", got, want)
	}
	// 未知 tag 的 error fallback 不得洩漏原值
	if strings.Contains(got, "leakme") {
		t.Errorf("Format leaked raw value via error fallback: %q", got)
	}
}

// TestFormatErrorFallbackMaskChar 驗證 parity 例外 #2：error fallback 使用 marshaler 的自訂遮罩字元，且絕不輸出原值。
func TestFormatErrorFallbackMaskChar(t *testing.T) {
	m := NewMaskerMarshaler(WithMaskChar('#'))
	in := struct {
		Unknown string `mask:"nosuchmasker"`
	}{Unknown: "leakme"}

	got := m.Format(in)
	const want = "{######}"
	if got != want {
		t.Errorf("Format with custom mask char = %q, want %q", got, want)
	}
	if strings.Contains(got, "leakme") {
		t.Errorf("Format leaked raw value: %q", got)
	}
}

type allocSeg struct {
	V string `mask:"name"`
}

// allocNest 含多個巢狀 struct 欄位：Struct() 須為每個巢狀 struct 各配置一個 masked struct（reflect.New + interface boxing），
// Format() 僅配置單一輸出緩衝，alloc 優勢隨巢狀數放大。此為「不建新 struct」訴求的代表性 fixture。
type allocNest struct {
	A allocSeg `mask:"struct"`
	B allocSeg `mask:"struct"`
	C allocSeg `mask:"struct"`
	D allocSeg `mask:"struct"`
	E allocSeg `mask:"struct"`
}

// TestFormatFewerAllocs 驗證 AC-9：Format() 的 allocs/op 低於 Struct()（不配置新 masked struct）。
func TestFormatFewerAllocs(t *testing.T) {
	m := NewMaskerMarshaler()
	in := allocNest{
		A: allocSeg{"John Doe"}, B: allocSeg{"Jane Doe"}, C: allocSeg{"Jim Doe"},
		D: allocSeg{"Joe Doe"}, E: allocSeg{"Jay Doe"},
	}
	resetTypeCache()
	_, _ = m.Struct(in)
	_ = m.Format(in)

	structAllocs := testing.AllocsPerRun(200, func() {
		_, _ = m.Struct(in)
	})
	formatAllocs := testing.AllocsPerRun(200, func() {
		_ = m.Format(in)
	})

	if formatAllocs >= structAllocs {
		t.Errorf("Format allocs/op (%.0f) not lower than Struct allocs/op (%.0f)", formatAllocs, structAllocs)
	}
}
