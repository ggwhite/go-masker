package ginmasker

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// bodyCaptureWriter 內嵌 gin.ResponseWriter，於寫出 response 時同步複製一份 body 副本，
// 供 middleware 組 log 使用。副本最多累積到 limit+1 bytes：超過即停止累積並標記 tooLarge，
// 但仍照常寫給原 writer，確保 client 收到完整未遮罩的 response。
type bodyCaptureWriter struct {
	gin.ResponseWriter
	buf      bytes.Buffer
	limit    int64
	tooLarge bool
}

// newBodyCaptureWriter 以原 writer 與累積上限建立 bodyCaptureWriter。
func newBodyCaptureWriter(w gin.ResponseWriter, limit int64) *bodyCaptureWriter {
	return &bodyCaptureWriter{ResponseWriter: w, limit: limit}
}

// Write 把 bytes 寫給原 writer，同時在未超限時複製進內部 buffer。
func (w *bodyCaptureWriter) Write(b []byte) (int, error) {
	w.capture(b)
	return w.ResponseWriter.Write(b)
}

// WriteString 把字串寫給原 writer，同時在未超限時複製進內部 buffer。
func (w *bodyCaptureWriter) WriteString(s string) (int, error) {
	w.capture([]byte(s))
	return w.ResponseWriter.WriteString(s)
}

// capture 在尚未超限時把 b 累積進 buffer，累積量達 limit+1 即停止並標記 tooLarge。
func (w *bodyCaptureWriter) capture(b []byte) {
	if w.tooLarge {
		return
	}
	remaining := w.limit + 1 - int64(w.buf.Len())
	if remaining <= 0 {
		w.tooLarge = true
		return
	}
	if int64(len(b)) > remaining {
		w.buf.Write(b[:remaining])
		w.tooLarge = true
		return
	}
	w.buf.Write(b)
}

// captured 回傳目前累積的 response body 副本與是否超限。
func (w *bodyCaptureWriter) captured() ([]byte, bool) {
	return w.buf.Bytes(), w.tooLarge || int64(w.buf.Len()) > w.limit
}
