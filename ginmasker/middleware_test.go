package ginmasker

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// echoHandler 讀取 request body 後原樣回傳，用於驗證下游仍讀得到完整 body、
// 且 client 收到的是未遮罩內容。
func echoHandler(c *gin.Context) {
	b, _ := io.ReadAll(c.Request.Body)
	c.Data(http.StatusOK, "application/json", b)
}

func newObservedLogger(level zapcore.Level) (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(level)
	return zap.New(core), logs
}

func doRequest(engine *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func lastEntryField(t *testing.T, logs *observer.ObservedLogs, key string) (any, bool) {
	t.Helper()
	all := logs.All()
	if len(all) == 0 {
		return nil, false
	}
	v, ok := all[len(all)-1].ContextMap()[key]
	return v, ok
}

func TestMiddlewareRequestBodyMaskedAndReplayed(t *testing.T) {
	logger, logs := newObservedLogger(zapcore.InfoLevel)
	r := gin.New()
	r.Use(Middleware(WithLogger(logger), WithKeywords("password")))
	r.POST("/echo", echoHandler)

	body := `{"password":"secret","user":"alice"}`
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := doRequest(r, req)

	// 下游 handler 讀得到完整 body → echo 回的是未遮罩原文
	if w.Body.String() != body {
		t.Errorf("client body = %q, want unmasked %q", w.Body.String(), body)
	}
	// log 內 request_body 已遮罩
	v, ok := lastEntryField(t, logs, "request_body")
	if !ok {
		t.Fatal("request_body field missing")
	}
	masked := v.(string)
	if strings.Contains(masked, "secret") {
		t.Errorf("request_body leaked secret: %q", masked)
	}
	if !strings.Contains(masked, `"password":"******"`) {
		t.Errorf("request_body not masked: %q", masked)
	}
}

func TestMiddlewareResponseBodyMasked(t *testing.T) {
	logger, logs := newObservedLogger(zapcore.InfoLevel)
	r := gin.New()
	r.Use(Middleware(WithLogger(logger), WithKeywords("token")))
	r.GET("/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"token": "abcdef", "ok": true})
	})

	w := doRequest(r, httptest.NewRequest(http.MethodGet, "/data", nil))

	// client 收到未遮罩 response
	if !strings.Contains(w.Body.String(), "abcdef") {
		t.Errorf("client should receive unmasked response, got %q", w.Body.String())
	}
	v, ok := lastEntryField(t, logs, "response_body")
	if !ok {
		t.Fatal("response_body field missing")
	}
	masked := v.(string)
	if strings.Contains(masked, "abcdef") {
		t.Errorf("response_body leaked token: %q", masked)
	}
	if !strings.Contains(masked, `"token":"******"`) {
		t.Errorf("response_body not masked: %q", masked)
	}
}

func TestMiddlewareGateSkipsWhenLevelDisabled(t *testing.T) {
	// observer 只收 Error 以上 → Info 被停用，middleware 應整段跳過
	logger, logs := newObservedLogger(zapcore.ErrorLevel)
	r := gin.New()
	r.Use(Middleware(WithLogger(logger), WithKeywords("password")))
	r.POST("/echo", echoHandler)

	body := `{"password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := doRequest(r, req)

	if logs.Len() != 0 {
		t.Errorf("expected no log entries when Info disabled, got %d", logs.Len())
	}
	// 即使跳過，下游仍須讀得到完整 body
	if w.Body.String() != body {
		t.Errorf("client body = %q, want %q", w.Body.String(), body)
	}
}

func TestMiddlewareNonJSONRequest(t *testing.T) {
	logger, logs := newObservedLogger(zapcore.InfoLevel)
	r := gin.New()
	r.Use(Middleware(WithLogger(logger), WithKeywords("password")))
	r.POST("/echo", echoHandler)

	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader("password=secret"))
	req.Header.Set("Content-Type", "text/plain")
	doRequest(r, req)

	v, _ := lastEntryField(t, logs, "request_body")
	if v.(string) != markerNonJSON {
		t.Errorf("request_body = %q, want %q", v, markerNonJSON)
	}
}

func TestMiddlewareRequestTooLarge(t *testing.T) {
	logger, logs := newObservedLogger(zapcore.InfoLevel)
	r := gin.New()
	r.Use(Middleware(WithLogger(logger), WithMaxBodySize(10), WithKeywords("password")))
	r.POST("/echo", echoHandler)

	body := `{"password":"this-is-a-very-long-secret-value"}`
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := doRequest(r, req)

	// 超限仍 replay 完整 body 給下游
	if w.Body.String() != body {
		t.Errorf("client body = %q, want full %q", w.Body.String(), body)
	}
	v, _ := lastEntryField(t, logs, "request_body")
	if v.(string) != markerTooLarge {
		t.Errorf("request_body = %q, want %q", v, markerTooLarge)
	}
}

func TestMiddlewareResponseTooLarge(t *testing.T) {
	logger, logs := newObservedLogger(zapcore.InfoLevel)
	r := gin.New()
	r.Use(Middleware(WithLogger(logger), WithMaxBodySize(10)))
	large := strings.Repeat("x", 100)
	r.GET("/data", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte(large))
	})

	w := doRequest(r, httptest.NewRequest(http.MethodGet, "/data", nil))

	// client 收到完整 response
	if w.Body.String() != large {
		t.Errorf("client should receive full response of len %d, got %d", len(large), w.Body.Len())
	}
	v, _ := lastEntryField(t, logs, "response_body")
	if v.(string) != markerTooLarge {
		t.Errorf("response_body = %q, want %q", v, markerTooLarge)
	}
}

func TestMiddlewareQueryAndHeaderMask(t *testing.T) {
	logger, logs := newObservedLogger(zapcore.InfoLevel)
	r := gin.New()
	r.Use(Middleware(
		WithLogger(logger),
		WithQueryMask("token"),
		WithHeaderMask("Authorization"),
	))
	r.GET("/data", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/data?token=abcdef&page=1", nil)
	req.Header.Set("Authorization", "Bearer xyz")
	doRequest(r, req)

	q, ok := lastEntryField(t, logs, "query")
	if !ok {
		t.Fatal("query field missing")
	}
	parsed, _ := url.ParseQuery(q.(string))
	if parsed.Get("token") == "abcdef" || parsed.Get("token") == "" {
		t.Errorf("query token not masked: %q", parsed.Get("token"))
	}
	if parsed.Get("page") != "1" {
		t.Errorf("query page = %q, want 1", parsed.Get("page"))
	}

	h, ok := lastEntryField(t, logs, "headers")
	if !ok {
		t.Fatal("headers field missing")
	}
	if strings.Contains(h.(string), "Bearer xyz") {
		t.Errorf("headers leaked Authorization: %q", h)
	}
	if !strings.Contains(h.(string), "Authorization=") {
		t.Errorf("headers missing Authorization key: %q", h)
	}
}

func TestMiddlewareNoQueryHeaderFieldsWhenUnset(t *testing.T) {
	logger, logs := newObservedLogger(zapcore.InfoLevel)
	r := gin.New()
	r.Use(Middleware(WithLogger(logger)))
	r.GET("/data", func(c *gin.Context) { c.Status(http.StatusOK) })

	doRequest(r, httptest.NewRequest(http.MethodGet, "/data?token=abc", nil))

	if _, ok := lastEntryField(t, logs, "query"); ok {
		t.Error("query field should be absent when WithQueryMask unset")
	}
	if _, ok := lastEntryField(t, logs, "headers"); ok {
		t.Error("headers field should be absent when WithHeaderMask unset")
	}
}

func TestMiddlewareDefaultLoggerNop(t *testing.T) {
	// 未配置 logger 時預設 NewNop，Info 停用 → 不應 panic、handler 正常運作
	r := gin.New()
	r.Use(Middleware())
	r.POST("/echo", echoHandler)

	body := `{"password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := doRequest(r, req)

	if w.Body.String() != body {
		t.Errorf("client body = %q, want %q", w.Body.String(), body)
	}
}
