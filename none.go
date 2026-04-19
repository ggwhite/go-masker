package masker

// NoneMasker is a masker that returns the value unchanged.
type NoneMasker struct{}

// Marshal returns the value unchanged.
//
// Example:
//
//	NoneMasker{}.Marshal("*", "secret") // returns "secret"
func (m *NoneMasker) Marshal(i string, value string) string {
	return value
}
