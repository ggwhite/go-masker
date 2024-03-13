package masker

import "testing"

func TestMobileMasker_Marshal(t *testing.T) {
	type args struct {
		s string
		i string
	}
	tests := []struct {
		name string
		m    *MobileMasker
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
			name: "Happy Pass",
			args: args{
				s: "*",
				i: "0978978978",
			},
			want: "0978***978",
		},
		{
			name: "Happy Pass",
			args: args{
				s: "*",
				i: "0912345678",
			},
			want: "0912***678",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MobileMasker{}
			if got := m.Marshal(tt.args.s, tt.args.i); got != tt.want {
				t.Errorf("MobileMasker.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
