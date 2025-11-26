package masker

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// AbuseWordLoader provides methods to load abuse words from different sources
type AbuseWordLoader struct{}

// NewAbuseWordLoader creates a new abuse word loader
func NewAbuseWordLoader() *AbuseWordLoader {
	return &AbuseWordLoader{}
}

// LoadFromFile loads abuse words from a local file
// Each word should be on a separate line
// Empty lines and lines starting with # are ignored (comments)
func (l *AbuseWordLoader) LoadFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return l.loadFromReader(file)
}

// LoadFromReader loads abuse words from an io.Reader
// Each word should be on a separate line
// Empty lines and lines starting with # are ignored (comments)
func (l *AbuseWordLoader) LoadFromReader(reader io.Reader) ([]string, error) {
	return l.loadFromReader(reader)
}

// LoadFromString loads abuse words from a string
// Each word should be on a separate line
// Empty lines and lines starting with # are ignored (comments)
func (l *AbuseWordLoader) LoadFromString(content string) ([]string, error) {
	reader := strings.NewReader(content)
	return l.loadFromReader(reader)
}

// LoadFromSlice loads abuse words from a string slice
// This is useful when you already have words in memory
func (l *AbuseWordLoader) LoadFromSlice(words []string) []string {
	var result []string
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" {
			result = append(result, word)
		}
	}
	return result
}

// loadFromReader is the internal method that processes the reader
func (l *AbuseWordLoader) loadFromReader(reader io.Reader) ([]string, error) {
	var words []string
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		words = append(words, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
