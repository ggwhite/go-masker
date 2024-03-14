package masker

import "testing"

func TestURLMasker_Marshal(t *testing.T) {
	type args struct {
		s string
		i string
	}
	tests := []struct {
		name string
		m    *URLMasker
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &URLMasker{}
			if got := m.Marshal(tt.args.s, tt.args.i); got != tt.want {
				t.Errorf("URLMasker.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
