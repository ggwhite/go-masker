package masker

import "testing"

func TestIDMasker_Mask(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "Empty Input", value: "", want: ""},
		{name: "Happy Pass", value: "A123456789", want: "A12345****"},
		{name: "Length 3", value: "A12", want: "A12****"},
		{name: "Length 1", value: "A", want: "A****"},
		{name: "Length 7", value: "A123456", want: "A12345****"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IDMasker{mask: "*"}
			if got := m.Mask(tt.value); got != tt.want {
				t.Errorf("IDMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
