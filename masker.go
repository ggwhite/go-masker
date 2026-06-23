package masker

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const tagName = "mask"

// MaskerType 是遮罩型別的字串別名，對應 struct tag `mask:"..."` 的值。
type MaskerType string

// MaskerType 常數。v3 採用精簡命名（保留 Type prefix 以避免與其他識別字衝突）。
const (
	TypeNone      MaskerType = "none"
	TypePassword  MaskerType = "password"
	TypeName      MaskerType = "name"
	TypeAddress   MaskerType = "addr"
	TypeEmail     MaskerType = "email"
	TypeMobile    MaskerType = "mobile"
	TypeTel       MaskerType = "tel"
	TypeID        MaskerType = "id"
	TypeCredit    MaskerType = "credit"
	TypeURL       MaskerType = "url"
	TypeAbuse     MaskerType = "abuse"
	TypeStruct    MaskerType = "struct"
	TypeAll       MaskerType = "all"
	TypeMapStruct MaskerType = "mapstruct"
)

// Masker 是遮罩單一字串的介面。
// 遮罩字元已於建構時注入各實作，呼叫者只需提供待遮罩的值。
// 實作自訂 masker 時，實作此介面並透過 MaskerMarshaler.Register 註冊即可。
type Masker interface {
	Mask(value string) string
}

// Option 是 MaskerMarshaler 的建構選項，遵循 functional options 模式。
type Option func(*MaskerMarshaler)

// WithMaskChar 設定此 marshaler 使用的遮罩字元，預設為 '*'。
// 設定後僅影響由此 marshaler 建立的 masker 實例，不會污染 DefaultMaskerMarshaler。
// 注意：URLMasker（密碼段固定輸出 xxxxx）與 NoneMasker（原樣回傳）不受此設定影響。
func WithMaskChar(c rune) Option {
	return func(m *MaskerMarshaler) {
		m.maskChar = string(c)
	}
}

// MaskerMarshaler 管理一組 masker 並提供 Marshal 與 Struct 等遮罩操作。
// 所有 exported 方法皆可安全並行使用。
type MaskerMarshaler struct {
	mu       sync.RWMutex
	maskers  map[MaskerType]Masker
	maskChar string
}

// Marshal 依 masker 型別回傳遮罩後的值，找不到型別時回傳 error。
// 同時支援動態 tag first-N / last-N，這類 tag 不需事先註冊。
// Example:
//
//	m := masker.NewMaskerMarshaler()
//	s, _ := m.Marshal(masker.TypeMobile, "0987654321") // 0987***321
func (m *MaskerMarshaler) Marshal(t MaskerType, value string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if result, matched, err := parseGenericMask(m.maskChar, string(t), value); matched {
		return result, err
	}
	masker, ok := m.maskers[t]
	if !ok {
		return "", fmt.Errorf("masker %v not found", t)
	}
	return masker.Mask(value), nil
}

// MustMarshal 與 Marshal 相同，但找不到 masker 型別時直接 panic。
// 適用於呼叫者確信型別存在、不想處理 error 的場景。
// Example:
//
//	m := masker.NewMaskerMarshaler()
//	s := m.MustMarshal(masker.TypeMobile, "0987654321") // 0987***321
func (m *MaskerMarshaler) MustMarshal(t MaskerType, value string) string {
	result, err := m.Marshal(t, value)
	if err != nil {
		panic(err)
	}
	return result
}

// Register 以 masker 型別新增或覆寫一個 masker。
func (m *MaskerMarshaler) Register(t MaskerType, masker Masker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.maskers[t] = masker
}

// Unregister 移除指定型別的 masker。
func (m *MaskerMarshaler) Unregister(t MaskerType) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.maskers, t)
}

// Get 回傳指定型別的 masker，找不到時回傳 error。
func (m *MaskerMarshaler) Get(t MaskerType) (Masker, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	masker, ok := m.maskers[t]
	if !ok {
		return nil, fmt.Errorf("masker %v not found", t)
	}
	return masker, nil
}

// List 回傳目前已註冊的所有 masker 型別。
func (m *MaskerMarshaler) List() []MaskerType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []MaskerType
	for t := range m.maskers {
		list = append(list, t)
	}
	return list
}

// Struct 接收 struct 或 struct 指標，依欄位上的 `mask` tag 遮罩後回傳同型別的指標。
// 原輸入不會被修改。
//
// Limitation:
//   - 不特別處理循環物件圖；自我參照結構在 `mask:"struct"` 或 `mask:"mapstruct"` 下可能無限遞迴。
//
// Example:
//
//	type Foo struct {
//		Name  string `mask:"name"`
//		Email string `mask:"email"`
//	}
//	m := masker.NewMaskerMarshaler()
//	out, _ := m.Struct(&Foo{Name: "John Doe", Email: "john@gmail.com"})
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

	for _, f := range cachedTypeInfo(selem.Type()).fields {
		if !f.exported {
			continue
		}
		if !f.hasTag {
			tptr.Elem().Field(f.index).Set(selem.Field(f.index))
			continue
		}

		switch f.kind {
		default:
			tptr.Elem().Field(f.index).Set(selem.Field(f.index))
		case reflect.String:
			v, err := m.Marshal(f.tag, selem.Field(f.index).String())
			if err != nil {
				return nil, err
			}
			tptr.Elem().Field(f.index).SetString(v)
		case reflect.Struct:
			if f.tag == TypeStruct {
				_t, err := m.Struct(selem.Field(f.index).Interface())
				if err != nil {
					return nil, err
				}
				tptr.Elem().Field(f.index).Set(reflect.ValueOf(_t).Elem())
			} else {
				tptr.Elem().Field(f.index).Set(selem.Field(f.index))
			}
		case reflect.Ptr:
			if selem.Field(f.index).IsNil() {
				continue
			}
			if f.tag == TypeStruct {
				_t, err := m.Struct(selem.Field(f.index).Interface())
				if err != nil {
					return nil, err
				}
				tptr.Elem().Field(f.index).Set(reflect.ValueOf(_t))
			} else {
				tptr.Elem().Field(f.index).Set(selem.Field(f.index))
			}
		case reflect.Map:
			if selem.Field(f.index).IsNil() {
				continue
			}
			if f.tag != TypeMapStruct {
				tptr.Elem().Field(f.index).Set(selem.Field(f.index))
				continue
			}

			maskedMap := reflect.MakeMapWithSize(selem.Field(f.index).Type(), selem.Field(f.index).Len())
			iter := selem.Field(f.index).MapRange()
			for iter.Next() {
				v := iter.Value()
				maskedValue, err := m.maskMapStructValue(v)
				if err != nil {
					return nil, err
				}

				maskedMap.SetMapIndex(iter.Key(), maskedValue)
			}
			tptr.Elem().Field(f.index).Set(maskedMap)
			continue

		case reflect.Slice:
			if selem.Field(f.index).IsNil() {
				continue
			}
			switch {
			case f.elemKind == reflect.String:
				orgval := selem.Field(f.index).Interface().([]string)
				newval := make([]string, len(orgval))
				for j, val := range orgval {
					v, err := m.Marshal(f.tag, val)
					if err != nil {
						return nil, err
					}
					newval[j] = v
				}
				tptr.Elem().Field(f.index).Set(reflect.ValueOf(newval))
			case f.elemKind == reflect.Struct && f.tag == TypeStruct:
				newval := reflect.MakeSlice(selem.Field(f.index).Type(), 0, selem.Field(f.index).Len())
				for j, l := 0, selem.Field(f.index).Len(); j < l; j++ {
					_n, err := m.Struct(selem.Field(f.index).Index(j).Interface())
					if err != nil {
						return nil, err
					}
					newval = reflect.Append(newval, reflect.ValueOf(_n).Elem())
				}
				tptr.Elem().Field(f.index).Set(newval)
			case f.elemKind == reflect.Ptr && f.tag == TypeStruct:
				newval := reflect.MakeSlice(selem.Field(f.index).Type(), 0, selem.Field(f.index).Len())
				for j, l := 0, selem.Field(f.index).Len(); j < l; j++ {
					if selem.Field(f.index).Index(j).IsNil() {
						newval = reflect.Append(newval, selem.Field(f.index).Index(j))
						continue
					}
					_n, err := m.Struct(selem.Field(f.index).Index(j).Interface())
					if err != nil {
						return nil, err
					}
					newval = reflect.Append(newval, reflect.ValueOf(_n))
				}
				tptr.Elem().Field(f.index).Set(newval)
			case f.elemKind == reflect.Interface && f.tag == TypeStruct:
				newval := reflect.MakeSlice(selem.Field(f.index).Type(), 0, selem.Field(f.index).Len())
				for j, l := 0, selem.Field(f.index).Len(); j < l; j++ {
					_n, err := m.Struct(selem.Field(f.index).Index(j).Interface())
					if err != nil {
						return nil, err
					}
					if reflect.TypeOf(selem.Field(f.index).Index(j).Interface()).Kind() != reflect.Ptr {
						newval = reflect.Append(newval, reflect.ValueOf(_n).Elem())
					} else {
						newval = reflect.Append(newval, reflect.ValueOf(_n))
					}
				}
				tptr.Elem().Field(f.index).Set(newval)
			default:
				tptr.Elem().Field(f.index).Set(selem.Field(f.index))
			}
		case reflect.Interface:
			if selem.Field(f.index).IsNil() {
				continue
			}
			if f.tag == TypeStruct {
				_t, err := m.Struct(selem.Field(f.index).Interface())
				if err != nil {
					return nil, err
				}
				if reflect.TypeOf(selem.Field(f.index).Interface()).Kind() != reflect.Ptr {
					tptr.Elem().Field(f.index).Set(reflect.ValueOf(_t).Elem())
				} else {
					tptr.Elem().Field(f.index).Set(reflect.ValueOf(_t))
				}
			} else {
				concrete := reflect.ValueOf(selem.Field(f.index).Interface())
				if concrete.Kind() == reflect.String {
					v, err := m.Marshal(f.tag, concrete.String())
					if err != nil {
						return nil, err
					}
					tptr.Elem().Field(f.index).Set(reflect.ValueOf(v))
				} else {
					tptr.Elem().Field(f.index).Set(selem.Field(f.index))
				}
			}
		}
	}

	return tptr.Interface(), nil
}

// maskMapStructValue 遞迴遮罩 mapstruct 的值，支援 struct、ptr、slice、map 形狀。
// 不支援的純量／容器葉節點原樣回傳以保留原始資料。
// 注意：不處理循環物件圖，可能造成 stack overflow。
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
	case reflect.Interface:
		if v.IsNil() {
			return v, nil
		}
		return m.maskMapStructValue(v.Elem())
	default:
		return v, nil
	}
}

// NewMaskerMarshaler 建立一個新的 MaskerMarshaler，並註冊全部預設 masker。
// 預設遮罩字元為 '*'，可透過 WithMaskChar 選項變更。
// 每次呼叫都會建立各自獨立的一組 masker 實例，因此設定 WithMaskChar 不會影響其他 marshaler。
// 動態 tag first-N / last-N 不需註冊即可使用。
func NewMaskerMarshaler(opts ...Option) *MaskerMarshaler {
	m := &MaskerMarshaler{maskChar: "*"}
	for _, opt := range opts {
		opt(m)
	}
	c := m.maskChar
	m.maskers = map[MaskerType]Masker{
		TypeNone:     &NoneMasker{},
		TypePassword: &PasswordMasker{mask: c},
		TypeName:     &NameMasker{mask: c},
		TypeAddress:  &AddressMasker{mask: c},
		TypeEmail:    &EmailMasker{mask: c},
		TypeMobile:   &MobileMasker{mask: c},
		TypeTel:      &TelephoneMasker{mask: c},
		TypeID:       &IDMasker{mask: c},
		TypeCredit:   &CreditMasker{mask: c},
		TypeURL:      &URLMasker{},
		TypeAbuse:    NewAbuseMasker(c),
		TypeAll:      &AllMasker{mask: c},
	}
	return m
}

// DefaultMaskerMarshaler 是預設的 MaskerMarshaler，遮罩字元為 '*'，供 package-level 便利函式使用。
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
