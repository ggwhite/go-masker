package masker

import "testing"

func TestTelephoneMasker_Mask_Default(t *testing.T) {
	tests := []struct {
		name   string
		masker *TelephoneMasker
		value  string
		want   string
	}{
		{
			name:   "Empty Input",
			masker: &TelephoneMasker{mask: "*"},
			value:  "",
			want:   "",
		},
		{
			name:   "With Special Characters - 10 digits",
			masker: &TelephoneMasker{mask: "*"},
			value:  "(02-)27   99-3--078",
			want:   "(02)2799-****",
		},
		{
			name:   "Happy Pass - 10 digits",
			masker: &TelephoneMasker{mask: "*"},
			value:  "0227993078",
			want:   "(02)2799-****",
		},
		{
			name:   "Happy Pass - 8 digits",
			masker: &TelephoneMasker{mask: "*"},
			value:  "27993078",
			want:   "2799-****",
		},
		{
			name:   "Length Mismatch - Returns Original",
			masker: &TelephoneMasker{mask: "*"},
			value:  "12345",
			want:   "12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.masker.Mask(tt.value); got != tt.want {
				t.Errorf("TelephoneMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTelephoneMasker_Mask_Configurable(t *testing.T) {
	tests := []struct {
		name   string
		masker *TelephoneMasker
		value  string
		want   string
	}{
		{
			name:   "Empty Input",
			masker: &TelephoneMasker{mask: "*", regionCodeLen: 2, numberLen: 8},
			value:  "",
			want:   "",
		},
		{
			name:   "With Special Characters",
			masker: &TelephoneMasker{mask: "*", regionCodeLen: 2, numberLen: 8},
			value:  "(02-)27   99-3--078",
			want:   "02 2799****",
		},
		{
			name:   "Happy Pass",
			masker: &TelephoneMasker{mask: "*", regionCodeLen: 2, numberLen: 8},
			value:  "0227993078",
			want:   "02 2799****",
		},
		{
			name:   "Happy Pass 2",
			masker: &TelephoneMasker{mask: "*", regionCodeLen: 2, numberLen: 8},
			value:  "0788079966",
			want:   "07 8807****",
		},
		{
			name:   "Length Mismatch - All Masked",
			masker: &TelephoneMasker{mask: "*", regionCodeLen: 2, numberLen: 8},
			value:  "12345",
			want:   "*****",
		},
		{
			name:   "NumberLen Less Than 4 - All Masked",
			masker: &TelephoneMasker{mask: "*", regionCodeLen: 2, numberLen: 3},
			value:  "02123",
			want:   "*****",
		},
		{
			name:   "Custom RegionCodeLen",
			masker: &TelephoneMasker{mask: "*", regionCodeLen: 3, numberLen: 7},
			value:  "0201234567",
			want:   "020 123****",
		},
		{
			name:   "Zero RegionCodeLen - No Leading Space",
			masker: &TelephoneMasker{mask: "*", regionCodeLen: 0, numberLen: 8},
			value:  "12345678",
			want:   "1234****",
		},
		{
			name:   "Zero RegionCodeLen With Special Characters",
			masker: &TelephoneMasker{mask: "*", regionCodeLen: 0, numberLen: 8},
			value:  "12-34 56-78",
			want:   "1234****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.masker.Mask(tt.value); got != tt.want {
				t.Errorf("TelephoneMasker.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTelephoneMasker_CustomMaskChar(t *testing.T) {
	t.Run("Default mode", func(t *testing.T) {
		m := &TelephoneMasker{mask: "#"}
		got := m.Mask("0227993078")
		want := "(02)2799-####"
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Configurable mode", func(t *testing.T) {
		m := &TelephoneMasker{mask: "#", regionCodeLen: 2, numberLen: 8}
		got := m.Mask("0227993078")
		want := "02 2799####"
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Zero regionCodeLen", func(t *testing.T) {
		m := &TelephoneMasker{mask: "#", regionCodeLen: 0, numberLen: 8}
		got := m.Mask("12345678")
		want := "1234####"
		if got != want {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}
