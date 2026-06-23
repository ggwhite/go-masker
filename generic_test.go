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

func TestMarshal_GenericTags_Error(t *testing.T) {
	m := NewMaskerMarshaler()
	if _, err := m.Marshal("first-abc", "hello"); err == nil {
		t.Errorf("expected error for invalid first-N tag")
	}
}
