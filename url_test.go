package masker

import "testing"

func TestURLMasker_Mask(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "With Password", value: "http://john:password@localhost:3000", want: "http://john:xxxxx@localhost:3000"},
		{name: "No Password", value: "http://localhost:3000", want: "http://localhost:3000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &URLMasker{}
			if got := m.Mask(tt.value); got != tt.want {
				t.Errorf("URLMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoneMasker_Mask(t *testing.T) {
	m := &NoneMasker{}
	if got := m.Mask("secret"); got != "secret" {
		t.Errorf("NoneMasker.Mask() = %v, want secret", got)
	}
}

func TestAllMasker_Mask(t *testing.T) {
	m := &AllMasker{mask: "*"}
	if got := m.Mask("secret"); got != "******" {
		t.Errorf("AllMasker.Mask() = %v, want ******", got)
	}
}
