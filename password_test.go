package masker

import "testing"

func TestPasswordMasker_Mask(t *testing.T) {
	tests := []struct {
		name  string
		mask  string
		value string
		want  string
	}{
		{name: "Empty Input", mask: "*", value: "", want: "**************"},
		{name: "Short", mask: "*", value: "1234567", want: "**************"},
		{name: "Mixed", mask: "*", value: "abcd!@#$%321", want: "**************"},
		{name: "Custom Mask Char", mask: "@", value: "abcd!@#$%321", want: "@@@@@@@@@@@@@@"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PasswordMasker{mask: tt.mask}
			if got := m.Mask(tt.value); got != tt.want {
				t.Errorf("PasswordMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
