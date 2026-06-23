package ginmasker

import "go.uber.org/zap"

// defaultMaxBodySize 為預設的 body 讀取／累積上限（1MB），避免讀爆記憶體。
const defaultMaxBodySize int64 = 1 << 20

// defaultMaskChar 為預設遮罩字元；維持 "*" 以與 v3 core 的 masker.All 輸出逐字一致。
const defaultMaskChar = "*"

// config 是 middleware 的內部設定，由 newConfig 套用預設值後再以 Option 覆寫。
type config struct {
	maxBodySize int64
	keywords    []string
	logger      *zap.Logger
	maskChar    string
	queryKeys   []string
	headerKeys  []string
}

// Option 以 functional options 方式調整 middleware 設定。
type Option func(*config)

// WithMaxBodySize 設定 request／response body 的讀取與累積上限（bytes）。
// 超過上限的 body 不做遮罩、log 改記長度與 too-large 標記。
// 防呆：n <= 0 時忽略本設定、保留預設 1MB，避免意外關閉保護。
func WithMaxBodySize(n int64) Option {
	return func(c *config) {
		if n <= 0 {
			return
		}
		c.maxBodySize = n
	}
}

// WithKeywords 設定要在 request／response body 內遮罩的 key 名稱清單。
// 比對為 case-insensitive 子字串（strings.Contains），命中的 value 一律全遮。
func WithKeywords(keywords ...string) Option {
	return func(c *config) {
		c.keywords = keywords
	}
}

// WithLogger 設定輸出 access log 的 zap.Logger。
// 防呆：傳入 nil 時保留預設 zap.NewNop()，不讓 logger 變成 nil。
func WithLogger(logger *zap.Logger) Option {
	return func(c *config) {
		if logger == nil {
			return
		}
		c.logger = logger
	}
}

// WithQueryMask 設定要遮罩的 query 參數名（精確比對、case-sensitive）。
// 未設定時不輸出 query log 欄位。
func WithQueryMask(keys ...string) Option {
	return func(c *config) {
		c.queryKeys = keys
	}
}

// WithHeaderMask 設定要遮罩的 header 名（精確比對、case-insensitive）。
// 未設定時不輸出 headers log 欄位。
func WithHeaderMask(keys ...string) Option {
	return func(c *config) {
		c.headerKeys = keys
	}
}

// newConfig 以預設值建立 config，再依序套用 opts。
func newConfig(opts ...Option) *config {
	c := &config{
		maxBodySize: defaultMaxBodySize,
		logger:      zap.NewNop(),
		maskChar:    defaultMaskChar,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
