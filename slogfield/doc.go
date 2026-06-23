// Package slogfield 提供 go-masker v3 與標準函式庫 log/slog 的整合 helper。
//
// 它讓 slog 使用者免去「先手動呼叫 masker.Mobile(...) 再包 slog.String(...)」的轉接：
// 直接寫 slogfield.Phone("phone", rawValue) 即可取得已遮罩的 slog.Attr；
// 也提供 SensitiveAttr adapter，把 masker.Sensitive[T] 直接轉成已遮罩的 slog.Attr。
//
// 依賴隔離設計：本 module 只依賴 stdlib log/slog 與 v3 core，不引入任何第三方套件，
// 因此 go.sum 預期為空或不存在。core 透過 replace 指回本地，不會額外下載。
//
// 所有 helper 一律以 slog.String(key, masked) 建構，不使用 slog.Any，
// 避免值被 reflection 處理或誤觸 raw；遮罩值一律委派 v3 core 的 package-level 函式，
// 確保輸出與 v3 core（也與 zapfield）逐字一致。
package slogfield
