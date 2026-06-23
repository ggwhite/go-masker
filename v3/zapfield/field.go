package zapfield

import (
	masker "github.com/ggwhite/go-masker/v3"
	"go.uber.org/zap"
)

// 本檔為每個 v3 core package-level 遮罩函式提供對應的 zap.Field helper。
// 命名與 Sensitive 建構子家族一致（mobile 用 Phone，與 masker.NewPhone 對齊）。
// 一律以 zap.String(key, masked) 建構、遮罩值一律委派 core 函式，確保輸出逐字一致。

// Phone 回傳已遮罩手機號碼的 zap.Field，遮罩委派 masker.Mobile。
// Example:
//
//	logger.Info("login", zapfield.Phone("phone", "0987654321")) // phone=0987***321
func Phone(key, value string) zap.Field {
	return zap.String(key, masker.Mobile(value))
}

// Email 回傳已遮罩 email 的 zap.Field，遮罩委派 masker.Email。
// Example:
//
//	logger.Info("signup", zapfield.Email("email", "ggw.chang@gmail.com")) // email=ggw****@gmail.com
func Email(key, value string) zap.Field {
	return zap.String(key, masker.Email(value))
}

// Password 回傳已遮罩密碼的 zap.Field，遮罩委派 masker.Password（固定 14 個遮罩字元）。
// Example:
//
//	logger.Info("auth", zapfield.Password("pwd", "password")) // pwd=**************
func Password(key, value string) zap.Field {
	return zap.String(key, masker.Password(value))
}

// Name 回傳已遮罩姓名的 zap.Field，遮罩委派 masker.Name。
// Example:
//
//	logger.Info("profile", zapfield.Name("name", "John Doe")) // name=J**n D**e
func Name(key, value string) zap.Field {
	return zap.String(key, masker.Name(value))
}

// Address 回傳已遮罩地址的 zap.Field，遮罩委派 masker.Address。
// Example:
//
//	logger.Info("order", zapfield.Address("addr", "台北市內湖區內湖路一段737巷1號1樓")) // addr=台北市內湖區******
func Address(key, value string) zap.Field {
	return zap.String(key, masker.Address(value))
}

// ID 回傳已遮罩證號的 zap.Field，遮罩委派 masker.ID。
// Example:
//
//	logger.Info("kyc", zapfield.ID("id", "A123456789")) // id=A12345****
func ID(key, value string) zap.Field {
	return zap.String(key, masker.ID(value))
}

// Credit 回傳已遮罩信用卡號的 zap.Field，遮罩委派 masker.Credit。
// Example:
//
//	logger.Info("pay", zapfield.Credit("card", "4111111111111111")) // card=411111******1111
func Credit(key, value string) zap.Field {
	return zap.String(key, masker.Credit(value))
}

// Tel 回傳已遮罩市話的 zap.Field，遮罩委派 masker.Tel。
// Example:
//
//	logger.Info("contact", zapfield.Tel("tel", "0227993078")) // tel=(02)2799-****
func Tel(key, value string) zap.Field {
	return zap.String(key, masker.Tel(value))
}

// URL 回傳已遮罩 URL 密碼段的 zap.Field，遮罩委派 masker.URL。
// Example:
//
//	logger.Info("hook", zapfield.URL("url", "http://john:password@localhost:3000")) // url=http://john:xxxxx@localhost:3000
func URL(key, value string) zap.Field {
	return zap.String(key, masker.URL(value))
}

// Abuse 回傳已遮罩敏感詞的 zap.Field，遮罩委派 masker.Abuse。
// 注意：DefaultMaskerMarshaler 的字典預設為空，未載入詞典時原樣輸出。
// Example:
//
//	logger.Info("msg", zapfield.Abuse("text", "hello")) // text=hello（空字典）
func Abuse(key, value string) zap.Field {
	return zap.String(key, masker.Abuse(value))
}

// None 回傳原值（不遮罩）的 zap.Field，行為委派 masker.None。
// Example:
//
//	logger.Info("meta", zapfield.None("note", "secret")) // note=secret
func None(key, value string) zap.Field {
	return zap.String(key, masker.None(value))
}

// All 回傳整串遮罩的 zap.Field，遮罩委派 masker.All（每個字元都換成遮罩字元）。
// Example:
//
//	logger.Info("token", zapfield.All("v", "secret")) // v=******
func All(key, value string) zap.Field {
	return zap.String(key, masker.All(value))
}
