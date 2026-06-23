package masker

// DefaultRedactText 是 redact 模式下的預設替換文字，未指定 WithRedactText 時採用。
const DefaultRedactText = "[REDACTED]"

// SensitiveOption 設定 Sensitive[T] 的建構選項，傳給 NewSensitive 或任一內建建構子。
// 選項在建構時套用一次後即固定，之後不可變更，避免併發切換造成 race。
type SensitiveOption func(*sensitiveConfig)

// sensitiveConfig 收集建構時的選項，僅在建構流程內部使用。
type sensitiveConfig struct {
	redact     bool
	redactText string
}

// WithRedact 啟用完全替換模式：所有自動輸出（String／GoString／MarshalJSON／MarshalText／LogValue）
// 改回傳替換文字（預設 [REDACTED]）而非部分遮罩，適用於合規或對外 API response 場景。
// Reveal 與 SQL 的 Value 不受影響，仍回傳原值。
// Example:
//
//	s := masker.NewPhone("0987654321", masker.WithRedact())
//	fmt.Println(s)   // [REDACTED]
//	_ = s.Reveal()   // 0987654321
func WithRedact() SensitiveOption {
	return func(c *sensitiveConfig) {
		c.redact = true
	}
}

// WithRedactText 啟用完全替換模式並自訂替換文字。隱含 WithRedact，單獨使用即可開啟 redact。
// Example:
//
//	s := masker.NewPhone("0987654321", masker.WithRedactText("***"))
//	fmt.Println(s) // ***
func WithRedactText(text string) SensitiveOption {
	return func(c *sensitiveConfig) {
		c.redact = true
		c.redactText = text
	}
}
