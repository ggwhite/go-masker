package masker

import (
	"fmt"
	"reflect"
)

const tagName = "mask"

// MaskerType is a string type for masker type
type MaskerType string

// MaskerType constants
const (
	MaskerTypeNone     MaskerType = "none"
	MaskerTypePassword MaskerType = "password"
	MaskerTypeName     MaskerType = "name"
	MaskerTypeAddress  MaskerType = "addr"
	MaskerTypeEmail    MaskerType = "email"
	MaskerTypeMobile   MaskerType = "mobile"
	MaskerTypeTel      MaskerType = "tel"
	MaskerTypeID       MaskerType = "id"
	MaskerTypeCredit   MaskerType = "credit"
	MaskerTypeURL      MaskerType = "url"
	MaskerTypeStruct   MaskerType = "struct"
)

// Masker is an interface for masking sensitive data
type Masker interface {
	Marshal(string, string) string
}

// MaskerMarshaler is a masker marshaler
type MaskerMarshaler struct {
	Maskers map[MaskerType]Masker
	masker  string // default masker
}

func (m *MaskerMarshaler) Marshal(t MaskerType, value string) (string, error) {
	masker, ok := m.Maskers[t]
	if !ok {
		return "", fmt.Errorf("masker %v not found", t)
	}
	return masker.Marshal(m.masker, value), nil
}

func (m *MaskerMarshaler) Register(t MaskerType, masker Masker) {
	m.Maskers[t] = masker
}

func (m *MaskerMarshaler) Unregister(t MaskerType) {
	delete(m.Maskers, t)
}

func (m *MaskerMarshaler) Get(t MaskerType) (Masker, error) {
	masker, ok := m.Maskers[t]
	if !ok {
		return nil, fmt.Errorf("masker %v not found", t)
	}
	return masker, nil
}

func (m *MaskerMarshaler) List() []MaskerType {
	var list []MaskerType
	for t := range m.Maskers {
		list = append(list, t)
	}
	return list
}

func (m *MaskerMarshaler) SetMasker(masker string) {
	m.masker = masker
}

// Struct must input a interface{}, add tag mask on struct fields, after Struct(), return a pointer interface{} of input type and it will be masked with the tag format type
//
// Example:
func (m *MaskerMarshaler) Struct(s interface{}) (interface{}, error) {
	if s == nil {
		return nil, fmt.Errorf("input is nil")
	}

	var selem, tptr reflect.Value

	st := reflect.TypeOf(s)

	if st.Kind() == reflect.Ptr {
		tptr = reflect.New(st.Elem())
		selem = reflect.ValueOf(s).Elem()
	} else {
		tptr = reflect.New(st)
		selem = reflect.ValueOf(s)
	}

	for i := 0; i < selem.NumField(); i++ {
		if !selem.Type().Field(i).IsExported() {
			continue
		}
		mtag := selem.Type().Field(i).Tag.Get(tagName)
		if len(mtag) == 0 {
			tptr.Elem().Field(i).Set(selem.Field(i))
			continue
		}
		switch selem.Field(i).Type().Kind() {
		default:
			tptr.Elem().Field(i).Set(selem.Field(i))
		case reflect.String:
			v, err := m.Marshal(MaskerType(mtag), selem.Field(i).String())
			if err != nil {
				return nil, err
			}
			tptr.Elem().Field(i).SetString(v)
		case reflect.Struct:
			if MaskerType(mtag) == MaskerTypeStruct {
				_t, err := m.Struct(selem.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t).Elem())
			}
		case reflect.Ptr:
			if selem.Field(i).IsNil() {
				continue
			}
			if MaskerType(mtag) == MaskerTypeStruct {
				_t, err := m.Struct(selem.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t))
			}
		case reflect.Slice:
			if selem.Field(i).IsNil() {
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.String {
				orgval := selem.Field(i).Interface().([]string)
				newval := make([]string, len(orgval))
				for i, val := range selem.Field(i).Interface().([]string) {
					v, err := m.Marshal(MaskerType(mtag), val)
					if err != nil {
						return nil, err
					}
					newval[i] = v
				}
				tptr.Elem().Field(i).Set(reflect.ValueOf(newval))
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.Struct && MaskerType(mtag) == MaskerTypeStruct {
				newval := reflect.MakeSlice(selem.Field(i).Type(), 0, selem.Field(i).Len())
				for j, l := 0, selem.Field(i).Len(); j < l; j++ {
					_n, err := m.Struct(selem.Field(i).Index(j).Interface())
					if err != nil {
						return nil, err
					}
					newval = reflect.Append(newval, reflect.ValueOf(_n).Elem())
				}
				tptr.Elem().Field(i).Set(newval)
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.Ptr && MaskerType(mtag) == MaskerTypeStruct {
				newval := reflect.MakeSlice(selem.Field(i).Type(), 0, selem.Field(i).Len())
				for j, l := 0, selem.Field(i).Len(); j < l; j++ {
					_n, err := m.Struct(selem.Field(i).Index(j).Interface())
					if err != nil {
						return nil, err
					}
					newval = reflect.Append(newval, reflect.ValueOf(_n))
				}
				tptr.Elem().Field(i).Set(newval)
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.Interface && MaskerType(mtag) == MaskerTypeStruct {
				newval := reflect.MakeSlice(selem.Field(i).Type(), 0, selem.Field(i).Len())
				for j, l := 0, selem.Field(i).Len(); j < l; j++ {
					_n, err := m.Struct(selem.Field(i).Index(j).Interface())
					if err != nil {
						return nil, err
					}
					if reflect.TypeOf(selem.Field(i).Index(j).Interface()).Kind() != reflect.Ptr {
						newval = reflect.Append(newval, reflect.ValueOf(_n).Elem())
					} else {
						newval = reflect.Append(newval, reflect.ValueOf(_n))
					}
				}
				tptr.Elem().Field(i).Set(newval)
				continue
			}
		case reflect.Interface:
			if selem.Field(i).IsNil() {
				continue
			}
			if MaskerType(mtag) != MaskerTypeStruct {
				continue
			}
			_t, err := m.Struct(selem.Field(i).Interface())
			if err != nil {
				return nil, err
			}
			if reflect.TypeOf(selem.Field(i).Interface()).Kind() != reflect.Ptr {
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t).Elem())
			} else {
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t))
			}
		}
	}

	return tptr.Interface(), nil
}

func NewMaskerMarshaler() *MaskerMarshaler {
	return &MaskerMarshaler{
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
	}
}

var DefaultMaskerMarshaler = &MaskerMarshaler{
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
}

func strLoop(str string, length int) string {
	var mask string
	for i := 1; i <= length; i++ {
		mask += str
	}
	return mask
}

func overlay(str string, overlay string, start int, end int) (overlayed string) {
	r := []rune(str)
	l := len([]rune(r))

	if l == 0 {
		return ""
	}

	if start < 0 {
		start = 0
	}
	if start > l {
		start = l
	}
	if end < 0 {
		end = 0
	}
	if end > l {
		end = l
	}
	if start > end {
		tmp := start
		start = end
		end = tmp
	}

	overlayed = ""
	overlayed += string(r[:start])
	overlayed += overlay
	overlayed += string(r[end:])
	return overlayed
}
