package masker

import "fmt"

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

func NewMaskerMarshaler() *MaskerMarshaler {
	return &MaskerMarshaler{
		Maskers: make(map[MaskerType]Masker),
		masker:  "*",
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
