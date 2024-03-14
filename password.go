package masker

// PasswordMasker is a masker for password
type PasswordMasker struct{}

// Marshal masks password
// It returns 14 asterisks
// Example:
//
//	PasswordMasker{}.Marshal("*", "password") // returns "**************"
//	PasswordMasker{}.Marshal("&", "password") // returns "&&&&&&&&&&&&&&"
func (m *PasswordMasker) Marshal(s string, i string) string {
	return strLoop(s, 14)
}
