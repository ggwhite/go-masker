package masker

import "testing"

func TestNameMasker_Marshal(t *testing.T) {
	type args struct {
		s string
		i string
	}
	tests := []struct {
		name string
		m    *NameMasker
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
			name: "Chinese Length 1",
			args: args{
				s: "*",
				i: "王",
			},
			want: "**",
		},
		{
			name: "Chinese Length 2",
			args: args{
				s: "*",
				i: "王蛋",
			},
			want: "王**",
		},
		{
			name: "Chinese Length 3",
			args: args{
				s: "*",
				i: "王八蛋",
			},
			want: "王**蛋",
		},
		{
			name: "Chinese Length 4",
			args: args{
				s: "*",
				i: "王七八蛋",
			},
			want: "王**蛋",
		},
		{
			name: "Chinese Length 5",
			args: args{
				s: "*",
				i: "王七八九蛋",
			},
			want: "王**九蛋",
		},
		{
			name: "Chinese Length 6",
			args: args{
				s: "*",
				i: "王七八九十蛋",
			},
			want: "王**九十蛋",
		},
		{
			name: "English Length 4",
			args: args{
				s: "*",
				i: "Alen",
			},
			want: "A**n",
		},
		{
			name: "English Full Name",
			args: args{
				s: "*",
				i: "Alen Lin",
			},
			want: "A**n L**n",
		},
		{
			name: "English Full Name",
			args: args{
				s: "*",
				i: "Jorge Marry",
			},
			want: "J**ge M**ry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &NameMasker{}
			if got := m.Marshal(tt.args.s, tt.args.i); got != tt.want {
				t.Errorf("NameMasker.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
