package masker

import "testing"

func TestTelephoneMasker_Mask(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "Empty Input", value: "", want: ""},
		{name: "With Special Chart", value: "(02-)27   99-3--078", want: "(02)2799-****"},
		{name: "Happy Pass", value: "0227993078", want: "(02)2799-****"},
		{name: "Happy Pass 2", value: "0788079966", want: "(07)8807-****"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TelephoneMasker{mask: "*"}
			if got := m.Mask(tt.value); got != tt.want {
				t.Errorf("TelephoneMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTelephoneMasker_CustomMaskChar(t *testing.T) {
	m := &TelephoneMasker{mask: "#"}
	got := m.Mask("0227993078")
	want := "(02)2799-####"
	if got != want {
		t.Errorf("got = %v, want %v", got, want)
	}
}
