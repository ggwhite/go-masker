package ginmasker

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	masker "github.com/ggwhite/go-masker/v3"
)

func TestMaskJSONBody(t *testing.T) {
	raw := []byte(`{"user":"alice","password":"secret123","nested":{"token":"abcdef"}}`)
	got, ok := maskJSONBody(raw, []string{"password", "token"}, "*")
	if !ok {
		t.Fatalf("maskJSONBody ok = false, want true")
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("unmarshal masked: %v", err)
	}
	if m["user"] != "alice" {
		t.Errorf("user = %v, want alice", m["user"])
	}
	if m["password"] != "*********" {
		t.Errorf("password = %v, want 9 stars", m["password"])
	}
	nested := m["nested"].(map[string]any)
	if nested["token"] != "******" {
		t.Errorf("nested.token = %v, want 6 stars", nested["token"])
	}
}

func TestMaskJSONBodyInvalid(t *testing.T) {
	got, ok := maskJSONBody([]byte("not json"), []string{"password"}, "*")
	if ok || got != "" {
		t.Errorf("maskJSONBody(invalid) = (%q, %v), want (\"\", false)", got, ok)
	}
}

func TestWalkAndMaskArray(t *testing.T) {
	var node any
	_ = json.Unmarshal([]byte(`[{"password":"a"},{"name":"bob"}]`), &node)
	out := walkAndMask(node, []string{"password"}, "*")
	arr := out.([]any)
	first := arr[0].(map[string]any)
	if first["password"] != "*" {
		t.Errorf("password = %v, want 1 star", first["password"])
	}
	second := arr[1].(map[string]any)
	if second["name"] != "bob" {
		t.Errorf("name = %v, want bob (not matched)", second["name"])
	}
}

func TestKeyMatches(t *testing.T) {
	cases := []struct {
		key      string
		keywords []string
		want     bool
	}{
		{"Password", []string{"password"}, true},
		{"access_token", []string{"TOKEN"}, true},
		{"username", []string{"password"}, false},
		{"anything", nil, false},
	}
	for _, c := range cases {
		if got := keyMatches(c.key, c.keywords); got != c.want {
			t.Errorf("keyMatches(%q, %v) = %v, want %v", c.key, c.keywords, got, c.want)
		}
	}
}

// TestMaskStringMatchesCore 驗證預設 maskChar 的 string 全遮與 v3 core 的 masker.All 逐字一致。
func TestMaskStringMatchesCore(t *testing.T) {
	for _, s := range []string{"secret", "", "中文密碼", "a"} {
		if got, want := maskString(s, "*"), masker.All(s); got != want {
			t.Errorf("maskString(%q, *) = %q, want %q (masker.All)", s, got, want)
		}
	}
}

func TestMaskStringCustomChar(t *testing.T) {
	if got := maskString("abc", "#"); got != "###" {
		t.Errorf("maskString(abc, #) = %q, want ###", got)
	}
	// rune 數而非 byte 數
	if got := maskString("中", "#"); got != "#" {
		t.Errorf("maskString(中, #) = %q, want #", got)
	}
}

func TestMaskValueNonString(t *testing.T) {
	cases := []any{float64(42), true, map[string]any{"a": 1}, []any{1, 2}}
	for _, v := range cases {
		if got := maskValue(v, "*"); got != "****" {
			t.Errorf("maskValue(%v) = %v, want **** (fixed redaction)", v, got)
		}
	}
}

func TestMaskedQuery(t *testing.T) {
	values := url.Values{"token": {"abcdef"}, "page": {"1"}}
	got := maskedQuery(values, []string{"token"}, "*")
	parsed, _ := url.ParseQuery(got)
	if parsed.Get("token") != "******" {
		t.Errorf("token = %q, want 6 stars", parsed.Get("token"))
	}
	if parsed.Get("page") != "1" {
		t.Errorf("page = %q, want 1 (untouched)", parsed.Get("page"))
	}
}

func TestMaskedQueryDoesNotMutateInput(t *testing.T) {
	values := url.Values{"token": {"abcdef"}}
	_ = maskedQuery(values, []string{"token"}, "*")
	if values.Get("token") != "abcdef" {
		t.Errorf("input mutated: token = %q, want abcdef", values.Get("token"))
	}
}

func TestMaskedHeaders(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer xyz")
	h.Set("X-Request-Id", "req-1")
	// case-insensitive 精確比對命中 Authorization
	got := maskedHeaders(h, []string{"authorization"}, "*")
	want := "Authorization=**********; X-Request-Id=req-1"
	if got != want {
		t.Errorf("maskedHeaders = %q, want %q", got, want)
	}
	// 不可修改原 header
	if h.Get("Authorization") != "Bearer xyz" {
		t.Errorf("original header mutated: %q", h.Get("Authorization"))
	}
}

func TestMaskedHeadersNilSafe(t *testing.T) {
	if got := maskedHeaders(nil, []string{"authorization"}, "*"); got != "" {
		t.Errorf("maskedHeaders(nil) = %q, want empty", got)
	}
}

func TestWalkAndMaskScalarRoot(t *testing.T) {
	var node any
	_ = json.Unmarshal([]byte(`"plain"`), &node)
	out := walkAndMask(node, []string{"password"}, "*")
	if !reflect.DeepEqual(out, "plain") {
		t.Errorf("scalar root = %v, want plain unchanged", out)
	}
}
