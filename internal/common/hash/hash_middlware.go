package hash

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"io"
	"net/http"
)

type hashResponseWriter struct {
	http.ResponseWriter
	key string
}

func ValidateHash(key string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			if r.Header.Get("HashSHA256") == "" {
				http.Error(w, "Пустой заголовок HashSHA256", http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
				return
			}
			bodyHash := Sign(body, key)
			headerHash, err := hex.DecodeString(r.Header.Get("HashSHA256"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if !hmac.Equal(headerHash, bodyHash) {
				http.Error(w, "Несоответствие хешей", http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			hw := &hashResponseWriter{
				ResponseWriter: w,
				key:            key,
			}
			next.ServeHTTP(hw, r)
		})
	}
}

func (h *hashResponseWriter) Write(b []byte) (int, error) {
	hash := Sign(b, h.key)
	h.Header().Set("HashSHA256", hex.EncodeToString(hash))
	return h.ResponseWriter.Write(b)
}
