package masker

// NameMasker is a masker for name
type NoneMasker struct{}

// Marshal masks name
// It returns the same value
// Example:
//
//	NoneMasker{}.Marshal("*", "name") // returns "name"
func (m *NoneMasker) Marshal(i string, value string) string {
	return value
}
