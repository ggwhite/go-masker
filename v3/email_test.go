package masker

import "testing"

func TestEmailMasker_Mask(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "Empty Input", value: "", want: ""},
		{name: "Happy Pass", value: "ggw.chang@gmail.com", want: "ggw****@gmail.com"},
		{name: "Address Less Than 3", value: "qq@gmail.com", want: "qq****@gmail.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &EmailMasker{mask: "*"}
			if got := m.Mask(tt.value); got != tt.want {
				t.Errorf("EmailMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
