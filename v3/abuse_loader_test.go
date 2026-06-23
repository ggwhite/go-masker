package masker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAbuseWordLoader_LoadFromString(t *testing.T) {
	loader := NewAbuseWordLoader()
	words, err := loader.LoadFromString("bad\n# comment\n\nterrible\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"bad", "terrible"}
	if strings.Join(words, ",") != strings.Join(want, ",") {
		t.Errorf("LoadFromString() = %v, want %v", words, want)
	}
}

func TestAbuseWordLoader_LoadFromReader(t *testing.T) {
	loader := NewAbuseWordLoader()
	words, err := loader.LoadFromReader(strings.NewReader("foo\nbar"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(words) != 2 || words[0] != "foo" || words[1] != "bar" {
		t.Errorf("LoadFromReader() = %v, want [foo bar]", words)
	}
}

func TestAbuseWordLoader_LoadFromSlice(t *testing.T) {
	loader := NewAbuseWordLoader()
	words := loader.LoadFromSlice([]string{" foo ", "", "bar"})
	if len(words) != 2 || words[0] != "foo" || words[1] != "bar" {
		t.Errorf("LoadFromSlice() = %v, want [foo bar]", words)
	}
}

func TestAbuseWordLoader_LoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "words.txt")
	if err := os.WriteFile(path, []byte("bad\nterrible\n"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	loader := NewAbuseWordLoader()
	words, err := loader.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(words) != 2 {
		t.Errorf("LoadFromFile() len = %d, want 2", len(words))
	}

	if _, err := loader.LoadFromFile(filepath.Join(dir, "missing.txt")); err == nil {
		t.Errorf("expected error for missing file")
	}
}
