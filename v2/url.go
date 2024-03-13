package masker

import "net/url"

// URLMasker is a masker for URL
type URLMasker struct{}

// Marshal masks URL
// It mask the password part of the URL if exists
func (m *URLMasker) Marshal(s, i string) string {
	u, err := url.Parse(i)
	if err != nil {
		return i
	}
	return u.Redacted()
}
