package masker

// NoneMasker 是不做任何遮罩、原樣回傳值的 masker。
type NoneMasker struct{}

// Mask 原樣回傳輸入值，不做任何遮罩。
// Example:
//
//	(&NoneMasker{}).Mask("secret") // returns "secret"
func (m *NoneMasker) Mask(value string) string {
	return value
}
