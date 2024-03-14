package masker

import "testing"

func TestCreditMasker_Marshal(t *testing.T) {
	type args struct {
		s string
		i string
	}
	tests := []struct {
		name string
		m    *CreditMasker
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
			name: "VISA JCB MasterCard",
			args: args{
				s: "*",
				i: "1234567890123456",
			},
			want: "123456******3456",
		},
		{
			name: "American Express",
			args: args{
				s: "*",
				i: "123456789012345",
			},
			want: "123456******345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &CreditMasker{}
			if got := m.Marshal(tt.args.s, tt.args.i); got != tt.want {
				t.Errorf("CreditMasker.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
