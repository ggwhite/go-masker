package masker

import "testing"

func TestCreditMasker_Mask(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "Empty Input", value: "", want: ""},
		{name: "VISA JCB MasterCard", value: "1234567890123456", want: "123456******3456"},
		{name: "American Express", value: "123456789012345", want: "123456******345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &CreditMasker{mask: "*"}
			if got := m.Mask(tt.value); got != tt.want {
				t.Errorf("CreditMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
