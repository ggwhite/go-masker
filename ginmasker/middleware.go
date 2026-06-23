package ginmasker

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// log 欄位與標記常數。
const (
	markerTooLarge = "<omitted: too large>"
	markerNonJSON  = "<non-json>"
	contentJSON    = "application/json"
	logMessage     = "access"
)

// Middleware 回傳一個 Gin middleware，於寫入 access log 前自動遮罩 request／response body、
// query 與 header 中的敏感欄位。
//
// 只影響 log 輸出：原始 request body 會完整 replay 回 c.Request.Body 供下游 handler 讀取，
// 原始 response 也會完整寫回 client，遮罩結果不外洩到實際資料流。
//
// 效能考量：僅在 logger 的 Info level 啟用時才讀取與遮罩 body；level 未啟用時直接放行。
//
// 用法：
//
//	r.Use(ginmasker.Middleware(
//	    ginmasker.WithLogger(logger),
//	    ginmasker.WithKeywords("password", "token"),
//	))
func Middleware(opts ...Option) gin.HandlerFunc {
	cfg := newConfig(opts...)
	return func(c *gin.Context) {
		if !cfg.logger.Core().Enabled(zapcore.InfoLevel) {
			c.Next()
			return
		}

		start := time.Now()
		requestBody := cfg.captureRequestBody(c)

		capture := newBodyCaptureWriter(c.Writer, cfg.maxBodySize)
		c.Writer = capture

		c.Next()

		cfg.writeLog(c, start, requestBody, capture)
	}
}

// captureRequestBody 讀取並遮罩 request body（僅 application/json），同時把完整原始 body
// replay 回 c.Request.Body 供下游使用；回傳供 log 使用的字串標記。
func (cfg *config) captureRequestBody(c *gin.Context) string {
	if !strings.Contains(c.ContentType(), contentJSON) {
		return markerNonJSON
	}

	limited := io.LimitReader(c.Request.Body, cfg.maxBodySize+1)
	head, err := io.ReadAll(limited)
	if err != nil {
		// 讀取失敗：把已讀部分接回剩餘 body，避免吃掉下游 body。
		c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(head), c.Request.Body))
		return markerNonJSON
	}

	// replay：把已讀的 head 接回原 body 剩餘部分，確保下游 handler 讀得到完整 body。
	c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(head), c.Request.Body))

	if int64(len(head)) > cfg.maxBodySize {
		return markerTooLarge
	}

	masked, ok := maskJSONBody(head, cfg.keywords, cfg.maskChar)
	if !ok {
		return markerNonJSON
	}
	return masked
}

// captureResponseBody 依 response content-type 與是否超限決定 log 用的 response 字串標記。
func (cfg *config) captureResponseBody(c *gin.Context, capture *bodyCaptureWriter) string {
	body, tooLarge := capture.captured()
	if tooLarge {
		return markerTooLarge
	}
	if !strings.Contains(c.Writer.Header().Get("Content-Type"), contentJSON) {
		return markerNonJSON
	}
	masked, ok := maskJSONBody(body, cfg.keywords, cfg.maskChar)
	if !ok {
		return markerNonJSON
	}
	return masked
}

// writeLog 組一筆 access log（固定 message "access"、Info level）並輸出。
// 全部使用明確型別的 zap.Field，不使用 zap.Any／zap.Reflect。
func (cfg *config) writeLog(c *gin.Context, start time.Time, requestBody string, capture *bodyCaptureWriter) {
	fields := make([]zap.Field, 0, 9)
	fields = append(fields,
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", c.Writer.Status()),
		zap.Duration("latency", time.Since(start)),
		zap.String("client_ip", c.ClientIP()),
		zap.String("request_body", requestBody),
		zap.String("response_body", cfg.captureResponseBody(c, capture)),
	)

	if len(cfg.queryKeys) > 0 {
		fields = append(fields, zap.String("query", maskedQuery(c.Request.URL.Query(), cfg.queryKeys, cfg.maskChar)))
	}
	if len(cfg.headerKeys) > 0 {
		fields = append(fields, zap.String("headers", maskedHeaders(c.Request.Header, cfg.headerKeys, cfg.maskChar)))
	}

	cfg.logger.Info(logMessage, fields...)
}
