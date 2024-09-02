package masker

import (
	"reflect"
	"testing"
)

func TestMaskerMarshaler_Marshal(t *testing.T) {
	type fields struct {
		Maskers map[MaskerType]Masker
		masker  string
	}
	type args struct {
		t     MaskerType
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test marshal",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone: &NoneMasker{},
					MaskerTypeID:   &IDMasker{},
				},
				masker: "*",
			},
			args: args{
				t:     MaskerTypeID,
				value: "A123456789",
			},
			want:    "A12345****",
			wantErr: false,
		},
		{
			name: "test marshal not found",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone: &NoneMasker{},
				},
				masker: "*",
			},
			args: args{
				t:     MaskerTypeID,
				value: "A123456789",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MaskerMarshaler{
				Maskers: tt.fields.Maskers,
				masker:  tt.fields.masker,
			}
			got, err := m.Marshal(tt.args.t, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaskerMarshaler.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MaskerMarshaler.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskerMarshaler_Register(t *testing.T) {
	type fields struct {
		Maskers map[MaskerType]Masker
		masker  string
	}
	type args struct {
		t      MaskerType
		masker Masker
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test register",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone: &NoneMasker{},
				},
				masker: "*",
			},
			args: args{
				t:      MaskerTypeID,
				masker: &IDMasker{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MaskerMarshaler{
				Maskers: tt.fields.Maskers,
				masker:  tt.fields.masker,
			}
			m.Register(tt.args.t, tt.args.masker)
		})
	}
}

func TestMaskerMarshaler_Unregister(t *testing.T) {
	type fields struct {
		Maskers map[MaskerType]Masker
		masker  string
	}
	type args struct {
		t MaskerType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test unregister",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone: &NoneMasker{},
					MaskerTypeID:   &IDMasker{},
				},
				masker: "*",
			},
			args: args{
				t: MaskerTypeID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MaskerMarshaler{
				Maskers: tt.fields.Maskers,
				masker:  tt.fields.masker,
			}
			m.Unregister(tt.args.t)
		})
	}
}

func TestMaskerMarshaler_Get(t *testing.T) {
	type fields struct {
		Maskers map[MaskerType]Masker
		masker  string
	}
	type args struct {
		t MaskerType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Masker
		wantErr bool
	}{
		{
			name: "test get",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone: &NoneMasker{},
					MaskerTypeID:   &IDMasker{},
				},
				masker: "*",
			},
			args: args{
				t: MaskerTypeID,
			},
			want:    &IDMasker{},
			wantErr: false,
		},
		{
			name: "test get not found",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone: &NoneMasker{},
				},
				masker: "*",
			},
			args: args{
				t: MaskerTypeID,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MaskerMarshaler{
				Maskers: tt.fields.Maskers,
				masker:  tt.fields.masker,
			}
			got, err := m.Get(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaskerMarshaler.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaskerMarshaler.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskerMarshaler_List(t *testing.T) {
	type fields struct {
		Maskers map[MaskerType]Masker
		masker  string
	}
	tests := []struct {
		name   string
		fields fields
		want   []MaskerType
	}{
		{
			name: "test list",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone: &NoneMasker{},
					MaskerTypeID:   &IDMasker{},
				},
				masker: "*",
			},
			want: []MaskerType{MaskerTypeNone, MaskerTypeID},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MaskerMarshaler{
				Maskers: tt.fields.Maskers,
				masker:  tt.fields.masker,
			}
			if got := m.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaskerMarshaler.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskerMarshaler_SetMasker(t *testing.T) {
	type fields struct {
		Maskers map[MaskerType]Masker
		masker  string
	}
	type args struct {
		masker string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test set masker",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone: &NoneMasker{},
					MaskerTypeID:   &IDMasker{},
				},
				masker: "*",
			},
			args: args{
				masker: "#",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MaskerMarshaler{
				Maskers: tt.fields.Maskers,
				masker:  tt.fields.masker,
			}
			m.SetMasker(tt.args.masker)
		})
	}
}

func TestMaskerMarshaler_Struct(t *testing.T) {
	type fields struct {
		Maskers map[MaskerType]Masker
		masker  string
	}
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test struct",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone:  &NoneMasker{},
					MaskerTypeID:    &IDMasker{},
					MaskerTypeName:  &NameMasker{},
					MaskerTypeEmail: &EmailMasker{},
				},
				masker: "*",
			},
			args: args{
				s: struct {
					ID      string `mask:"id"`
					Name    string `mask:"name"`
					Account struct {
						Name  string `mask:"name"`
						Email string `mask:"email"`
					} `mask:"struct"`
				}{
					ID:   "A123456789",
					Name: "John Doe",
					Account: struct {
						Name  string `mask:"name"`
						Email string `mask:"email"`
					}{
						Name:  "John Doe",
						Email: "ggw.chang@gmail.com",
					},
				},
			},
			want: &struct {
				ID      string `mask:"id"`
				Name    string `mask:"name"`
				Account struct {
					Name  string `mask:"name"`
					Email string `mask:"email"`
				} `mask:"struct"`
			}{
				ID:   "A12345****",
				Name: "J**n D**e",
				Account: struct {
					Name  string `mask:"name"`
					Email string `mask:"email"`
				}{
					Name:  "J**n D**e",
					Email: "ggw****@gmail.com",
				},
			},
		},
		{
			name: "test struct not found",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone: &NoneMasker{},
				},
				masker: "*",
			},
			args: args{
				s: struct {
					ID string `mask:"id"`
				}{
					ID: "A123456789",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test struct with slice string",
			fields: fields{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone: &NoneMasker{},
					MaskerTypeID:   &IDMasker{},
				},
				masker: "*",
			},
			args: args{
				s: struct {
					IDs []string `mask:"id"`
				}{
					IDs: []string{"A123456789", "A123456789"},
				},
			},
			want: &struct {
				IDs []string `mask:"id"`
			}{
				IDs: []string{"A12345****", "A12345****"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MaskerMarshaler{
				Maskers: tt.fields.Maskers,
				masker:  tt.fields.masker,
			}
			got, err := m.Struct(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaskerMarshaler.Struct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaskerMarshaler.Struct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMaskerMarshaler(t *testing.T) {
	tests := []struct {
		name string
		want *MaskerMarshaler
	}{
		{
			name: "test new masker marshaler",
			want: &MaskerMarshaler{
				Maskers: map[MaskerType]Masker{
					MaskerTypeNone:     &NoneMasker{},
					MaskerTypePassword: &PasswordMasker{},
					MaskerTypeName:     &NameMasker{},
					MaskerTypeAddress:  &AddressMasker{},
					MaskerTypeEmail:    &EmailMasker{},
					MaskerTypeMobile:   &MobileMasker{},
					MaskerTypeTel:      &TelephoneMasker{},
					MaskerTypeID:       &IDMasker{},
					MaskerTypeCredit:   &CreditMasker{},
					MaskerTypeURL:      &URLMasker{},
				},
				masker: "*",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMaskerMarshaler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMaskerMarshaler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strLoop(t *testing.T) {
	type args struct {
		str    string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test str loop",
			args: args{
				str:    "*",
				length: 6,
			},
			want: "******",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strLoop(tt.args.str, tt.args.length); got != tt.want {
				t.Errorf("strLoop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_overlay(t *testing.T) {
	type args struct {
		str     string
		overlay string
		start   int
		end     int
	}
	tests := []struct {
		name          string
		args          args
		wantOverlayed string
	}{
		{
			name: "test overlay",
			args: args{
				str:     "A123456789",
				overlay: "*",
				start:   6,
				end:     10,
			},
			wantOverlayed: "A12345*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOverlayed := overlay(tt.args.str, tt.args.overlay, tt.args.start, tt.args.end); gotOverlayed != tt.wantOverlayed {
				t.Errorf("overlay() = %v, want %v", gotOverlayed, tt.wantOverlayed)
			}
		})
	}
}
