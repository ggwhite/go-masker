package masker

import "testing"

func TestTelephoneMasker_Marshal(t *testing.T) {
	type args struct {
		s string
		i string
	}
	tests := []struct {
		name string
		m    *TelephoneMasker
		args args
		want string
	}{
		{
			name: "Empty Input",
			args: args{
				s: "*",
				i: "",
			},
			want: "",
		},
		{
			name: "With Special Chart",
			args: args{
				s: "*",
				i: "(02-)27   99-3--078",
			},
			want: "(02)2799-****",
		},
		{
			name: "Happy Pass",
			args: args{
				s: "*",
				i: "0227993078",
			},
			want: "(02)2799-****",
		},
		{
			name: "Happy Pass",
			args: args{
				s: "*",
				i: "0788079966",
			},
			want: "(07)8807-****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &TelephoneMasker{}
			if got := m.Marshal(tt.args.s, tt.args.i); got != tt.want {
				t.Errorf("TelephoneMasker.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
