// Package ginmasker 提供 go-masker v3 與 Gin 的整合 middleware，
// 在寫入 access log 前自動遮罩 request／response body、query、header 中的敏感欄位。
//
// 用法概要：
//
//	r := gin.New()
//	r.Use(ginmasker.Middleware(
//	    ginmasker.WithLogger(logger),
//	    ginmasker.WithKeywords("password", "token", "secret"),
//	    ginmasker.WithQueryMask("access_token"),
//	    ginmasker.WithHeaderMask("Authorization"),
//	))
//
// 只遮 log、不改 body：middleware 攔截到的 body 僅用於組 log 欄位，
// 原始 request body 會完整 replay 回 c.Request.Body 供下游 handler 讀取，
// 原始 response 也會完整寫回 client，遮罩結果絕不外洩到實際資料流。
//
// 依賴隔離設計：gin 與 zap 依賴被刻意限制在這個獨立 sub-module 內，
// 不污染 v3 core（v3/go.mod 不含 gin／zap）。只有需要 Gin 整合的使用者才會引入這些依賴。
//
// keyword 遮罩語意一律委派 v3 core 的 package-level 函式（如 masker.All），
// 確保輸出與 v3 core 逐字一致。
//
// ⚠️ keyword 比對本質是「靠 key 名稱猜測敏感欄位」，必然不精準：
//   - false positive：key 名碰巧含關鍵字（如 password_count）會被誤遮。
//   - false negative：敏感欄位的 key 名不在清單內則漏遮。
//
// 因此 keyword 比對只適合當作改不動業務 code 時的被動防線；
// 若業務端已在 handler 內以 Sensitive[T]、masker.Struct() 等顯式遮罩敏感欄位，
// 本層即為多餘，不應視為主力方案。
package ginmasker
