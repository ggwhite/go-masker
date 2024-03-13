package masker

import "testing"

func TestAddressMasker_Marshal(t *testing.T) {
	type args struct {
		s string
		i string
	}
	tests := []struct {
		name string
		m    *AddressMasker
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
				i: "台北市大安區敦化南路五段7788號378樓",
			},
			want: "台北市大安區******",
		},
		{
			name: "Length Less Than 6",
			args: args{
				s: "*",
				i: "台北市",
			},
			want: "******",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AddressMasker{}
			if got := m.Marshal(tt.args.s, tt.args.i); got != tt.want {
				t.Errorf("AddressMasker.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
