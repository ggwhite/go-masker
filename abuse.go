package masker

import (
	"strings"
	"unicode"
)

// TrieNode represents a node in the trie data structure
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// AbuseTrie represents a trie data structure for storing abuse words
type AbuseTrie struct {
	root *TrieNode
}

// NewAbuseTrie creates a new empty trie
func NewAbuseTrie() *AbuseTrie {
	return &AbuseTrie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

// Insert adds a word to the trie
func (t *AbuseTrie) Insert(word string) {
	if word == "" {
		return
	}

	node := t.root
	word = strings.ToLower(strings.TrimSpace(word))

	for _, char := range word {
		if node.children[char] == nil {
			node.children[char] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// InsertAll adds multiple words to the trie
func (t *AbuseTrie) InsertAll(words []string) {
	for _, word := range words {
		t.Insert(word)
	}
}

// Contains checks if a word exists in the trie
func (t *AbuseTrie) Contains(word string) bool {
	if word == "" {
		return false
	}

	node := t.root
	word = strings.ToLower(strings.TrimSpace(word))

	for _, char := range word {
		if node.children[char] == nil {
			return false
		}
		node = node.children[char]
	}
	return node.isEnd
}

// findAbuseWords finds all abuse words in the given text
func (t *AbuseTrie) findAbuseWords(text string) []string {
	var found []string
	words := strings.Fields(text)

	for _, word := range words {
		// Clean the word (remove punctuation)
		cleanWord := cleanWord(word)
		if cleanWord != "" && t.Contains(cleanWord) {
			found = append(found, word)
		}
	}

	return found
}

// cleanWord removes punctuation and converts to lowercase
func cleanWord(word string) string {
	var result strings.Builder
	for _, char := range word {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			result.WriteRune(unicode.ToLower(char))
		}
	}
	return result.String()
}

// AbuseMasker is a masker for abuse words
type AbuseMasker struct {
	trie *AbuseTrie
}

// NewAbuseMasker creates a new abuse masker with an empty trie
func NewAbuseMasker() *AbuseMasker {
	return &AbuseMasker{
		trie: NewAbuseTrie(),
	}
}

// NewAbuseMaskerWithWords creates a new abuse masker with predefined words
func NewAbuseMaskerWithWords(words []string) *AbuseMasker {
	trie := NewAbuseTrie()
	trie.InsertAll(words)
	return &AbuseMasker{
		trie: trie,
	}
}

// AddWords adds abuse words to the masker
func (m *AbuseMasker) AddWords(words []string) {
	m.trie.InsertAll(words)
}

// AddWord adds a single abuse word to the masker
func (m *AbuseMasker) AddWord(word string) {
	m.trie.Insert(word)
}

// Marshal masks abuse words in the given text
// It replaces abuse words with the specified mask character
// Example:
//
//	abuseMasker := NewAbuseMaskerWithWords([]string{"bad", "terrible"})
//	abuseMasker.Marshal("*", "This is a bad word and terrible") // returns "This is a *** word and ********"
func (m *AbuseMasker) Marshal(maskChar string, text string) string {
	if text == "" {
		return text
	}

	// Find all abuse words in the text
	abuseWords := m.trie.findAbuseWords(text)

	// Replace each abuse word with masked version
	result := text
	for _, abuseWord := range abuseWords {
		maskedWord := strLoop(maskChar, len([]rune(abuseWord)))
		result = strings.ReplaceAll(result, abuseWord, maskedWord)
	}

	return result
}

// ContainsAbuse checks if the text contains any abuse words
func (m *AbuseMasker) ContainsAbuse(text string) bool {
	abuseWords := m.trie.findAbuseWords(text)
	return len(abuseWords) > 0
}

// GetAbuseWords returns all abuse words found in the text
func (m *AbuseMasker) GetAbuseWords(text string) []string {
	return m.trie.findAbuseWords(text)
}
