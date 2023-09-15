package hash

import (
	"bytes"
	"crypto/hmac"
	"io"
	"net/http"
)

type hashResponseWriter struct {
	http.ResponseWriter
	capturedData []byte
}

func ValidateHash(key string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			receivedHash := r.Header.Get("HashSHA256")
			data, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
				return
			}
			localHash := Sign(data, key)
			if !hmac.Equal([]byte(receivedHash), []byte(localHash)) {
				http.Error(w, "Несоответствие хешей", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(data))
			hw := hashResponseWriter{w, data}

			next.ServeHTTP(&hw, r)

			responseHash := Sign(hw.capturedData, key)
			w.Header().Set("HashSHA256", responseHash)
		})
	}
}
