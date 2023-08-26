package compression

import (
	"io"
	"net/http"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

// желательно создать конструктор
func (gzw *gzipResponseWriter) Write(data []byte) (int, error) {
	return gzw.Writer.Write(data)
}
