package masker

import "testing"

func TestNameMasker_Mask(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "Empty Input", value: "", want: ""},
		{name: "Chinese Length 1", value: "王", want: "**"},
		{name: "Chinese Length 2", value: "王蛋", want: "王**"},
		{name: "Chinese Length 3", value: "王八蛋", want: "王**蛋"},
		{name: "Chinese Length 4", value: "王七八蛋", want: "王**蛋"},
		{name: "Chinese Length 5", value: "王七八九蛋", want: "王**九蛋"},
		{name: "Chinese Length 6", value: "王七八九十蛋", want: "王**九十蛋"},
		{name: "English Length 4", value: "Alen", want: "A**n"},
		{name: "English Full Name", value: "Alen Lin", want: "A**n L**n"},
		{name: "English Full Name 2", value: "Jorge Marry", want: "J**ge M**ry"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &NameMasker{mask: "*"}
			if got := m.Mask(tt.value); got != tt.want {
				t.Errorf("NameMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
