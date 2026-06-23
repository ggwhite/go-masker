package masker

import (
	"reflect"
	"testing"
)

// 本檔移植 v2 root package 的 Struct() parity 測試案例（對照 masker_test.go
// 的 TestMaskerMarshaler_Struct / TestMaskerMarshaler_MapStructs / TestStruct_*）。
// 期望輸出值與 v2 逐字一致，僅將 API 名稱適配為 v3（常數改名、建構子改用 NewMaskerMarshaler）。

// TestStruct_StringFieldsAllTypes 對照 AC-2：各 masker type 於 string 欄位的遮罩輸出。
func TestStruct_StringFieldsAllTypes(t *testing.T) {
	type allTypes struct {
		Name     string `mask:"name"`
		Email    string `mask:"email"`
		Password string `mask:"password"`
		ID       string `mask:"id"`
		Mobile   string `mask:"mobile"`
		Tel      string `mask:"tel"`
		Credit   string `mask:"credit"`
		Addr     string `mask:"addr"`
		URL      string `mask:"url"`
		None     string `mask:"none"`
		All      string `mask:"all"`
		Plain    string
	}

	m := NewMaskerMarshaler()
	out, err := m.Struct(&allTypes{
		Name:     "John Doe",
		Email:    "ggw.chang@gmail.com",
		Password: "secret",
		ID:       "A123456789",
		Mobile:   "0987654321",
		Tel:      "0227993078",
		Credit:   "4111111111111111",
		Addr:     "台北市大安區敦化南路五段7788號378樓",
		URL:      "http://john:password@localhost:3000",
		None:     "secret",
		All:      "secret",
		Plain:    "keep",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := allTypes{
		Name:     "J**n D**e",
		Email:    "ggw****@gmail.com",
		Password: "**************",
		ID:       "A12345****",
		Mobile:   "0987***321",
		Tel:      "(02)2799-****",
		Credit:   "411111******1111",
		Addr:     "台北市大安區******",
		URL:      "http://john:xxxxx@localhost:3000",
		None:     "secret",
		All:      "******",
		Plain:    "keep",
	}
	if got := out.(*allTypes); *got != want {
		t.Errorf("Struct() = %+v, want %+v", *got, want)
	}
}

// TestStruct_NestedStructValue 對照 v2 TestMaskerMarshaler_Struct 的 nested value type（AC-3）。
func TestStruct_NestedStructValue(t *testing.T) {
	type account struct {
		Name  string `mask:"name"`
		Email string `mask:"email"`
	}
	type person struct {
		ID      string  `mask:"id"`
		Name    string  `mask:"name"`
		Account account `mask:"struct"`
	}

	m := NewMaskerMarshaler()
	out, err := m.Struct(&person{
		ID:      "A123456789",
		Name:    "John Doe",
		Account: account{Name: "John Doe", Email: "ggw.chang@gmail.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := person{
		ID:      "A12345****",
		Name:    "J**n D**e",
		Account: account{Name: "J**n D**e", Email: "ggw****@gmail.com"},
	}
	if got := out.(*person); *got != want {
		t.Errorf("Struct() = %+v, want %+v", *got, want)
	}
}

// TestStruct_NestedStructPointer 對照 AC-3 的 *Struct + mask:"struct"。
func TestStruct_NestedStructPointer(t *testing.T) {
	type account struct {
		Name  string `mask:"name"`
		Email string `mask:"email"`
	}
	type person struct {
		Name    string   `mask:"name"`
		Account *account `mask:"struct"`
	}

	m := NewMaskerMarshaler()
	out, err := m.Struct(&person{
		Name:    "John Doe",
		Account: &account{Name: "John Doe", Email: "ggw.chang@gmail.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.(*person)
	if got.Name != "J**n D**e" {
		t.Errorf("Name = %v, want J**n D**e", got.Name)
	}
	if got.Account == nil {
		t.Fatal("Account is nil, want masked pointer")
	}
	if *got.Account != (account{Name: "J**n D**e", Email: "ggw****@gmail.com"}) {
		t.Errorf("Account = %+v, want {J**n D**e ggw****@gmail.com}", *got.Account)
	}
}

// TestStruct_SliceString 對照 v2 "test struct with slice string"（AC-4）。
func TestStruct_SliceString(t *testing.T) {
	type holder struct {
		IDs []string `mask:"id"`
	}
	m := NewMaskerMarshaler()
	out, err := m.Struct(&holder{IDs: []string{"A123456789", "A123456789"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := &holder{IDs: []string{"A12345****", "A12345****"}}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("Struct() = %#v, want %#v", out, want)
	}
}

// TestStruct_SliceStruct 對照 AC-4 的 []struct + mask:"struct"。
func TestStruct_SliceStruct(t *testing.T) {
	type inner struct {
		ID string `mask:"id"`
	}
	type holder struct {
		Items []inner `mask:"struct"`
	}
	m := NewMaskerMarshaler()
	out, err := m.Struct(&holder{Items: []inner{{ID: "A123456789"}, {ID: "B223456789"}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := &holder{Items: []inner{{ID: "A12345****"}, {ID: "B22345****"}}}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("Struct() = %#v, want %#v", out, want)
	}
}

// TestStruct_SlicePtrStruct 對照 AC-4 的 []*struct + mask:"struct"。
func TestStruct_SlicePtrStruct(t *testing.T) {
	type inner struct {
		ID string `mask:"id"`
	}
	type holder struct {
		Items []*inner `mask:"struct"`
	}
	m := NewMaskerMarshaler()
	out, err := m.Struct(&holder{Items: []*inner{{ID: "A123456789"}, {ID: "B223456789"}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.(*holder)
	if len(got.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(got.Items))
	}
	if got.Items[0].ID != "A12345****" || got.Items[1].ID != "B22345****" {
		t.Errorf("Items = [%v %v], want [A12345**** B22345****]", got.Items[0].ID, got.Items[1].ID)
	}
}

// TestStruct_SliceInterface 對照 AC-4 的 []interface{} + mask:"struct"（含 value 與 pointer 元素）。
func TestStruct_SliceInterface(t *testing.T) {
	type inner struct {
		Name string `mask:"name"`
	}
	type holder struct {
		Items []interface{} `mask:"struct"`
	}
	m := NewMaskerMarshaler()
	out, err := m.Struct(&holder{Items: []interface{}{
		inner{Name: "John Doe"},
		&inner{Name: "John Doe"},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.(*holder)
	if v, ok := got.Items[0].(inner); !ok || v.Name != "J**n D**e" {
		t.Errorf("Items[0] = %#v, want inner{J**n D**e}", got.Items[0])
	}
	if v, ok := got.Items[1].(*inner); !ok || v.Name != "J**n D**e" {
		t.Errorf("Items[1] = %#v, want *inner{J**n D**e}", got.Items[1])
	}
}

// TestStruct_InterfaceField 對照 AC-5：interface{} 欄位（value 與 pointer）+ mask:"struct"。
func TestStruct_InterfaceField(t *testing.T) {
	type inner struct {
		Name string `mask:"name"`
	}
	type box struct {
		Data interface{} `mask:"struct"`
	}
	m := NewMaskerMarshaler()

	// value 內容
	outVal, err := m.Struct(&box{Data: inner{Name: "John Doe"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := outVal.(*box).Data.(inner); !ok || v.Name != "J**n D**e" {
		t.Errorf("value Data = %#v, want inner{J**n D**e}", outVal.(*box).Data)
	}

	// pointer 內容
	outPtr, err := m.Struct(&box{Data: &inner{Name: "John Doe"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := outPtr.(*box).Data.(*inner); !ok || v.Name != "J**n D**e" {
		t.Errorf("pointer Data = %#v, want *inner{J**n D**e}", outPtr.(*box).Data)
	}
}

// TestStruct_MapStructs 逐字移植 v2 TestMaskerMarshaler_MapStructs（AC-6）。
func TestStruct_MapStructs(t *testing.T) {
	type item struct {
		ID string `mask:"id"`
	}

	tests := []struct {
		name string
		in   interface{}
		want interface{}
	}{
		{
			name: "map struct with mapstruct tag",
			in: struct {
				Items map[int]item `mask:"mapstruct"`
			}{
				Items: map[int]item{1: {ID: "A123456789"}},
			},
			want: &struct {
				Items map[int]item `mask:"mapstruct"`
			}{
				Items: map[int]item{1: {ID: "A12345****"}},
			},
		},
		{
			name: "nil map",
			in: struct {
				Items map[int]item `mask:"mapstruct"`
			}{
				Items: nil,
			},
			want: &struct {
				Items map[int]item `mask:"mapstruct"`
			}{
				Items: nil,
			},
		},
		{
			name: "empty non-nil map",
			in: struct {
				Items map[int]item `mask:"mapstruct"`
			}{Items: map[int]item{}},
			want: &struct {
				Items map[int]item `mask:"mapstruct"`
			}{Items: map[int]item{}},
		},
		{
			name: "map without mapstruct tag",
			in: struct {
				Items map[int]item `mask:"struct"`
			}{
				Items: map[int]item{1: {ID: "A123456789"}},
			},
			want: &struct {
				Items map[int]item `mask:"struct"`
			}{
				Items: map[int]item{1: {ID: "A123456789"}},
			},
		},
		{
			name: "map with scalar values passes through unchanged",
			in: struct {
				Items map[string]string `mask:"mapstruct"`
			}{Items: map[string]string{"key": "value"}},
			want: &struct {
				Items map[string]string `mask:"mapstruct"`
			}{Items: map[string]string{"key": "value"}},
		},
		{
			name: "map pointer values",
			in: struct {
				Items map[int]*item `mask:"mapstruct"`
			}{
				Items: map[int]*item{
					1: {ID: "A123456789"},
					2: nil,
				},
			},
			want: &struct {
				Items map[int]*item `mask:"mapstruct"`
			}{
				Items: map[int]*item{
					1: {ID: "A12345****"},
					2: nil,
				},
			},
		},
		{
			name: "map slice values",
			in: struct {
				Items map[int][]item `mask:"mapstruct"`
			}{
				Items: map[int][]item{
					1: {{ID: "A123456789"}, {ID: "B223456789"}},
					2: {{ID: "C323456789"}},
				},
			},
			want: &struct {
				Items map[int][]item `mask:"mapstruct"`
			}{
				Items: map[int][]item{
					1: {{ID: "A12345****"}, {ID: "B22345****"}},
					2: {{ID: "C32345****"}},
				},
			},
		},
		{
			name: "nil map slice values",
			in: struct {
				Items map[int][]item `mask:"mapstruct"`
			}{
				Items: nil,
			},
			want: &struct {
				Items map[int][]item `mask:"mapstruct"`
			}{
				Items: nil,
			},
		},
		{
			name: "map pointer to slice values",
			in: func() interface{} {
				items := []item{{ID: "A123456789"}, {ID: "B223456789"}}
				return struct {
					Items map[int]*[]item `mask:"mapstruct"`
				}{
					Items: map[int]*[]item{
						1: &items,
						2: nil,
					},
				}
			}(),
			want: func() interface{} {
				items := []item{{ID: "A12345****"}, {ID: "B22345****"}}
				return &struct {
					Items map[int]*[]item `mask:"mapstruct"`
				}{
					Items: map[int]*[]item{
						1: &items,
						2: nil,
					},
				}
			}(),
		},
		{
			name: "map of map with slice values",
			in: struct {
				Items map[int]map[int][]item `mask:"mapstruct"`
			}{
				Items: map[int]map[int][]item{
					1: {
						10: {{ID: "A123456789"}, {ID: "B223456789"}},
						11: nil,
					},
					2: nil,
				},
			},
			want: &struct {
				Items map[int]map[int][]item `mask:"mapstruct"`
			}{
				Items: map[int]map[int][]item{
					1: {
						10: {{ID: "A12345****"}, {ID: "B22345****"}},
						11: nil,
					},
					2: nil,
				},
			},
		},
		{
			name: "multi-level map with pointer slice leaves",
			in: func() interface{} {
				a := []item{{ID: "A123456789"}}
				b := []item{{ID: "B223456789"}, {ID: "C323456789"}}
				return struct {
					Items map[int]map[string]map[int]*[]item `mask:"mapstruct"`
				}{
					Items: map[int]map[string]map[int]*[]item{
						1: {
							"x": {
								100: &a,
								101: nil,
							},
						},
						2: {
							"y": {
								200: &b,
							},
						},
					},
				}
			}(),
			want: func() interface{} {
				a := []item{{ID: "A12345****"}}
				b := []item{{ID: "B22345****"}, {ID: "C32345****"}}
				return &struct {
					Items map[int]map[string]map[int]*[]item `mask:"mapstruct"`
				}{
					Items: map[int]map[string]map[int]*[]item{
						1: {
							"x": {
								100: &a,
								101: nil,
							},
						},
						2: {
							"y": {
								200: &b,
							},
						},
					},
				}
			}(),
		},
	}

	m := NewMaskerMarshaler()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.Struct(tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Struct() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

// TestStruct_NilPreservation 對照 AC-7：nil pointer / nil slice / nil map 欄位維持 nil 不 panic。
func TestStruct_NilPreservation(t *testing.T) {
	type inner struct {
		Name string `mask:"name"`
	}
	type holder struct {
		Ptr   *inner            `mask:"struct"`
		Slice []string          `mask:"id"`
		Map   map[int]inner     `mask:"mapstruct"`
		Plain string            `mask:"name"`
	}
	m := NewMaskerMarshaler()
	out, err := m.Struct(&holder{Plain: "John Doe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.(*holder)
	if got.Ptr != nil {
		t.Errorf("Ptr = %v, want nil", got.Ptr)
	}
	if got.Slice != nil {
		t.Errorf("Slice = %v, want nil", got.Slice)
	}
	if got.Map != nil {
		t.Errorf("Map = %v, want nil", got.Map)
	}
	if got.Plain != "J**n D**e" {
		t.Errorf("Plain = %v, want J**n D**e", got.Plain)
	}
}

// TestStruct_NonStructInput 對照 v2 PR #35：非 struct 輸入回傳 error 不 panic（AC-8）。
func TestStruct_NonStructInput(t *testing.T) {
	m := NewMaskerMarshaler()
	str := "hello"
	cases := []struct {
		name string
		in   interface{}
	}{
		{"string", "string input"},
		{"int", 42},
		{"ptr to string", &str},
		{"nil", nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := m.Struct(c.in); err == nil {
				t.Errorf("Struct(%v) expected error, got nil", c.name)
			}
		})
	}
}

// TestStruct_UnexportedSkipped 對照 AC-9：unexported 欄位跳過不處理、不報錯。
func TestStruct_UnexportedSkipped(t *testing.T) {
	type withUnexported struct {
		Name   string `mask:"name"`
		secret string `mask:"name"`
	}
	m := NewMaskerMarshaler()
	out, err := m.Struct(&withUnexported{Name: "John Doe", secret: "hidden"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := out.(*withUnexported)
	if got.Name != "J**n D**e" {
		t.Errorf("Name = %v, want J**n D**e", got.Name)
	}
	if got.secret != "" {
		t.Errorf("secret = %q, want \"\" (unexported 欄位跳過維持零值)", got.secret)
	}
}

// TestStruct_TypedNilPointer 對照 v2 同名測試（AC-10）。
func TestStruct_TypedNilPointer(t *testing.T) {
	type inner struct {
		Name string `mask:"name"`
	}
	m := NewMaskerMarshaler()
	if _, err := m.Struct((*inner)(nil)); err == nil {
		t.Fatal("expected error for typed-nil pointer, got nil")
	}
}

// TestStruct_SlicePtrNilElement 對照 v2 同名測試（AC-10）：slice 內含 nil 元素。
func TestStruct_SlicePtrNilElement(t *testing.T) {
	type inner struct {
		Name string `mask:"name"`
	}
	type outer struct {
		Items []*inner `mask:"struct"`
	}
	m := NewMaskerMarshaler()
	got, err := m.Struct(outer{Items: []*inner{{Name: "John"}, nil, {Name: "Jane"}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := got.(*outer)
	if result.Items[0].Name != "J**n" {
		t.Errorf("Items[0].Name = %v, want J**n", result.Items[0].Name)
	}
	if result.Items[1] != nil {
		t.Errorf("Items[1] = %v, want nil", result.Items[1])
	}
	if result.Items[2].Name != "J**e" {
		t.Errorf("Items[2].Name = %v, want J**e", result.Items[2].Name)
	}
}

// TestStruct_StructFieldWithNonStructTag 對照 v2 同名測試（AC-10）。
func TestStruct_StructFieldWithNonStructTag(t *testing.T) {
	type inner struct {
		Value string
	}
	type outer struct {
		Data inner `mask:"name"`
	}
	m := NewMaskerMarshaler()
	got, err := m.Struct(outer{Data: inner{Value: "hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.(*outer).Data.Value != "hello" {
		t.Errorf("Data.Value = %v, want hello (應原樣複製)", got.(*outer).Data.Value)
	}
}

// TestStruct_PtrFieldWithNonStructTag 對照 v2 同名測試（AC-10）。
func TestStruct_PtrFieldWithNonStructTag(t *testing.T) {
	type inner struct {
		Value string
	}
	type outer struct {
		Data *inner `mask:"name"`
	}
	m := NewMaskerMarshaler()
	got, err := m.Struct(outer{Data: &inner{Value: "hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := got.(*outer)
	if result.Data == nil {
		t.Fatal("Data is nil, expected to be copied")
	}
	if result.Data.Value != "hello" {
		t.Errorf("Data.Value = %v, want hello (應原樣複製)", result.Data.Value)
	}
}

// TestStruct_SliceOfIntPreserved 對照 v2 同名測試（AC-10）。
func TestStruct_SliceOfIntPreserved(t *testing.T) {
	type data struct {
		Scores []int `mask:"all"`
	}
	m := NewMaskerMarshaler()
	got, err := m.Struct(data{Scores: []int{1, 2, 3}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got.(*data).Scores, []int{1, 2, 3}) {
		t.Errorf("Scores = %v, want [1 2 3] (應原樣複製)", got.(*data).Scores)
	}
}

// TestStruct_GenericTags 對照 v2 generic_test.go 的 TestStruct_GenericTags（AC-11）。
func TestStruct_GenericTags(t *testing.T) {
	type secret struct {
		SSN    string `mask:"all"`
		Code   string `mask:"first-3"`
		Suffix string `mask:"last-4"`
	}
	m := NewMaskerMarshaler()
	out, err := m.Struct(&secret{SSN: "123456789", Code: "ABCDEF", Suffix: "ABCDEF"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := out.(*secret)
	if s.SSN != "*********" {
		t.Errorf("SSN = %v, want *********", s.SSN)
	}
	if s.Code != "***DEF" {
		t.Errorf("Code = %v, want ***DEF", s.Code)
	}
	if s.Suffix != "AB****" {
		t.Errorf("Suffix = %v, want AB****", s.Suffix)
	}
}
