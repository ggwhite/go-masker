package masker

import "testing"

func TestPasswordMasker_Marshal(t *testing.T) {
	type args struct {
		s string
		i string
	}
	tests := []struct {
		name string
		m    *PasswordMasker
		args args
		want string
	}{
		{
			name: "Empty Input",
			args: args{
				s: "*",
				i: "",
			},
			want: "**************",
		},
		{
			name: "Happy Pass",
			args: args{
				s: "*",
				i: "1234567",
			},
			want: "**************",
		},
		{
			name: "Happy Pass",
			args: args{
				s: "*",
				i: "abcd!@#$%321",
			},
			want: "**************",
		},
		{
			name: "Happy Pass",
			args: args{
				s: "@",
				i: "abcd!@#$%321",
			},
			want: "@@@@@@@@@@@@@@",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PasswordMasker{}
			if got := m.Marshal(tt.args.s, tt.args.i); got != tt.want {
				t.Errorf("PasswordMasker.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
