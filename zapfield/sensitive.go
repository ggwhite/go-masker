package zapfield

import (
	masker "github.com/ggwhite/go-masker/v3"
	"go.uber.org/zap"
)

// Sensitive 把任意 masker.Sensitive[T] 轉成已遮罩的 zap.Field。
//
// 它只讀 Sensitive[T] 快取的 masked 值（透過 String()，與 fmt／json／slog 路徑同源），
// 絕不呼叫 Reveal()、絕不碰原值。zero-value Sensitive[T]（未經建構子、mask 為 nil）的
// String() 回空字串，adapter 即輸出 zap.String(key, "")，不洩漏原值（安全退化）。
//
// Example:
//
//	phone := masker.NewPhone("0987654321")
//	logger.Info("login", zapfield.Sensitive("phone", phone)) // phone=0987***321
func Sensitive[T any](key string, s masker.Sensitive[T]) zap.Field {
	return zap.String(key, s.String())
}
