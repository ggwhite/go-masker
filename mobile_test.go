package masker

import "testing"

func TestMobileMasker_Mask(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "Empty Input", value: "", want: ""},
		{name: "Happy Pass", value: "0978978978", want: "0978***978"},
		{name: "Happy Pass 2", value: "0912345678", want: "0912***678"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MobileMasker{mask: "*"}
			if got := m.Mask(tt.value); got != tt.want {
				t.Errorf("MobileMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
