package ginservice

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// gin 服务扩展
type ginResWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ginResWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ginResWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
