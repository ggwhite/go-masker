// Package zapfield 提供 go-masker v3 與 go.uber.org/zap 的整合 helper。
//
// 它讓 zap 使用者免去「先手動呼叫 masker.Mobile(...) 再包 zap.String(...)」的轉接：
// 直接寫 zapfield.Phone("phone", rawValue) 即可取得已遮罩的 zap.Field；
// 也提供 Sensitive adapter，把 masker.Sensitive[T] 直接轉成已遮罩的 zap.Field。
//
// 依賴隔離設計：zap 依賴被刻意限制在這個獨立 sub-module 內，
// 不污染 v3 core（v3/go.mod 不含 zap）。只有需要 zap 整合的使用者才會引入 zap 依賴。
//
// 所有 helper 一律以 zap.String(key, masked) 建構，不使用 zap.Any／zap.Reflect，
// 避免 reflection 與原值洩漏；遮罩值一律委派 v3 core 的 package-level 函式，
// 確保輸出與 v3 core 逐字一致。
package zapfield
