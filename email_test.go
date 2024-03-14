package masker

import "testing"

func TestEmailMasker_Marshal(t *testing.T) {
	type args struct {
		s string
		i string
	}
	tests := []struct {
		name string
		m    *EmailMasker
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
				i: "ggw.chang@gmail.com",
			},
			want: "ggw****ng@gmail.com",
		},
		{
			name: "Address Less Than 3",
			args: args{
				s: "*",
				i: "qq@gmail.com",
			},
			want: "qq****@gmail.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &EmailMasker{}
			if got := m.Marshal(tt.args.s, tt.args.i); got != tt.want {
				t.Errorf("EmailMasker.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
