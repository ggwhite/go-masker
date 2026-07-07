package masker

import "testing"

func TestParseGenericMask(t *testing.T) {
	tests := []struct {
		tag     string
		value   string
		want    string
		matched bool
		wantErr bool
	}{
		{tag: "all", value: "hello", want: "", matched: false},
		{tag: "all", value: "", want: "", matched: false},
		{tag: "first-3", value: "hello", want: "***lo", matched: true},
		{tag: "first-0", value: "hello", want: "hello", matched: true},
		{tag: "first-10", value: "hello", want: "*****", matched: true},
		{tag: "last-3", value: "hello", want: "he***", matched: true},
		{tag: "last-0", value: "hello", want: "hello", matched: true},
		{tag: "last-10", value: "hello", want: "*****", matched: true},
		{tag: "name", value: "hello", want: "", matched: false},
		{tag: "first-abc", value: "hello", want: "", matched: false, wantErr: true},
		{tag: "last-abc", value: "hello", want: "", matched: false, wantErr: true},
		{tag: "first--1", value: "hello", want: "", matched: false, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.tag+"/"+tt.value, func(t *testing.T) {
			got, matched, err := parseGenericMask("*", tt.tag, tt.value)
			if matched != tt.matched {
				t.Errorf("matched = %v, want %v", matched, tt.matched)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if matched && !tt.wantErr && got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseGenericMask_Tel(t *testing.T) {
	tests := []struct {
		tag     string
		value   string
		want    string
		matched bool
		wantErr bool
	}{
		// AC-1: tel-2-8, issue #40 範例
		{tag: "tel-2-8", value: "0227993078", want: "02-2799-****", matched: true},
		// AC-2: tel-3-8, issue #40 範例
		{tag: "tel-3-8", value: "75588888888", want: "755-8888-****", matched: true},
		// AC-3: tel-3-7, issue #40 範例
		{tag: "tel-3-7", value: "2125551234", want: "212-555-****", matched: true},
		// AC-4: 3 段含國際碼
		{tag: "tel-2-3-8", value: "8675588888888", want: "+86-755-8888-****", matched: true},
		// AC-5: 4 段含國際碼 + space
		{tag: "tel-2-3-8-space", value: "8675588888888", want: "+86 755 8888 ****", matched: true},
		// AC-6: 2 段 + space
		{tag: "tel-2-8-space", value: "0227993078", want: "02 2799 ****", matched: true},
		// AC-7: numberLen == 4 邊界
		{tag: "tel-2-4", value: "021234", want: "02-****", matched: true},
		// AC-8: 已格式化未遮罩輸入
		{tag: "tel-2-8", value: "(02)2799-3078", want: "02-2799-****", matched: true},
		// AC-9: 含 + 的已格式化輸入
		{tag: "tel-2-3-8", value: "+86-755-8888-8888", want: "+86-755-8888-****", matched: true},
		// AC-10: 長度不符（多 1 碼）
		{tag: "tel-2-8", value: "00227993078", want: "00227993078", matched: true},
		// AC-11: 長度不符（少 1 碼）
		{tag: "tel-2-8", value: "022799307", want: "022799307", matched: true},
		// AC-12: 空字串
		{tag: "tel-2-8", value: "", want: "", matched: true},
		// AC-13: regionLen <= 0
		{tag: "tel-0-8", value: "0227993078", matched: false, wantErr: true},
		// AC-14: numberLen < 4
		{tag: "tel-2-3", value: "0227993078", matched: false, wantErr: true},
		// AC-15: 分段數錯誤（1 段）
		{tag: "tel-2", value: "0227993078", matched: false, wantErr: true},
		// AC-16: 分段數錯誤（5 段）
		{tag: "tel-2-3-8-dash-extra", value: "0227993078", matched: false, wantErr: true},
		// AC-17: 第三段消歧義失敗
		{tag: "tel-2-8-foo", value: "0227993078", matched: false, wantErr: true},
		// AC-18: 第四段非法分隔符
		{tag: "tel-2-3-8-comma", value: "0227993078", matched: false, wantErr: true},
		// AC-19: intlLen <= 0（3 段時）
		{tag: "tel-0-3-8", value: "0227993078", matched: false, wantErr: true},
		// AC-25: numberLen == 4 + 含國際碼
		{tag: "tel-2-3-4", value: "867551234", want: "+86-755-****", matched: true},
		// AC-26: intlLen <= 0（4 段）
		{tag: "tel-0-3-8-dash", value: "0227993078", matched: false, wantErr: true},
		// AC-27: regionLen <= 0（3 段含國際碼）
		{tag: "tel-2-0-8", value: "0227993078", matched: false, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.tag+"/"+tt.value, func(t *testing.T) {
			got, matched, err := parseGenericMask("*", tt.tag, tt.value)
			if matched != tt.matched {
				t.Errorf("matched = %v, want %v", matched, tt.matched)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if matched && !tt.wantErr && got != tt.want {
				t.Errorf("got = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMarshal_GenericTags(t *testing.T) {
	m := NewMaskerMarshaler()
	tests := []struct {
		tag   MaskerType
		value string
		want  string
	}{
		{TypeAll, "hello", "*****"},
		{"first-2", "hello", "**llo"},
		{"last-2", "hello", "hel**"},
	}
	for _, tt := range tests {
		t.Run(string(tt.tag)+"/"+tt.value, func(t *testing.T) {
			got, err := m.Marshal(tt.tag, tt.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshal_GenericTags_Tel(t *testing.T) {
	m := NewMaskerMarshaler()
	tests := []struct {
		tag   MaskerType
		value string
		want  string
	}{
		// AC-20: 端對端整合
		{"tel-2-8", "0227993078", "02-2799-****"},
		// AC-21: 回歸——既有 tel 行為不變
		{TypeTel, "0227993078", "(02)2799-****"},
		{TypeTel, "27993078", "2799-****"},
	}
	for _, tt := range tests {
		t.Run(string(tt.tag)+"/"+tt.value, func(t *testing.T) {
			got, err := m.Marshal(tt.tag, tt.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseGenericMask_MobileID(t *testing.T) {
	tests := []struct {
		tag     string
		value   string
		want    string
		matched bool
		wantErr bool
	}{
		// mobile 正常案例
		{tag: "mobile-3-4", value: "09012345678", want: "090****5678", matched: true},
		{tag: "mobile-3-4", value: "2025551234", want: "202***1234", matched: true},
		{tag: "mobile-0-4", value: "447911123456", want: "********3456", matched: true},
		{tag: "mobile-4-3", value: "0987654321", want: "0987***321", matched: true},
		{tag: "mobile-4-0", value: "0987654321", want: "0987******", matched: true},

		// id 正常案例
		{tag: "id-0-4", value: "123456789", want: "*****6789", matched: true},
		{tag: "id-4-0", value: "123456789012", want: "1234********", matched: true},
		{tag: "id-3-3", value: "S1234567D", want: "S12***67D", matched: true},
		{tag: "id-6-0", value: "A123456789", want: "A12345****", matched: true},
		{tag: "id-2-1", value: "AB1234567", want: "AB******7", matched: true},

		// 邊界案例
		{tag: "mobile-5-5", value: "1234567890", want: "1234567890", matched: true},
		{tag: "mobile-3-4", value: "12345", want: "12345", matched: true},
		{tag: "id-3-0", value: "", want: "", matched: true},
		{tag: "mobile-1-1", value: "AB", want: "AB", matched: true},

		// 多位元組字元（AC-26）
		{tag: "mobile-2-2", value: "你好世界測試", want: "你好**測試", matched: true},
		{tag: "id-1-1", value: "甲乙丙", want: "甲*丙", matched: true},

		// mid 通用標籤
		{tag: "mid-3-4", value: "09012345678", want: "090****5678", matched: true},
		{tag: "mid-0-4", value: "123456789", want: "*****6789", matched: true},
		{tag: "mid-2-2", value: "ABCDEF", want: "AB**EF", matched: true},
		{tag: "mid-1-0", value: "secret", want: "s*****", matched: true},
		{tag: "mid-5-5", value: "1234567890", want: "1234567890", matched: true},
		{tag: "mid-2-2", value: "你好世界測試", want: "你好**測試", matched: true},
		{tag: "mid-0-0", value: "123", matched: false, wantErr: true},

		// 錯誤案例
		{tag: "mobile-0-0", value: "123", matched: false, wantErr: true},
		{tag: "id-0-0", value: "123", matched: false, wantErr: true},
		{tag: "mobile-3", value: "123", matched: false, wantErr: true},
		{tag: "id-1-2-3", value: "123", matched: false, wantErr: true},
		{tag: "mobile-abc-4", value: "123", matched: false, wantErr: true},
		{tag: "id-3-def", value: "123", matched: false, wantErr: true},
		{tag: "mobile--1-4", value: "123", matched: false, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.tag+"/"+tt.value, func(t *testing.T) {
			got, matched, err := parseGenericMask("*", tt.tag, tt.value)
			if matched != tt.matched {
				t.Errorf("matched = %v, want %v", matched, tt.matched)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if matched && !tt.wantErr && got != tt.want {
				t.Errorf("got = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMarshal_GenericTags_MobileID(t *testing.T) {
	m := NewMaskerMarshaler()
	tests := []struct {
		tag   MaskerType
		value string
		want  string
	}{
		// 端對端整合
		{"mobile-3-4", "09012345678", "090****5678"},
		{"id-0-4", "123456789", "*****6789"},
		{"mid-2-2", "ABCDEF", "AB**EF"},
		// 回歸——既有行為不變
		{TypeMobile, "0987654321", "0987***321"},
		{TypeID, "A123456789", "A12345****"},
	}
	for _, tt := range tests {
		t.Run(string(tt.tag)+"/"+tt.value, func(t *testing.T) {
			got, err := m.Marshal(tt.tag, tt.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMarshal_GenericTags_Error(t *testing.T) {
	m := NewMaskerMarshaler()
	if _, err := m.Marshal("first-abc", "hello"); err == nil {
		t.Errorf("expected error for invalid first-N tag")
	}
}
