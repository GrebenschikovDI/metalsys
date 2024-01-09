package crypto

import (
	"bytes"
	"io"
	"net/http"
)

func DecryptMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if MetDecryptor != nil {
				buf, _ := io.ReadAll(r.Body)
				message, err := MetDecryptor.Decrypt(buf)
				if err != nil {
					http.Error(w, "Decryption error", http.StatusBadRequest)
					return
				}
				rwd := r.Clone(r.Context())
				rwd.Body = io.NopCloser(bytes.NewBuffer(message))
				next.ServeHTTP(w, rwd)
			}
			next.ServeHTTP(w, r)
		})
	}
}
