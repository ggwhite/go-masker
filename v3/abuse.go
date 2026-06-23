package masker

import (
	"strings"
	"unicode"
)

// TrieNode 是 trie 資料結構中的節點。
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// AbuseTrie 是儲存髒話／敏感詞的 trie 資料結構。
type AbuseTrie struct {
	root *TrieNode
}

// NewAbuseTrie 建立一個空的 trie。
func NewAbuseTrie() *AbuseTrie {
	return &AbuseTrie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

// Insert 將一個詞加入 trie，加入前會轉小寫並去除前後空白。
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

// InsertAll 將多個詞一次加入 trie。
func (t *AbuseTrie) InsertAll(words []string) {
	for _, word := range words {
		t.Insert(word)
	}
}

// Contains 檢查某個詞是否存在於 trie。
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

// findAbuseWords 找出文字中所有命中字典的詞。
func (t *AbuseTrie) findAbuseWords(text string) []string {
	var found []string
	words := strings.Fields(text)

	for _, word := range words {
		cleanWord := cleanWord(word)
		if cleanWord != "" && t.Contains(cleanWord) {
			found = append(found, word)
		}
	}

	return found
}

// cleanWord 去除標點並轉小寫，只保留字母與數字。
func cleanWord(word string) string {
	var result strings.Builder
	for _, char := range word {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			result.WriteRune(unicode.ToLower(char))
		}
	}
	return result.String()
}

// AbuseMasker 是以字典命中後遮罩髒話／敏感詞的 masker。
// 須先載入詞典（透過 AddWords / AddWord 或 NewAbuseMaskerWithWords）才會遮罩，空字典時原樣回傳。
type AbuseMasker struct {
	trie *AbuseTrie
	mask string
}

// NewAbuseMasker 建立一個空字典的 abuse masker，並指定遮罩字元。
func NewAbuseMasker(maskChar string) *AbuseMasker {
	return &AbuseMasker{
		trie: NewAbuseTrie(),
		mask: maskChar,
	}
}

// NewAbuseMaskerWithWords 建立一個預先載入詞典的 abuse masker，並指定遮罩字元。
func NewAbuseMaskerWithWords(maskChar string, words []string) *AbuseMasker {
	trie := NewAbuseTrie()
	trie.InsertAll(words)
	return &AbuseMasker{
		trie: trie,
		mask: maskChar,
	}
}

// AddWords 將多個敏感詞加入 masker。
func (m *AbuseMasker) AddWords(words []string) {
	m.trie.InsertAll(words)
}

// AddWord 將單一敏感詞加入 masker。
func (m *AbuseMasker) AddWord(word string) {
	m.trie.Insert(word)
}

// Mask 遮罩文字中命中字典的敏感詞。
// 只有以空白分隔的完整 token 會被遮罩，較長字內的子字串不會被遮罩。
// 空字典時原樣回傳。
// Example:
//
//	m := NewAbuseMaskerWithWords("*", []string{"bad", "terrible"})
//	m.Mask("This is a bad word and terrible") // returns "This is a *** word and ********"
func (m *AbuseMasker) Mask(text string) string {
	if text == "" {
		return text
	}

	runes := []rune(text)
	var result strings.Builder
	result.Grow(len(text))
	i := 0

	for i < len(runes) {
		if unicode.IsSpace(runes[i]) {
			result.WriteRune(runes[i])
			i++
			continue
		}

		j := i
		for j < len(runes) && !unicode.IsSpace(runes[j]) {
			j++
		}

		word := string(runes[i:j])
		cleaned := cleanWord(word)
		if cleaned != "" && m.trie.Contains(cleaned) {
			result.WriteString(strLoop(m.mask, j-i))
		} else {
			result.WriteString(word)
		}
		i = j
	}

	return result.String()
}

// ContainsAbuse 檢查文字是否含有任何命中字典的敏感詞。
func (m *AbuseMasker) ContainsAbuse(text string) bool {
	abuseWords := m.trie.findAbuseWords(text)
	return len(abuseWords) > 0
}

// GetAbuseWords 回傳文字中所有命中字典的敏感詞。
func (m *AbuseMasker) GetAbuseWords(text string) []string {
	return m.trie.findAbuseWords(text)
}
