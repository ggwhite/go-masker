package masker

import "math"

// AddressMasker is a masker for address
type AddressMasker struct{}

// Marshal masks address
// It masks last 6 digits of address
// Example:
//
//	AddressMasker{}.Marshal("*", "台北市內湖區內湖路一段737巷1號1樓") // returns "台北市內湖區******"
func (m *AddressMasker) Marshal(s string, i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	if l <= 6 {
		return strLoop(s, 6)
	}
	return overlay(i, strLoop(s, 6), 6, math.MaxInt)
}
