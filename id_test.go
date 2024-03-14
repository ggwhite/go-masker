package masker

import "testing"

func TestIDMasker_Marshal(t *testing.T) {
	type args struct {
		s string
		i string
	}
	tests := []struct {
		name string
		m    *IDMasker
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
				i: "A123456789",
			},
			want: "A12345****",
		},
		{
			name: "Length Less Than 6",
			args: args{
				s: "*",
				i: "A12",
			},
			want: "A12****",
		},
		{
			name: "Length Less Than 6",
			args: args{
				s: "*",
				i: "A",
			},
			want: "A****",
		},
		{
			name: "Length Between 6 and 10",
			args: args{
				s: "*",
				i: "A123456",
			},
			want: "A12345****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IDMasker{}
			if got := m.Marshal(tt.args.s, tt.args.i); got != tt.want {
				t.Errorf("IDMasker.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
