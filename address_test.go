package masker

import "testing"

func TestAddressMasker_Mask(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "Empty Input", value: "", want: ""},
		{name: "Happy Pass", value: "台北市大安區敦化南路五段7788號378樓", want: "台北市大安區******"},
		{name: "Length Less Than 6", value: "台北市", want: "******"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AddressMasker{mask: "*"}
			if got := m.Mask(tt.value); got != tt.want {
				t.Errorf("AddressMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
