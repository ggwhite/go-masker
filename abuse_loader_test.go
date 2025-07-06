package masker

import (
	"strings"
	"testing"
)

func TestAbuseWordLoaderLoadFromString(t *testing.T) {
	loader := NewAbuseWordLoader()

	content := `# This is a comment
bad
terrible
# Another comment

awful
horrible
`

	words, err := loader.LoadFromString(content)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expected := []string{"bad", "terrible", "awful", "horrible"}
	if len(words) != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), len(words))
	}

	for i, word := range expected {
		if words[i] != word {
			t.Errorf("Expected word %s at position %d, got %s", word, i, words[i])
		}
	}
}

func TestAbuseWordLoaderLoadFromSlice(t *testing.T) {
	loader := NewAbuseWordLoader()

	input := []string{"bad", "  terrible  ", "", "awful", "   "}
	words := loader.LoadFromSlice(input)

	expected := []string{"bad", "terrible", "awful"}
	if len(words) != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), len(words))
	}

	for i, word := range expected {
		if words[i] != word {
			t.Errorf("Expected word %s at position %d, got %s", word, i, words[i])
		}
	}
}

func TestAbuseWordLoaderLoadFromReader(t *testing.T) {
	loader := NewAbuseWordLoader()

	content := "bad\nterrible\nawful"
	reader := strings.NewReader(content)

	words, err := loader.LoadFromReader(reader)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expected := []string{"bad", "terrible", "awful"}
	if len(words) != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), len(words))
	}

	for i, word := range expected {
		if words[i] != word {
			t.Errorf("Expected word %s at position %d, got %s", word, i, words[i])
		}
	}
}

func TestAbuseWordLoaderEmptyInput(t *testing.T) {
	loader := NewAbuseWordLoader()

	// Test empty string
	words, err := loader.LoadFromString("")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(words) != 0 {
		t.Errorf("Expected 0 words, got %d", len(words))
	}

	// Test empty slice
	words = loader.LoadFromSlice([]string{})
	if len(words) != 0 {
		t.Errorf("Expected 0 words, got %d", len(words))
	}

	// Test slice with only empty strings
	words = loader.LoadFromSlice([]string{"", "   ", "\t", "\n"})
	if len(words) != 0 {
		t.Errorf("Expected 0 words, got %d", len(words))
	}
}

func TestAbuseWordLoaderOnlyComments(t *testing.T) {
	loader := NewAbuseWordLoader()

	content := `# This is a comment
# Another comment
# Yet another comment
`

	words, err := loader.LoadFromString(content)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(words) != 0 {
		t.Errorf("Expected 0 words, got %d", len(words))
	}
}

func TestAbuseWordLoaderMixedContent(t *testing.T) {
	loader := NewAbuseWordLoader()

	content := `# Header comment
bad
# Middle comment
terrible
# End comment
awful
`

	words, err := loader.LoadFromString(content)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expected := []string{"bad", "terrible", "awful"}
	if len(words) != len(expected) {
		t.Errorf("Expected %d words, got %d", len(expected), len(words))
	}

	for i, word := range expected {
		if words[i] != word {
			t.Errorf("Expected word %s at position %d, got %s", word, i, words[i])
		}
	}
}

func BenchmarkAbuseWordLoaderLoadFromString(b *testing.B) {
	loader := NewAbuseWordLoader()
	content := strings.Repeat("word\n", 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loader.LoadFromString(content)
	}
}

func BenchmarkAbuseWordLoaderLoadFromSlice(b *testing.B) {
	loader := NewAbuseWordLoader()
	words := make([]string, 5000)
	for i := 0; i < 5000; i++ {
		words[i] = "word"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loader.LoadFromSlice(words)
	}
}
