package masker

import (
	"strings"
	"testing"
)

func TestAbuseTrie(t *testing.T) {
	trie := NewAbuseTrie()

	// Test empty trie
	if trie.Contains("test") {
		t.Error("Empty trie should not contain any words")
	}

	// Test inserting words
	words := []string{"bad", "terrible", "awful", "horrible"}
	trie.InsertAll(words)

	// Test contains
	for _, word := range words {
		if !trie.Contains(word) {
			t.Errorf("Trie should contain word: %s", word)
		}
	}

	// Test case insensitive
	if !trie.Contains("BAD") {
		t.Error("Trie should be case insensitive")
	}

	// Test non-existent words
	nonExistent := []string{"good", "nice", "wonderful"}
	for _, word := range nonExistent {
		if trie.Contains(word) {
			t.Errorf("Trie should not contain word: %s", word)
		}
	}
}

func TestAbuseMasker(t *testing.T) {
	// Test with predefined words
	abuseWords := []string{"bad", "terrible", "awful"}
	masker := NewAbuseMaskerWithWords(abuseWords)

	// Test basic masking
	text := "This is a bad word and terrible"
	expected := "This is a *** word and ********"
	result := masker.Marshal("*", text)

	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}

	// Test case insensitive
	text = "This is a BAD word and TERRIBLE"
	expected = "This is a *** word and ********"
	result = masker.Marshal("*", text)

	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}

	// Test with different mask character
	text = "This is a bad word"
	expected = "This is a ### word"
	result = masker.Marshal("#", text)

	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}

	// Test empty text
	result = masker.Marshal("*", "")
	if result != "" {
		t.Errorf("Expected empty string, Got: %s", result)
	}

	// Test text with no abuse words
	text = "This is a good word"
	result = masker.Marshal("*", text)
	if result != text {
		t.Errorf("Expected: %s, Got: %s", text, result)
	}
}

func TestAbuseMaskerAddWords(t *testing.T) {
	masker := NewAbuseMasker()

	// Add words after creation
	words := []string{"bad", "terrible"}
	masker.AddWords(words)

	text := "This is a bad word"
	expected := "This is a *** word"
	result := masker.Marshal("*", text)

	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}

	// Add single word
	masker.AddWord("awful")
	text = "This is an awful word"
	expected = "This is an ***** word"
	result = masker.Marshal("*", text)

	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
}

func TestAbuseMaskerContainsAbuse(t *testing.T) {
	masker := NewAbuseMaskerWithWords([]string{"bad", "terrible"})

	// Test contains abuse
	text := "This is a bad word"
	if !masker.ContainsAbuse(text) {
		t.Error("Should detect abuse in text")
	}

	// Test no abuse
	text = "This is a good word"
	if masker.ContainsAbuse(text) {
		t.Error("Should not detect abuse in clean text")
	}
}

func TestAbuseMaskerGetAbuseWords(t *testing.T) {
	masker := NewAbuseMaskerWithWords([]string{"bad", "terrible", "awful"})

	text := "This is a bad word and terrible"
	abuseWords := masker.GetAbuseWords(text)

	expected := []string{"bad", "terrible"}
	if len(abuseWords) != len(expected) {
		t.Errorf("Expected %d abuse words, Got %d", len(expected), len(abuseWords))
	}

	for _, word := range expected {
		found := false
		for _, foundWord := range abuseWords {
			if foundWord == word {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find abuse word: %s", word)
		}
	}
}

func TestAbuseMaskerMultipleOccurrences(t *testing.T) {
	masker := NewAbuseMaskerWithWords([]string{"bad"})

	// Test multiple occurrences
	text := "This is bad and that is also bad"
	expected := "This is *** and that is also ***"
	result := masker.Marshal("*", text)

	if result != expected {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
}

func TestAbuseMaskerEmptyTrie(t *testing.T) {
	masker := NewAbuseMasker() // Empty trie

	text := "This is a test"
	result := masker.Marshal("*", text)

	if result != text {
		t.Errorf("Expected: %s, Got: %s", text, result)
	}

	if masker.ContainsAbuse(text) {
		t.Error("Empty trie should not detect abuse")
	}

	abuseWords := masker.GetAbuseWords(text)
	if len(abuseWords) != 0 {
		t.Error("Empty trie should return no abuse words")
	}
}

func TestCleanWord(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"bad", "bad"},
		{"BAD", "bad"},
		{"bad!", "bad"},
		{"bad.", "bad"},
		{"bad,", "bad"},
		{"bad;", "bad"},
		{"bad:", "bad"},
		{"bad?", "bad"},
		{"bad!", "bad"},
		{"bad@", "bad"},
		{"bad#", "bad"},
		{"bad$", "bad"},
		{"bad%", "bad"},
		{"bad^", "bad"},
		{"bad&", "bad"},
		{"bad*", "bad"},
		{"bad(", "bad"},
		{"bad)", "bad"},
		{"bad-", "bad"},
		{"bad_", "bad"},
		{"bad+", "bad"},
		{"bad=", "bad"},
		{"bad[", "bad"},
		{"bad]", "bad"},
		{"bad{", "bad"},
		{"bad}", "bad"},
		{"bad|", "bad"},
		{"bad\\", "bad"},
		{"bad/", "bad"},
		{"bad<", "bad"},
		{"bad>", "bad"},
		{"bad`", "bad"},
		{"bad~", "bad"},
		{"bad123", "bad123"},
		{"123bad", "123bad"},
		{"", ""},
		{"   ", ""},
		{"   bad   ", "bad"},
	}

	for _, test := range tests {
		result := cleanWord(test.input)
		if result != test.expected {
			t.Errorf("cleanWord(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func BenchmarkAbuseMasker(b *testing.B) {
	// Create a large list of abuse words
	abuseWords := make([]string, 5000)
	for i := 0; i < 1000; i++ {
		abuseWords[i] = strings.Repeat("a", i%10+1) + strings.Repeat("b", i%5+1)
	}

	masker := NewAbuseMaskerWithWords(abuseWords)
	text := "This is a test with some words that might be abuse words"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		masker.Marshal("*", text)
	}
}

func TestAbuseMasker_NoSubstringMatch(t *testing.T) {
	masker := NewAbuseMaskerWithWords([]string{"bad"})
	got := masker.Marshal("*", "bad badminton")
	want := "*** badminton"
	if got != want {
		t.Errorf("got = %q, want %q", got, want)
	}
}

func TestAbuseMasker_PreservesWhitespace(t *testing.T) {
	masker := NewAbuseMaskerWithWords([]string{"bad"})
	got := masker.Marshal("*", "  bad  good  bad  ")
	want := "  ***  good  ***  "
	if got != want {
		t.Errorf("got = %q, want %q", got, want)
	}
}
