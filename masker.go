package masker

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const tagName = "mask"

// MaskerType is a string type for masker type
type MaskerType string

// MaskerType constants
const (
	MaskerTypeNone      MaskerType = "none"
	MaskerTypePassword  MaskerType = "password"
	MaskerTypeName      MaskerType = "name"
	MaskerTypeAddress   MaskerType = "addr"
	MaskerTypeEmail     MaskerType = "email"
	MaskerTypeMobile    MaskerType = "mobile"
	MaskerTypeTel       MaskerType = "tel"
	MaskerTypeID        MaskerType = "id"
	MaskerTypeCredit    MaskerType = "credit"
	MaskerTypeURL       MaskerType = "url"
	MaskerTypeAbuse     MaskerType = "abuse"
	MaskerTypeStruct    MaskerType = "struct"
	MaskerTypeAll       MaskerType = "all"
	MaskerTypeMapStruct MaskerType = "mapstruct"
)

// Masker is an interface for masking sensitive data
type Masker interface {
	Marshal(string, string) string
}

// MaskerMarshaler is a masker marshaler.
// All exported methods are safe for concurrent use.
// Direct access to the Maskers field is NOT concurrency-safe; use Register/Get instead.
type MaskerMarshaler struct {
	mu      sync.RWMutex
	Maskers map[MaskerType]Masker
	masker  string // default masker
}

// Marshal returns a masked value by masker type
// It is used for masking sensitive data
// Example:
//
//	m := masker.NewMaskerMarshaler()
//	log.Println(m.Marshal(masker.MaskerTypeNone, "none"))                               // none <nil>
//	log.Println(m.Marshal(masker.MaskerTypePassword, "password"))                       // ************** <nil>
//	log.Println(m.Marshal(masker.MaskerTypeName, "name"))                               // n**e <nil>
//	log.Println(m.Marshal(masker.MaskerTypeAddress, "address"))                         // addres****** <nil>
//	log.Println(m.Marshal(masker.MaskerTypeEmail, "email"))                             // ema**** <nil>
//	log.Println(m.Marshal(masker.MaskerTypeMobile, "mobile"))                           // mobi*** <nil>
//	log.Println(m.Marshal(masker.MaskerTypeTel, "tel"))                                 // tel <nil>
//	log.Println(m.Marshal(masker.MaskerTypeID, "id"))                                   // id**** <nil>
//	log.Println(m.Marshal(masker.MaskerTypeCredit, "4111111111111111"))                 // 411111******1111 <nil>
//	log.Println(m.Marshal(masker.MaskerTypeURL, "http://john:password@localhost:3000")) // http://john:xxxxx@localhost:3000 <nil>
func (m *MaskerMarshaler) Marshal(t MaskerType, value string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if result, matched, err := parseGenericMask(m.masker, string(t), value); matched {
		return result, err
	}
	masker, ok := m.Maskers[t]
	if !ok {
		return "", fmt.Errorf("masker %v not found", t)
	}
	return masker.Marshal(m.masker, value), nil
}

// Register adds a masker by masker type
// It is used for adding or override a masker by masker type
// Example:
//
//	m := masker.NewMaskerMarshaler()
//	m.Register(masker.MaskerTypePassword, &PasswordMasker{})
//	log.Println(m.List()) // [password name addr email tel id url none mobile credit]
func (m *MaskerMarshaler) Register(t MaskerType, masker Masker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Maskers[t] = masker
}

// Unregister removes a masker by masker type
// It is used for removing a masker by masker type
// Example:
//
//	m := masker.NewMaskerMarshaler()
//	m.Unregister(masker.MaskerTypePassword)
//	log.Println(m.List()) // [name addr email tel id url none mobile credit]
func (m *MaskerMarshaler) Unregister(t MaskerType) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Maskers, t)
}

// Get returns a masker by masker type
// It is used for getting a masker by masker type
// Example:
//
//	m := masker.NewMaskerMarshaler()
//	masker, _ := m.Get(masker.MaskerTypePassword)
//	log.Println(masker) // &{PasswordMasker}
func (m *MaskerMarshaler) Get(t MaskerType) (Masker, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	masker, ok := m.Maskers[t]
	if !ok {
		return nil, fmt.Errorf("masker %v not found", t)
	}
	return masker, nil
}

// List returns a list of masker types
// It is used for listing all masker types
// Example:
//
//	m := masker.NewMaskerMarshaler()
//	log.Println(m.List()) // [password name addr email tel id url none mobile credit]
func (m *MaskerMarshaler) List() []MaskerType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []MaskerType
	for t := range m.Maskers {
		list = append(list, t)
	}
	return list
}

func (m *MaskerMarshaler) SetMasker(masker string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.masker = masker
}

// Struct must input a interface{}, add tag mask on struct fields, after Struct(), return a pointer interface{} of input type and it will be masked with the tag format type
//
// Limitation:
//   - Cyclic object graphs are not specially handled. Self-referencing structures can
//     recurse indefinitely when using `mask:"struct"` or `mask:"mapstruct"`.
//   - Avoid cyclic references in masking paths.
//
// Example:
//
//	type Foo struct {
//		Name      string `mask:"name"`
//		Email     string `mask:"email"`
//		Password  string `mask:"password"`
//		ID        string `mask:"id"`
//		Address   string `mask:"addr"`
//		Mobile    string `mask:"mobile"`
//		Telephone string `mask:"tel"`
//		Credit    string `mask:"credit"`
//		URL       string `mask:"url"`
//		Foo       *Foo   `mask:"struct"`
//	}
//
//	func main() {
//		m := masker.NewMaskerMarshaler()
//		log.Println(m.List()) // [password name addr email tel id url none mobile credit]
//		foo1 := &Foo{
//			Name:      "John Doe",
//			Email:     "john@gmail.com",
//			Password:  "password",
//			ID:        "1234567890",
//			Address:   "123 Main St",
//			Mobile:    "1234567890",
//			Telephone: "1234567890",
//			Credit:    "4111111111111111",
//			URL:       "http://john:password@localhost:3000",
//			Foo: &Foo{
//				Name:      "John Doe",
//				Email:     "john@gmail.com",
//				Password:  "password",
//				ID:        "1234567890",
//				Address:   "123 Main St",
//				Mobile:    "1234567890",
//				Telephone: "1234567890",
//				Credit:    "4111111111111111",
//				URL:       "http://john:password@localhost:3000",
//			},
//		}
//
//		foo2, _ := m.Struct(foo1)
//
//		log.Println(foo1) // &{John Doe john@gmail.com password 1234567890 123 Main St 1234567890 1234567890 4111111111111111 http://john:password@localhost:3000 0xc0000001e0}
//		log.Println(foo1.Foo) // &{John Doe john@gmail.com password 1234567890 123 Main St 1234567890 1234567890 4111111111111111 http://john:password@localhost:3000 <nil>}
//		log.Println(foo2) // &{J**n D**e joh****@gmail.com ************** 123456**** 123 Ma****** 1234***890 (12)3456-**** 411111******1111 http://john:xxxxx@localhost:3000 0xc000000320}
//		log.Println(foo2.(*Foo).Foo) // &{J**n D**e joh****@gmail.com ************** 123456**** 123 Ma****** 1234***890 (12)3456-**** 411111******1111 http://john:xxxxx@localhost:3000 <nil>}
//	}
func (m *MaskerMarshaler) Struct(s interface{}) (interface{}, error) {
	if s == nil {
		return nil, fmt.Errorf("input is nil")
	}

	var selem, tptr reflect.Value

	st := reflect.TypeOf(s)
	sv := reflect.ValueOf(s)

	actualKind := st.Kind()
	if actualKind == reflect.Ptr {
		if sv.IsNil() {
			return nil, fmt.Errorf("input is nil pointer")
		}
		actualKind = st.Elem().Kind()
	}
	if actualKind != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to struct, got %v", actualKind)
	}

	if st.Kind() == reflect.Ptr {
		tptr = reflect.New(st.Elem())
		selem = sv.Elem()
	} else {
		tptr = reflect.New(st)
		selem = sv
	}

	info := cachedTypeInfo(selem.Type())
	for _, f := range info.fields {
		if !f.exported {
			continue
		}
		sf := selem.Field(f.index)
		dst := tptr.Elem().Field(f.index)
		if !f.hasTag {
			dst.Set(sf)
			continue
		}

		switch f.kind {
		default:
			dst.Set(sf)
		case reflect.String:
			v, err := m.Marshal(f.tag, sf.String())
			if err != nil {
				return nil, err
			}
			dst.SetString(v)
		case reflect.Struct:
			if f.tag == MaskerTypeStruct {
				_t, err := m.Struct(sf.Interface())
				if err != nil {
					return nil, err
				}
				dst.Set(reflect.ValueOf(_t).Elem())
			} else {
				dst.Set(sf)
			}
		case reflect.Ptr:
			if sf.IsNil() {
				continue
			}
			if f.tag == MaskerTypeStruct {
				_t, err := m.Struct(sf.Interface())
				if err != nil {
					return nil, err
				}
				dst.Set(reflect.ValueOf(_t))
			} else {
				dst.Set(sf)
			}
		case reflect.Map:
			if sf.IsNil() {
				continue
			}
			if f.tag != MaskerTypeMapStruct {
				dst.Set(sf)
				continue
			}

			maskedMap := reflect.MakeMapWithSize(sf.Type(), sf.Len())
			iter := sf.MapRange()
			for iter.Next() {
				maskedValue, err := m.maskMapStructValue(iter.Value())
				if err != nil {
					return nil, err
				}

				maskedMap.SetMapIndex(iter.Key(), maskedValue)
			}
			dst.Set(maskedMap)
			continue

		case reflect.Slice:
			if sf.IsNil() {
				continue
			}
			switch {
			case f.elemKind == reflect.String:
				orgval := sf.Interface().([]string)
				newval := make([]string, len(orgval))
				for i, val := range orgval {
					v, err := m.Marshal(f.tag, val)
					if err != nil {
						return nil, err
					}
					newval[i] = v
				}
				dst.Set(reflect.ValueOf(newval))
			case f.elemKind == reflect.Struct && f.tag == MaskerTypeStruct:
				newval := reflect.MakeSlice(sf.Type(), 0, sf.Len())
				for j, l := 0, sf.Len(); j < l; j++ {
					_n, err := m.Struct(sf.Index(j).Interface())
					if err != nil {
						return nil, err
					}
					newval = reflect.Append(newval, reflect.ValueOf(_n).Elem())
				}
				dst.Set(newval)
			case f.elemKind == reflect.Ptr && f.tag == MaskerTypeStruct:
				newval := reflect.MakeSlice(sf.Type(), 0, sf.Len())
				for j, l := 0, sf.Len(); j < l; j++ {
					if sf.Index(j).IsNil() {
						newval = reflect.Append(newval, sf.Index(j))
						continue
					}
					_n, err := m.Struct(sf.Index(j).Interface())
					if err != nil {
						return nil, err
					}
					newval = reflect.Append(newval, reflect.ValueOf(_n))
				}
				dst.Set(newval)
			case f.elemKind == reflect.Interface && f.tag == MaskerTypeStruct:
				newval := reflect.MakeSlice(sf.Type(), 0, sf.Len())
				for j, l := 0, sf.Len(); j < l; j++ {
					_n, err := m.Struct(sf.Index(j).Interface())
					if err != nil {
						return nil, err
					}
					if reflect.TypeOf(sf.Index(j).Interface()).Kind() != reflect.Ptr {
						newval = reflect.Append(newval, reflect.ValueOf(_n).Elem())
					} else {
						newval = reflect.Append(newval, reflect.ValueOf(_n))
					}
				}
				dst.Set(newval)
			default:
				dst.Set(sf)
			}
		case reflect.Interface:
			if sf.IsNil() {
				continue
			}
			if f.tag != MaskerTypeStruct {
				continue
			}
			_t, err := m.Struct(sf.Interface())
			if err != nil {
				return nil, err
			}
			if reflect.TypeOf(sf.Interface()).Kind() != reflect.Ptr {
				dst.Set(reflect.ValueOf(_t).Elem())
			} else {
				dst.Set(reflect.ValueOf(_t))
			}
		}
	}

	return tptr.Interface(), nil
}

// maskMapStructValue recursively masks mapstruct values for struct, ptr, and slice shapes.
// Unsupported scalar/container leaf values are returned as-is to preserve original data.
// Note: cyclic object graphs are not handled and may cause a stack overflow.
func (m *MaskerMarshaler) maskMapStructValue(v reflect.Value) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Map:
		if v.IsNil() {
			return v, nil
		}
		maskedMap := reflect.MakeMapWithSize(v.Type(), v.Len())
		iter := v.MapRange()
		for iter.Next() {
			maskedValue, err := m.maskMapStructValue(iter.Value())
			if err != nil {
				return reflect.Value{}, err
			}
			maskedMap.SetMapIndex(iter.Key(), maskedValue)
		}
		return maskedMap, nil
	case reflect.Struct:
		nextValue, err := m.Struct(v.Interface())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(nextValue).Elem(), nil
	case reflect.Ptr:
		if v.IsNil() {
			return v, nil
		}
		maskedElem, err := m.maskMapStructValue(v.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		maskedPtr := reflect.New(v.Elem().Type())
		maskedPtr.Elem().Set(maskedElem)
		return maskedPtr, nil
	case reflect.Slice:
		if v.IsNil() {
			return v, nil
		}
		maskedSlice := reflect.MakeSlice(v.Type(), 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			maskedElem, err := m.maskMapStructValue(v.Index(i))
			if err != nil {
				return reflect.Value{}, err
			}
			maskedSlice = reflect.Append(maskedSlice, maskedElem)
		}
		return maskedSlice, nil
	default:
		return v, nil
	}
}

// NewMaskerMarshaler returns a new masker marshaler
// It has default maskers and default masker
// Default maskers are:
//   - NoneMasker
//   - PasswordMasker
//   - NameMasker
//   - AddressMasker
//   - EmailMasker
//   - MobileMasker
//   - TelephoneMasker
//   - IDMasker
//   - CreditMasker
//   - URLMasker
//   - AbuseMasker
//   - AllMasker
//
// Dynamic tags "first-N" and "last-N" are also supported without registration.
// Default masker is "*"
// It is used for masking sensitive data
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
			MaskerTypeAbuse:    NewAbuseMasker(),
			MaskerTypeAll:      &AllMasker{},
		},
		masker: "*",
	}
}

// DefaultMaskerMarshaler is a default masker marshaler.
// It has the same default maskers as NewMaskerMarshaler().
var DefaultMaskerMarshaler = NewMaskerMarshaler()

func strLoop(str string, length int) string {
	return strings.Repeat(str, length)
}

func overlay(str string, overlay string, start int, end int) string {
	r := []rune(str)
	l := len(r)

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
		start, end = end, start
	}

	return string(r[:start]) + overlay + string(r[end:])
}
