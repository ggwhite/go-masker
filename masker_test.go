package masker

import (
	"reflect"
	"testing"
)

func TestMasker_Name(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "A",
			m:    &Masker{},
			args: args{
				i: "王蛋",
			},
			want: "王**",
		},
		{
			name: "B",
			m:    &Masker{},
			args: args{
				i: "王八蛋",
			},
			want: "王**",
		},
		{
			name: "C",
			m:    &Masker{},
			args: args{
				i: "王七八蛋",
			},
			want: "王**蛋",
		},
		{
			name: "D",
			m:    &Masker{},
			args: args{
				i: "王七八九蛋",
			},
			want: "王**九蛋",
		},
		{
			name: "E",
			m:    &Masker{},
			args: args{
				i: "王七八九十蛋",
			},
			want: "王**九十蛋",
		},
		{
			name: "F",
			m:    &Masker{},
			args: args{
				i: "Alen",
			},
			want: "A**n",
		},
		{
			name: "G",
			m:    &Masker{},
			args: args{
				i: "Alen Lin",
			},
			want: "A**n Lin",
		},
		{
			name: "H",
			m:    &Masker{},
			args: args{
				i: "Jorge Marry",
			},
			want: "J**ge Marry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Name(tt.args.i); got != tt.want {
				t.Errorf("Masker.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_ID(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "A",
			m:    &Masker{},
			args: args{
				i: "A123456789",
			},
			want: "A12345****",
		},
		{
			name: "B",
			m:    &Masker{},
			args: args{
				i: "A12",
			},
			want: "A12****",
		},
		{
			name: "C",
			m:    &Masker{},
			args: args{
				i: "A",
			},
			want: "A****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ID(tt.args.i); got != tt.want {
				t.Errorf("Masker.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Address(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "A",
			m:    &Masker{},
			args: args{
				i: "台北市大安區敦化南路五段7788號378樓",
			},
			want: "台北市大安區******",
		},
		{
			name: "B",
			m:    &Masker{},
			args: args{
				i: "台北市",
			},
			want: "台北市******",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Address(tt.args.i); got != tt.want {
				t.Errorf("Masker.Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_CreditCard(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "A",
			m:    &Masker{},
			args: args{
				i: "1234567890123456",
			},
			want: "123456******3456",
		},
		{
			name: "B",
			m:    &Masker{},
			args: args{
				i: "123456789012345",
			},
			want: "123456******345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.CreditCard(tt.args.i); got != tt.want {
				t.Errorf("Masker.CreditCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Email(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "A",
			m:    &Masker{},
			args: args{
				i: "ggw.chang@gmail.com",
			},
			want: "ggw****ng@gmail.com",
		},
		{
			name: "B",
			m:    &Masker{},
			args: args{
				i: "qq@gmail.com",
			},
			want: "qq****@gmail.com",
		},
		{
			name: "C",
			m:    &Masker{},
			args: args{
				i: "qqabcd@yahoo.com.tw",
			},
			want: "qqa****@yahoo.com.tw",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Email(tt.args.i); got != tt.want {
				t.Errorf("Masker.Email() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Mobile(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "A",
			m:    &Masker{},
			args: args{
				i: "0978978978",
			},
			want: "0978***978",
		},
		{
			name: "B",
			m:    &Masker{},
			args: args{
				i: "0912345678",
			},
			want: "0912***678",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Mobile(tt.args.i); got != tt.want {
				t.Errorf("Masker.Mobile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Telephone(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "A",
			m:    &Masker{},
			args: args{
				i: "(02-)27   99-3--078",
			},
			want: "(02)2799-****",
		},
		{
			name: "B",
			m:    &Masker{},
			args: args{
				i: "0227993078",
			},
			want: "(02)2799-****",
		},
		{
			name: "C",
			m:    &Masker{},
			args: args{
				i: "0788079966",
			},
			want: "(07)8807-****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Telephone(tt.args.i); got != tt.want {
				t.Errorf("Masker.Telephone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Password(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "A",
			m:    &Masker{},
			args: args{
				i: "1234567",
			},
			want: "************",
		},
		{
			name: "B",
			m:    &Masker{},
			args: args{
				i: "abcd!@#$%321",
			},
			want: "************",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Password(tt.args.i); got != tt.want {
				t.Errorf("Masker.Password() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Masker
	}{
		{
			name: "New Instance",
			want: &Masker{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "A",
			args: args{
				i: "王蛋",
			},
			want: "王**",
		},
		{
			name: "B",
			args: args{
				i: "王八蛋",
			},
			want: "王**",
		},
		{
			name: "C",
			args: args{
				i: "王七八蛋",
			},
			want: "王**蛋",
		},
		{
			name: "D",
			args: args{
				i: "王七八九蛋",
			},
			want: "王**九蛋",
		},
		{
			name: "E",
			args: args{
				i: "王七八九十蛋",
			},
			want: "王**九十蛋",
		},
		{
			name: "F",
			args: args{
				i: "Alen",
			},
			want: "A**n",
		},
		{
			name: "G",
			args: args{
				i: "Alen Lin",
			},
			want: "A**n Lin",
		},
		{
			name: "H",
			args: args{
				i: "Jorge Marry",
			},
			want: "J**ge Marry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Name(tt.args.i); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestID(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "A",
			args: args{
				i: "A123456789",
			},
			want: "A12345****",
		},
		{
			name: "B",
			args: args{
				i: "A12",
			},
			want: "A12****",
		},
		{
			name: "C",
			args: args{
				i: "A",
			},
			want: "A****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ID(tt.args.i); got != tt.want {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddress(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "A",
			args: args{
				i: "台北市大安區敦化南路五段7788號378樓",
			},
			want: "台北市大安區******",
		},
		{
			name: "B",
			args: args{
				i: "台北市",
			},
			want: "台北市******",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Address(tt.args.i); got != tt.want {
				t.Errorf("Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreditCard(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "A",
			args: args{
				i: "1234567890123456",
			},
			want: "123456******3456",
		},
		{
			name: "B",
			args: args{
				i: "123456789012345",
			},
			want: "123456******345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreditCard(tt.args.i); got != tt.want {
				t.Errorf("CreditCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "A",
			args: args{
				i: "ggw.chang@gmail.com",
			},
			want: "ggw****ng@gmail.com",
		},
		{
			name: "B",
			args: args{
				i: "qq@gmail.com",
			},
			want: "qq****@gmail.com",
		},
		{
			name: "C",
			args: args{
				i: "qqabcd@yahoo.com.tw",
			},
			want: "qqa****@yahoo.com.tw",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Email(tt.args.i); got != tt.want {
				t.Errorf("Email() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMobile(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "A",
			args: args{
				i: "0978978978",
			},
			want: "0978***978",
		},
		{
			name: "B",
			args: args{
				i: "0912345678",
			},
			want: "0912***678",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mobile(tt.args.i); got != tt.want {
				t.Errorf("Mobile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTelephone(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "A",
			args: args{
				i: "(02-)27   99-3--078",
			},
			want: "(02)2799-****",
		},
		{
			name: "B",
			args: args{
				i: "0227993078",
			},
			want: "(02)2799-****",
		},
		{
			name: "C",
			args: args{
				i: "0788079966",
			},
			want: "(07)8807-****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Telephone(tt.args.i); got != tt.want {
				t.Errorf("Telephone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "A",
			args: args{
				i: "1234567",
			},
			want: "************",
		},
		{
			name: "B",
			args: args{
				i: "abcd!@#$%321",
			},
			want: "************",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Password(tt.args.i); got != tt.want {
				t.Errorf("Password() = %v, want %v", got, tt.want)
			}
		})
	}
}
