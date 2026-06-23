package ginmasker

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewConfigDefaults(t *testing.T) {
	cfg := newConfig()
	if cfg.maxBodySize != defaultMaxBodySize {
		t.Errorf("maxBodySize = %d, want %d", cfg.maxBodySize, defaultMaxBodySize)
	}
	if cfg.maskChar != defaultMaskChar {
		t.Errorf("maskChar = %q, want %q", cfg.maskChar, defaultMaskChar)
	}
	if cfg.logger == nil {
		t.Fatal("logger is nil, want zap.NewNop()")
	}
	if len(cfg.keywords) != 0 || len(cfg.queryKeys) != 0 || len(cfg.headerKeys) != 0 {
		t.Error("default slices should be empty")
	}
}

func TestWithMaxBodySizeIgnoresNonPositive(t *testing.T) {
	cfg := newConfig(WithMaxBodySize(0), WithMaxBodySize(-100))
	if cfg.maxBodySize != defaultMaxBodySize {
		t.Errorf("maxBodySize = %d, want default %d (n<=0 ignored)", cfg.maxBodySize, defaultMaxBodySize)
	}
	cfg = newConfig(WithMaxBodySize(2048))
	if cfg.maxBodySize != 2048 {
		t.Errorf("maxBodySize = %d, want 2048", cfg.maxBodySize)
	}
}

func TestWithLoggerNilKeepsNop(t *testing.T) {
	cfg := newConfig(WithLogger(nil))
	if cfg.logger == nil {
		t.Fatal("logger is nil after WithLogger(nil), want zap.NewNop()")
	}
	custom := zap.NewExample()
	cfg = newConfig(WithLogger(custom))
	if cfg.logger != custom {
		t.Error("WithLogger(custom) did not set custom logger")
	}
}

func TestWithKeywordsQueryHeader(t *testing.T) {
	cfg := newConfig(
		WithKeywords("password", "token"),
		WithQueryMask("access_token"),
		WithHeaderMask("Authorization"),
	)
	if len(cfg.keywords) != 2 {
		t.Errorf("keywords = %v, want 2", cfg.keywords)
	}
	if len(cfg.queryKeys) != 1 || cfg.queryKeys[0] != "access_token" {
		t.Errorf("queryKeys = %v", cfg.queryKeys)
	}
	if len(cfg.headerKeys) != 1 || cfg.headerKeys[0] != "Authorization" {
		t.Errorf("headerKeys = %v", cfg.headerKeys)
	}
}
