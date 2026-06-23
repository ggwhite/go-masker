package masker

import (
	"strings"
	"testing"
)

func TestAbuseMasker_Mask(t *testing.T) {
	m := NewAbuseMaskerWithWords("*", []string{"bad", "terrible"})
	got := m.Mask("This is a bad word and terrible")
	want := "This is a *** word and ********"
	if got != want {
		t.Errorf("AbuseMasker.Mask() = %v, want %v", got, want)
	}
}

func TestAbuseMasker_EmptyDictionary(t *testing.T) {
	m := NewAbuseMasker("*")
	if got := m.Mask("hello world"); got != "hello world" {
		t.Errorf("empty dict Mask() = %v, want hello world", got)
	}
}

func TestAbuseMasker_AddWords(t *testing.T) {
	m := NewAbuseMasker("#")
	m.AddWord("foo")
	m.AddWords([]string{"bar"})
	if got := m.Mask("foo bar baz"); got != "### ### baz" {
		t.Errorf("Mask() = %v, want ### ### baz", got)
	}
}

// TestAbuseConvenience_EmptyDict 驗證 AC-9b：default 空字典下 Abuse 回原值。
func TestAbuseConvenience_EmptyDict(t *testing.T) {
	if got := Abuse("hello"); got != "hello" {
		t.Errorf("Abuse() with empty dict = %v, want hello", got)
	}
}

// TestAbuse_LoadedDictionary 驗證 AC-9b：透過 loader 載入詞典後命中詞被遮罩。
func TestAbuse_LoadedDictionary(t *testing.T) {
	loader := NewAbuseWordLoader()
	words, err := loader.LoadFromReader(strings.NewReader("bad\n# comment\nterrible\n"))
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	m := NewMaskerMarshaler()
	abuse, err := m.Get(TypeAbuse)
	if err != nil {
		t.Fatalf("get abuse error: %v", err)
	}
	abuse.(*AbuseMasker).AddWords(words)

	got, err := m.Marshal(TypeAbuse, "a bad day")
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if got != "a *** day" {
		t.Errorf("Marshal(TypeAbuse) = %v, want a *** day", got)
	}
}

func TestAbuseTrie_Contains(t *testing.T) {
	trie := NewAbuseTrie()
	trie.InsertAll([]string{"hello", "world"})
	if !trie.Contains("hello") {
		t.Errorf("Contains(hello) = false, want true")
	}
	if trie.Contains("nope") {
		t.Errorf("Contains(nope) = true, want false")
	}
}
