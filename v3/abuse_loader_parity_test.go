package masker

import "testing"

// 本檔補齊 v2 abuse_loader_test.go 中尚未涵蓋於 v3 的 parity 案例：
// 空輸入、純註解、混合內容。期望輸出與 v2 逐字一致。

// TestAbuseWordLoader_EmptyInput 對照 v2 TestAbuseWordLoaderEmptyInput。
func TestAbuseWordLoader_EmptyInput(t *testing.T) {
	loader := NewAbuseWordLoader()

	words, err := loader.LoadFromString("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(words) != 0 {
		t.Errorf("LoadFromString(\"\") = %v, want 0 words", words)
	}

	if got := loader.LoadFromSlice([]string{}); len(got) != 0 {
		t.Errorf("LoadFromSlice([]) = %v, want 0 words", got)
	}

	if got := loader.LoadFromSlice([]string{"", "   ", "\t", "\n"}); len(got) != 0 {
		t.Errorf("LoadFromSlice(blank) = %v, want 0 words", got)
	}
}

// TestAbuseWordLoader_OnlyComments 對照 v2 TestAbuseWordLoaderOnlyComments。
func TestAbuseWordLoader_OnlyComments(t *testing.T) {
	loader := NewAbuseWordLoader()
	content := `# This is a comment
# Another comment
# Yet another comment
`
	words, err := loader.LoadFromString(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(words) != 0 {
		t.Errorf("LoadFromString(only comments) = %v, want 0 words", words)
	}
}

// TestAbuseWordLoader_MixedContent 對照 v2 TestAbuseWordLoaderMixedContent。
func TestAbuseWordLoader_MixedContent(t *testing.T) {
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
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"bad", "terrible", "awful"}
	if len(words) != len(want) {
		t.Fatalf("LoadFromString len = %d, want %d", len(words), len(want))
	}
	for i, w := range want {
		if words[i] != w {
			t.Errorf("words[%d] = %s, want %s", i, words[i], w)
		}
	}
}
