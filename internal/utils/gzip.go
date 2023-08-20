package utils

import (
	"io"
	"net/http"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (gzw *gzipResponseWriter) Write(data []byte) (int, error) {
	return gzw.Writer.Write(data)
}
