package masker

import "github.com/ggwhite/go-masker/cmd"

// Masker mask something
type Masker interface {
	Struct(interface{}, interface{}) error
	Name(string) string
	ID(string) string
	Address(string) string
	Credit(string) string
	Email(string) string
	Mobile(string) string
	Telephone(string) string
}

type instance struct{}

func (*instance) Struct(s, t interface{}) error {
	return cmd.Struct(s, t)
}

func (*instance) Name(i string) string {
	return cmd.Name(i)
}

func (*instance) ID(i string) string {
	return cmd.ID(i)
}

func (*instance) Address(i string) string {
	return cmd.Address(i)
}

func (*instance) Credit(i string) string {
	return cmd.Credit(i)
}

func (*instance) Email(i string) string {
	return cmd.Email(i)
}

func (*instance) Mobile(i string) string {
	return cmd.Mobile(i)
}

func (*instance) Telephone(i string) string {
	return cmd.Telephone(i)
}

// New create Masker
func New() Masker {
	return &instance{}
}

var defaultmasker Masker

func init() {
	defaultmasker = New()
}

// Struct mask target from source
func Struct(s, t interface{}) error {
	return defaultmasker.Struct(s, t)
}

// Name mask name of middle
func Name(i string) string {
	return defaultmasker.Name(i)
}

// ID mask ID Nbr of middle
func ID(i string) string {
	return defaultmasker.ID(i)
}

// Address mask address
func Address(i string) string {
	return defaultmasker.Address(i)
}

// Credit mask credit card number
func Credit(i string) string {
	return defaultmasker.Credit(i)
}

// Email mask email
func Email(i string) string {
	return defaultmasker.Email(i)
}

// Mobile mask mobile
func Mobile(i string) string {
	return defaultmasker.Mobile(i)
}

// Telephone mask telephone
func Telephone(i string) string {
	return cmd.Telephone(i)
}
