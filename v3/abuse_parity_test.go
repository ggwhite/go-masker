package masker

import "testing"

// 本檔移植 v2 root package 的 abuse 子系統 parity 測試（對照 abuse_test.go），
// API 名稱適配 v3：NewAbuseMasker(maskChar) / NewAbuseMaskerWithWords(maskChar, words)，
// 遮罩入口由 v2 的 Marshal(maskChar, text) 改為 v3 的 Mask(text)（mask 字元於建構時注入）。

// TestAbuseTrie_Full 對照 v2 TestAbuseTrie：Insert/InsertAll/Contains/大小寫不敏感/不存在詞。
func TestAbuseTrie_Full(t *testing.T) {
	trie := NewAbuseTrie()

	if trie.Contains("test") {
		t.Error("空 trie 不應命中任何詞")
	}

	words := []string{"bad", "terrible", "awful", "horrible"}
	trie.InsertAll(words)

	for _, word := range words {
		if !trie.Contains(word) {
			t.Errorf("trie 應包含詞: %s", word)
		}
	}

	if !trie.Contains("BAD") {
		t.Error("trie 應大小寫不敏感")
	}

	trie.Insert("single")
	if !trie.Contains("single") {
		t.Error("Insert 後應命中 single")
	}

	for _, word := range []string{"good", "nice", "wonderful"} {
		if trie.Contains(word) {
			t.Errorf("trie 不應包含詞: %s", word)
		}
	}
}

// TestAbuseTrie_FindAbuseWords 對照 v2 findAbuseWords 行為。
func TestAbuseTrie_FindAbuseWords(t *testing.T) {
	trie := NewAbuseTrie()
	trie.InsertAll([]string{"bad", "terrible"})
	got := trie.findAbuseWords("This is a bad and terrible day")
	if len(got) != 2 || got[0] != "bad" || got[1] != "terrible" {
		t.Errorf("findAbuseWords() = %v, want [bad terrible]", got)
	}
	if found := trie.findAbuseWords("all good here"); len(found) != 0 {
		t.Errorf("findAbuseWords(clean) = %v, want []", found)
	}
}

// TestAbuseMasker_WithWords 對照 v2 TestAbuseMasker：基本遮罩、大小寫、不同 mask 字元、空字串、無命中。
func TestAbuseMasker_WithWords(t *testing.T) {
	abuseWords := []string{"bad", "terrible", "awful"}
	m := NewAbuseMaskerWithWords("*", abuseWords)

	if got := m.Mask("This is a bad word and terrible"); got != "This is a *** word and ********" {
		t.Errorf("basic mask = %q, want %q", got, "This is a *** word and ********")
	}

	if got := m.Mask("This is a BAD word and TERRIBLE"); got != "This is a *** word and ********" {
		t.Errorf("case-insensitive mask = %q, want %q", got, "This is a *** word and ********")
	}

	// 不同 mask 字元（v3 於建構時注入）
	mHash := NewAbuseMaskerWithWords("#", abuseWords)
	if got := mHash.Mask("This is a bad word"); got != "This is a ### word" {
		t.Errorf("custom mask char = %q, want %q", got, "This is a ### word")
	}

	if got := m.Mask(""); got != "" {
		t.Errorf("empty text = %q, want empty", got)
	}

	if got := m.Mask("This is a good word"); got != "This is a good word" {
		t.Errorf("no-abuse text = %q, want unchanged", got)
	}
}

// TestAbuseMasker_AddWordsAndAddWord 對照 v2 TestAbuseMaskerAddWords。
func TestAbuseMasker_AddWordsAndAddWord(t *testing.T) {
	m := NewAbuseMasker("*")

	m.AddWords([]string{"bad", "terrible"})
	if got := m.Mask("This is a bad word"); got != "This is a *** word" {
		t.Errorf("AddWords mask = %q, want %q", got, "This is a *** word")
	}

	m.AddWord("awful")
	if got := m.Mask("This is an awful word"); got != "This is an ***** word" {
		t.Errorf("AddWord mask = %q, want %q", got, "This is an ***** word")
	}
}

// TestAbuseMasker_ContainsAbuse 對照 v2 TestAbuseMaskerContainsAbuse。
func TestAbuseMasker_ContainsAbuse(t *testing.T) {
	m := NewAbuseMaskerWithWords("*", []string{"bad", "terrible"})
	if !m.ContainsAbuse("This is a bad word") {
		t.Error("應偵測到敏感詞")
	}
	if m.ContainsAbuse("This is a good word") {
		t.Error("乾淨文字不應偵測到敏感詞")
	}
}

// TestAbuseMasker_GetAbuseWords 對照 v2 TestAbuseMaskerGetAbuseWords。
func TestAbuseMasker_GetAbuseWords(t *testing.T) {
	m := NewAbuseMaskerWithWords("*", []string{"bad", "terrible", "awful"})
	got := m.GetAbuseWords("This is a bad word and terrible")

	want := []string{"bad", "terrible"}
	if len(got) != len(want) {
		t.Fatalf("GetAbuseWords() len = %d, want %d", len(got), len(want))
	}
	for _, w := range want {
		found := false
		for _, g := range got {
			if g == w {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("應找到敏感詞: %s", w)
		}
	}
}

// TestAbuseMasker_MultipleOccurrences 對照 v2 TestAbuseMaskerMultipleOccurrences。
func TestAbuseMasker_MultipleOccurrences(t *testing.T) {
	m := NewAbuseMaskerWithWords("*", []string{"bad"})
	if got := m.Mask("This is bad and that is also bad"); got != "This is *** and that is also ***" {
		t.Errorf("multiple occurrences = %q, want %q", got, "This is *** and that is also ***")
	}
}

// TestAbuseMasker_EmptyTrieBehavior 對照 v2 TestAbuseMaskerEmptyTrie：空字典原樣回傳且偵測為空。
func TestAbuseMasker_EmptyTrieBehavior(t *testing.T) {
	m := NewAbuseMasker("*")
	text := "This is a test"

	if got := m.Mask(text); got != text {
		t.Errorf("empty trie mask = %q, want %q", got, text)
	}
	if m.ContainsAbuse(text) {
		t.Error("空字典不應偵測到敏感詞")
	}
	if got := m.GetAbuseWords(text); len(got) != 0 {
		t.Error("空字典應回傳無敏感詞")
	}
}

// TestAbuseMasker_NoSubstringMatchParity 對照 v2 TestAbuseMasker_NoSubstringMatch。
func TestAbuseMasker_NoSubstringMatchParity(t *testing.T) {
	m := NewAbuseMaskerWithWords("*", []string{"bad"})
	if got := m.Mask("bad badminton"); got != "*** badminton" {
		t.Errorf("no-substring-match = %q, want %q", got, "*** badminton")
	}
}

// TestAbuseMasker_PreservesWhitespaceParity 對照 v2 TestAbuseMasker_PreservesWhitespace。
func TestAbuseMasker_PreservesWhitespaceParity(t *testing.T) {
	m := NewAbuseMaskerWithWords("*", []string{"bad"})
	if got := m.Mask("  bad  good  bad  "); got != "  ***  good  ***  " {
		t.Errorf("preserve whitespace = %q, want %q", got, "  ***  good  ***  ")
	}
}

// TestCleanWordParity 對照 v2 TestCleanWord：去標點、轉小寫、保留字母數字。
func TestCleanWordParity(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"bad", "bad"},
		{"BAD", "bad"},
		{"bad!", "bad"},
		{"bad.", "bad"},
		{"bad,", "bad"},
		{"bad?", "bad"},
		{"bad@", "bad"},
		{"bad#", "bad"},
		{"bad-", "bad"},
		{"bad_", "bad"},
		{"bad123", "bad123"},
		{"123bad", "123bad"},
		{"", ""},
		{"   ", ""},
		{"   bad   ", "bad"},
	}
	for _, tt := range tests {
		if got := cleanWord(tt.input); got != tt.expected {
			t.Errorf("cleanWord(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// TestStruct_AbuseField 對照 AC-14：欄位 mask:"abuse" 經 Struct() 正確遮罩。
func TestStruct_AbuseField(t *testing.T) {
	type comment struct {
		Body string `mask:"abuse"`
	}

	m := NewMaskerMarshaler()
	abuse, err := m.Get(TypeAbuse)
	if err != nil {
		t.Fatalf("get abuse error: %v", err)
	}
	abuse.(*AbuseMasker).AddWords([]string{"bad", "terrible"})

	out, err := m.Struct(&comment{Body: "this is bad and terrible"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.(*comment).Body; got != "this is *** and ********" {
		t.Errorf("abuse field = %q, want %q", got, "this is *** and ********")
	}
}
